package p2p

import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/inconshreveable/log15"
	"github.com/libp2p/go-libp2p-pubsub"

	"github.com/phoreproject/synapse/net"
	"github.com/phoreproject/synapse/pb"
)

type p2prpcServer struct {
	service        *net.NetworkingService
	subscriptions  map[uint64]*pubsub.Subscription
	subChannels    map[uint64]chan []byte
	cancelChannels map[uint64]chan bool
	currentSubID   uint64
}

// NewP2PRPCServer sets up a server for handling P2P module RPC requests.
func NewP2PRPCServer(netService *net.NetworkingService) p2prpcServer {
	return p2prpcServer{
		service:        netService,
		subscriptions:  make(map[uint64]*pubsub.Subscription),
		subChannels:    make(map[uint64]chan []byte),
		cancelChannels: make(map[uint64]chan bool),
		currentSubID:   0,
	}
}

func (p p2prpcServer) GetConnectionStatus(ctx context.Context, in *empty.Empty) (*pb.ConnectionStatus, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p p2prpcServer) GetPeers(ctx context.Context, in *empty.Empty) (*pb.GetPeersResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p p2prpcServer) ListenForMessages(in *pb.Subscription, out pb.P2PRPC_ListenForMessagesServer) error {
	if _, success := p.subscriptions[in.ID]; !success {
		return fmt.Errorf("could not find subscription with ID %d", in.ID)
	}

	log15.Debug("listening to new messages on sub", "subID", in.ID)

	for {
		select {
		case msg := <-p.subChannels[in.ID]:
			err := out.Send(&pb.Message{Data: msg})
			if err != nil {
				return err
			}
		case <-p.cancelChannels[in.ID]:
			return io.EOF
		}
	}

}

func (p p2prpcServer) Subscribe(ctx context.Context, in *pb.SubscriptionRequest) (*pb.Subscription, error) {
	subID := p.currentSubID
	p.currentSubID++

	log15.Debug("subscribed to new messages", "topic", in.Topic, "subID", subID)

	p.subChannels[subID] = make(chan []byte)
	p.cancelChannels[subID] = make(chan bool)

	s, err := p.service.RegisterHandler(in.Topic, func(b []byte) error {
		select {
		case p.subChannels[subID] <- b:
		default:
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	p.subscriptions[subID] = s

	return &pb.Subscription{ID: subID}, nil
}

func (p p2prpcServer) Unsubscribe(ctx context.Context, in *pb.Subscription) (*empty.Empty, error) {
	if _, success := p.subscriptions[in.ID]; !success {
		return nil, fmt.Errorf("could not find subscription with ID %d", in.ID)
	}

	log15.Debug("unsubscribed to subID", "subID", in.ID)

	// either send it or not. we don't really care if it works.
	// this is dependent on whether the channel is being listened on
	select {
	case p.cancelChannels[in.ID] <- true:
	default:
	}

	close(p.cancelChannels[in.ID])
	close(p.subChannels[in.ID])
	p.subscriptions[in.ID].Cancel()

	return &empty.Empty{}, nil
}

func (p p2prpcServer) Broadcast(ctx context.Context, in *pb.MessageAndTopic) (*empty.Empty, error) {
	return &empty.Empty{}, p.service.Broadcast(in.Topic, in.Data)
}

func (p p2prpcServer) Connect(ctx context.Context, in *pb.Peers) (*pb.ConnectResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p p2prpcServer) Disconnect(ctx context.Context, in *empty.Empty) (*pb.DisconnectResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p p2prpcServer) GetSettings(ctx context.Context, in *empty.Empty) (*pb.P2PSettings, error) {
	return nil, fmt.Errorf("not implemented")
}