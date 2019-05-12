package p2p

import (
	"context"

	host 		"github.com/libp2p/go-libp2p-host"
	libp2p		"github.com/libp2p/go-libp2p"
	pubsub 		"github.com/libp2p/go-libp2p-pubsub"
	hostconf	"github.com/libp2p/go-libp2p/config"
	netCommon	"github.com/AlexandreBelling/go-boojum/network"
)

// PubSub is a wrapper around pubsub to make it compatible with our api
type PubSub struct {
	WhiteList			*WhiteList
	Host				*host.Host
	Engine 				*pubsub.PubSub
	Ctx     			context.Context
}

// NewPubSub constructs a new version of the 
func NewPubSub(ctx context.Context) *PubSub {
	return &PubSub{Ctx: ctx}
}

// DefaultPubSub construct a default PubSub
func DefaultPubSub(ctx context.Context) (*PubSub, error) {
	ps := NewPubSub(ctx)

	err := ps.BuildHost()
	if err != nil {
		return nil, err
	}

	err = ps.BuildEngine()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

// BuildHost reads the configuration and build a host using the previously provided context (by the constructor)
func (ps *PubSub) BuildHost(opts ...hostconf.Option) (error) {

	listenAddress := GetListenAddress()
	identity, err := GetIdentity()
	if err != nil {
		return err
	}

	opts = append(opts, listenAddress, identity)
	hst, err := libp2p.New(ps.Ctx, opts...)
	if err != nil {
		return err
	}

	ps.Host = &hst
	return nil
}

// BuildEngine builds the engine using the previously provided host and context
func (ps *PubSub) BuildEngine(opts ...pubsub.Option) (error) {
	engine, err := pubsub.NewGossipSub(ps.Ctx, *ps.Host, opts...)
	if err != nil {
		return err
	}

	ps.Engine = engine
	return nil
}

// Permissionize prevents connection from non authorized peers
func (ps *PubSub) Permissionize(provider netCommon.WhiteListProvider) {
	ps.WhiteList = NewWhiteList(provider)
	ps.WhiteList.Maintain()

	hst := *(ps.Host)
	hst.Network().SetConnHandler(ps.WhiteList.OnlyWhiteListed)
}


