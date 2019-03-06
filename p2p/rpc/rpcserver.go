package rpc

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"

	"github.com/phoreproject/synapse/p2p"
	"github.com/phoreproject/synapse/pb"
)

// RPCServer is a server to manager the P2P module
// RPC server.
type RPCServer struct {
	hostNode       *p2p.HostNode
	subscriptions  map[uint64]*pubsub.Subscription
	subChannels    map[uint64]chan []byte
	cancelChannels map[uint64]chan bool
	currentSubID   *uint64 // this is weird, but all of the methods have to pass the struct in by value
	lock           *sync.Mutex
}

// NewRPCServer sets up a server for handling P2P module RPC requests.
func NewRPCServer(hostNode *p2p.HostNode) RPCServer {
	p := RPCServer{
		hostNode:       hostNode,
		subscriptions:  make(map[uint64]*pubsub.Subscription),
		subChannels:    make(map[uint64]chan []byte),
		cancelChannels: make(map[uint64]chan bool),
		currentSubID:   new(uint64),
		lock:           new(sync.Mutex),
	}
	*p.currentSubID = 0
	return p
}

// GetConnectionStatus gets the status of the P2P connection.
func (p RPCServer) GetConnectionStatus(ctx context.Context, in *empty.Empty) (*pb.ConnectionStatus, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	return &pb.ConnectionStatus{Connected: len(p.hostNode.GetLivePeerList()) > 0}, nil
}

// GetPeers gets the peers for the P2P connection.
func (p RPCServer) GetPeers(ctx context.Context, in *empty.Empty) (*pb.GetPeersResponse, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	peers := p.hostNode.GetLivePeerList()
	peersPb := []*pb.Peer{}
	for range peers {
		//peersPb = append(peersPb, &pb.Peer{Address: p.String()})
		// TODO: we need either to see how to get the peer address in handler of SetStreamHandler
		// or to see if we need to return the peer address, we may return more meaningful information
		// such as public keys
		// Now the PeerNode can't know the address if the connection is established in SetStreamHandler
		peersPb = append(peersPb, &pb.Peer{Address: ""})
	}
	return &pb.GetPeersResponse{Peers: peersPb}, nil
}

// ListenForMessages listens to a subscription and receives
// a stream of messages.
func (p RPCServer) ListenForMessages(in *pb.Subscription, out pb.P2PRPC_ListenForMessagesServer) error {
	p.lock.Lock()
	if _, success := p.subscriptions[in.ID]; !success {
		return fmt.Errorf("could not find subscription with ID %d", in.ID)
	}

	messages := p.subChannels[in.ID]
	cancelChan := p.cancelChannels[in.ID]

	p.lock.Unlock()

	for {
		select {
		case msg := <-messages:
			err := out.Send(&pb.Message{Data: msg})
			if err != nil {
				return err
			}
		case <-cancelChan:
			return io.EOF
		}
	}

}

// Subscribe subscribes to a topic returning a subscription ID.
func (p RPCServer) Subscribe(ctx context.Context, in *pb.SubscriptionRequest) (*pb.Subscription, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	subID := *p.currentSubID
	*p.currentSubID++

	subChan := make(chan []byte)
	p.subChannels[subID] = subChan
	p.cancelChannels[subID] = make(chan bool)

	s, err := p.hostNode.SubscribeMessage(in.Topic, func(peer *p2p.PeerNode, data []byte) {
		select {
		case subChan <- data:
		default:
		}
	})

	if err != nil {
		return nil, err
	}

	p.subscriptions[subID] = s

	return &pb.Subscription{ID: subID}, nil
}

// Unsubscribe unsubscribes from a subscription given a subscription ID.
func (p RPCServer) Unsubscribe(ctx context.Context, in *pb.Subscription) (*empty.Empty, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, success := p.subscriptions[in.ID]; !success {
		return nil, fmt.Errorf("could not find subscription with ID %d", in.ID)
	}

	// either send it or not. we don't really care if it works.
	// this is dependent on whether the channel is being listened on
	select {
	case p.cancelChannels[in.ID] <- true:
	default:
	}

	close(p.cancelChannels[in.ID])
	close(p.subChannels[in.ID])
	p.hostNode.UnsubscribeMessage(p.subscriptions[in.ID])

	delete(p.cancelChannels, in.ID)
	delete(p.subChannels, in.ID)
	delete(p.subscriptions, in.ID)

	return &empty.Empty{}, nil
}

// Broadcast broadcasts a message to a topic.
func (p RPCServer) Broadcast(ctx context.Context, in *pb.MessageAndTopic) (*empty.Empty, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	return &empty.Empty{}, p.hostNode.Broadcast(in.Topic, in.Data)
}

// Connect connects to more peers.
func (p RPCServer) Connect(ctx context.Context, in *pb.Peers) (*pb.ConnectResponse, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	success := true
	for _, peer := range in.Peers {
		pInfo, err := p2p.StringToPeerInfo(peer.Address)
		if err != nil {
			return nil, err
		}
		_, err = p.hostNode.Connect(pInfo)
		if err != nil {
			success = false
			logrus.WithField("addr", peer.Address).Warn("could not connect to peer")
			continue
		}
	}
	return &pb.ConnectResponse{Success: success}, nil
}

// GetSettings gets the settings of the P2P connection.
func (p RPCServer) GetSettings(ctx context.Context, in *empty.Empty) (*pb.P2PSettings, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	return &pb.P2PSettings{}, nil
}