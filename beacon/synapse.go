package main

import (
	"crypto/rand"
	"flag"
	"fmt"

	"github.com/phoreproject/synapse/rpc"

	"github.com/libp2p/go-libp2p-crypto"

	"github.com/libp2p/go-libp2p-peerstore"

	"github.com/multiformats/go-multiaddr"

	"github.com/phoreproject/synapse/blockchain"
	"github.com/phoreproject/synapse/db"
	"github.com/phoreproject/synapse/net"

	logger "github.com/inconshreveable/log15"
	iaddr "github.com/ipfs/go-ipfs-addr"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

func parseInitialConnections(in string) ([]*peerstore.PeerInfo, error) {
	currentAddr := ""

	peers := []*peerstore.PeerInfo{}

	for i := range in {
		if in[i] == ',' {
			addr, err := iaddr.ParseString(currentAddr)
			currentAddr = ""
			if err != nil {
				return nil, err
			}
			peerinfo, err := pstore.InfoFromP2pAddr(addr.Multiaddr())
			if err != nil {
				return nil, err
			}

			peers = append(peers, peerinfo)
		}
		currentAddr = currentAddr + string(in[i])
	}

	return peers, nil
}

const clientVersion = "0.0.1"

func main() {
	listen := flag.String("listen", "/ip4/0.0.0.0/tcp/11781", "specifies the address to listen on")
	initialConnections := flag.String("connect", "", "comma separated multiaddrs")
	rpcConnect := flag.String("rpclisten", "127.0.0.1:11782", "host and port for RPC server to listen on")
	flag.Parse()

	logger.Info("initializing client", "version", clientVersion)

	logger.Info("initializing database")
	database := db.NewInMemoryDB()
	c := blockchain.MainNetConfig

	logger.Info("initializing blockchain")
	blockchain := blockchain.NewBlockchain(database, &c)

	logger.Info("initializing net")
	ps, err := parseInitialConnections(*initialConnections)
	if err != nil {
		panic(err)
	}

	sourceMultiAddr, err := multiaddr.NewMultiaddr(*listen)
	if err != nil {
		fmt.Printf("address %s is invalid", *listen)
		return
	}

	priv, _, err := crypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		panic(err)
	}

	network, err := net.NewNetworkingService(&sourceMultiAddr, priv)
	if err != nil {
		panic(err)
	}

	logger.Info("connecting to bootnodes")

	for _, p := range ps {
		err = network.Connect(p)
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}

	logger.Info("listening for blocks")

	blocks := network.GetBlocksChannel()

	go blockchain.HandleNewBlocks(blocks)

	logger.Info("initializing RPC")

	err = rpc.Serve(*rpcConnect, &blockchain)
	if err != nil {
		panic(err)
	}
}