package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	beaconconfig "github.com/phoreproject/synapse/beacon/config"
	beaconmodule "github.com/phoreproject/synapse/beacon/module"
	"github.com/phoreproject/synapse/cfg"
	"github.com/phoreproject/synapse/utils"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	shardconfig "github.com/phoreproject/synapse/shard/config"
	shardmodule "github.com/phoreproject/synapse/shard/module"

	validatorconfig "github.com/phoreproject/synapse/validator/config"
	validatormodule "github.com/phoreproject/synapse/validator/module"

	relayerconfig "github.com/phoreproject/synapse/relayer/config"
	relayermodule "github.com/phoreproject/synapse/relayer/module"

	_ "net/http/pprof"
)

// SynapseOptions are the options for all module configs.
type SynapseOptions struct {
	ModuleConfigs []AnyModuleConfig `yaml:"modules"`
}

// AnyModuleConfig is a wrapper around a module config allowing it to be unmarshalled from a string.
type AnyModuleConfig struct {
	Module interface{}
}

// UnmarshalYAML unmarshals YAML from
func (a *AnyModuleConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kvs map[string]interface{}

	if err := unmarshal(&kvs); err != nil {
		return err
	}

	t, found := kvs["type"]
	if !found {
		return fmt.Errorf("error while decoding YAML, module missing required parameter \"type\"")
	}

	var moduleValue reflect.Value
	var moduleType reflect.Type

	switch t {
	case "beacon":
		module := &beaconconfig.Options{}
		moduleValue = reflect.ValueOf(module).Elem()
		moduleType = reflect.TypeOf(module).Elem()
		a.Module = module
	case "shard":
		module := &shardconfig.Options{}
		moduleValue = reflect.ValueOf(module).Elem()
		moduleType = reflect.TypeOf(module).Elem()
		a.Module = module
	case "validator":
		module := &validatorconfig.Options{}
		moduleValue = reflect.ValueOf(module).Elem()
		moduleType = reflect.TypeOf(module).Elem()
		a.Module = module
	case "relayer":
		module := &relayerconfig.Options{}
		moduleValue = reflect.ValueOf(module).Elem()
		moduleType = reflect.TypeOf(module).Elem()
		a.Module = module
	default:
		return fmt.Errorf("error while decoding YAML: module has invalid type: %s", t)
	}

	for i := 0; i < moduleValue.NumField(); i++ {
		moduleField := moduleValue.Field(i)
		yamlTag, ok := moduleType.Field(i).Tag.Lookup("yaml")
		if !ok {
			continue
		}

		yamlField := strings.Split(yamlTag, ",")
		if len(yamlField) == 0 {
			continue
		}

		if val, found := kvs[yamlField[0]]; found {
			switch newVal := val.(type) {
			case []interface{}:
				newSlice := reflect.MakeSlice(moduleField.Type(), len(newVal), len(newVal))
				for i, listVal := range newVal {
					newSlice.Index(i).Set(reflect.ValueOf(listVal))
				}
				moduleField.Set(newSlice)
			default:
				moduleField.Set(reflect.ValueOf(val))
			}
		}
	}

	return nil
}

var _ yaml.Unmarshaler = &AnyModuleConfig{}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	moduleConfigs := SynapseOptions{}
	globalConfig := cfg.GlobalOptions{}
	err := cfg.LoadFlags(&moduleConfigs, &globalConfig)
	if err != nil {
		logger.Fatal(err)
	}

	utils.CheckNTP()

	lvl, err := logger.ParseLevel(globalConfig.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetLevel(lvl)

	logger.StandardLogger().SetFormatter(&logger.TextFormatter{
		ForceColors: globalConfig.ForceColors,
	})

	changed, newLimit, err := utils.ManageFdLimit()
	if err != nil {
		logger.Fatal(err)
	}
	if changed {
		logger.Infof("changed open file limit to: %d", newLimit)
	}

	beaconConfigs := make([]*beaconconfig.Options, 0, len(moduleConfigs.ModuleConfigs))
	validatorConfigs := make([]*validatorconfig.Options, 0, len(moduleConfigs.ModuleConfigs))
	shardConfigs := make([]*shardconfig.Options, 0, len(moduleConfigs.ModuleConfigs))
	relayerConfigs := make([]*relayerconfig.Options, 0, len(moduleConfigs.ModuleConfigs))

	for _, v := range moduleConfigs.ModuleConfigs {
		switch c := v.Module.(type) {
		case *beaconconfig.Options:
			beaconConfigs = append(beaconConfigs, c)
		case *validatorconfig.Options:
			validatorConfigs = append(validatorConfigs, c)
		case *shardconfig.Options:
			shardConfigs = append(shardConfigs, c)
		case *relayerconfig.Options:
			relayerConfigs = append(relayerConfigs, c)
		}
	}

	// first initialize all of the apps using the configs
	beaconApps := make([]*beaconmodule.BeaconApp, len(beaconConfigs))
	validatorApps := make([]*validatormodule.ValidatorApp, len(validatorConfigs))
	shardApps := make([]*shardmodule.ShardApp, len(shardConfigs))
	relayerApps := make([]*relayermodule.RelayerModule, len(relayerConfigs))

	for i, c := range beaconConfigs {
		app, err := beaconmodule.NewBeaconApp(*c)
		if err != nil {
			logger.Fatal(errors.Wrap(err, "error initializing beacon module"))
		}
		beaconApps[i] = app
	}
	errChan := make(chan error)

	for i, a := range beaconApps {
		go func() {
			logger.Infof("starting beacon module #%d", i)
			errChan <- a.Run()
		}()
	}

	for i, c := range shardConfigs {
		app, err := shardmodule.NewShardApp(*c)
		if err != nil {
			logger.Fatal(errors.Wrap(err, "error initializing shard module"))
		}
		shardApps[i] = app
	}

	for i, a := range shardApps {
		go func() {
			logger.Infof("starting shard module #%d", i)
			errChan <- a.Run()
		}()
	}

	for i, c := range validatorConfigs {
		app, err := validatormodule.NewValidatorApp(*c)
		if err != nil {
			logger.Fatal(errors.Wrap(err, "error initializing validator module"))
		}
		validatorApps[i] = app
	}

	for i, c := range relayerConfigs {
		app, err := relayermodule.NewRelayerModule(*c)
		if err != nil {
			logger.Fatal(errors.Wrap(err, "error initializing relayer module"))
		}
		relayerApps[i] = app
	}

	// order goes: beacon, shard, validator

	for i, a := range validatorApps {
		go func() {
			logger.Infof("starting validator module #%d", i)
			errChan <- a.Run()
		}()
	}

	for i, a := range relayerApps {
		go func() {
			logger.Infof("starting relayer module #%d", i)
			errChan <- a.Run()
		}()
	}

	// if globalConfig.Visualize {
	// 	chain.GlobalVis.Start()
	// }

	for {
		err := <-errChan
		if err != nil {
			logger.Error(err)
			for _, a := range beaconApps {
				a.Exit()
			}
			for _, a := range validatorApps {
				a.Exit()
			}
			for _, a := range shardApps {
				a.Exit()
			}
			continue
		}
		logger.Infof("module exited gracefully")
	}
}
