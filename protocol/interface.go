package protocol

import (
	// log "github.com/sirupsen/logrus"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/AlexandreBelling/go-boojum/network"
	"github.com/AlexandreBelling/go-boojum/identity"
)

// MemberProvider is a helper for getting co-peers in permissioned network
type MemberProvider interface {
	GetMembers() []identity.ID
}

// DefaultMembersProvider is an implementation of MembersProvider
type DefaultMembersProvider struct {
	WLP network.WhiteListProvider
}

// GetMembers return the list of all the members
func (d *DefaultMembersProvider) GetMembers() []identity.ID {
	pis, _ := d.WLP.GetPeers()
	members := make([]identity.ID, len(pis))
	for index, b := range pis {
		pi := peer.AddrInfo{}
		_ = pi.UnmarshalJSON(b)
		members[index] = identity.ID(pi.ID)
	}
	return members
}

