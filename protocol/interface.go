package protocol

import (
	// log "github.com/sirupsen/logrus"
	"github.com/AlexandreBelling/go-boojum/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// MemberProvider is a helper for getting co-peers in permissioned network
type MemberProvider interface {
	GetMembers() []ID
}

// DefaultMembersProvider is an implementation of MembersProvider
type DefaultMembersProvider struct {
	WLP network.WhiteListProvider
}

// GetMembers return the list of all the members
func (d *DefaultMembersProvider) GetMembers() []ID {
	pis, _ := d.WLP.GetPeers()
	members := make([]ID, len(pis))
	for index, pi := range pis {
		piUnmarshalled := &peer.AddrInfo{}
		_ = piUnmarshalled.UnmarshalJSON(pi)
		copy(members[index][:32], piUnmarshalled.ID)
	}
	return members
}

// PeerIDtoProtocolID converts a libp2p peer.ID to an ID
func PeerIDtoProtocolID(id peer.ID) ID {
	var res ID
	idByte, _ := id.Marshal()
	copy(res[:], idByte)
	return res
}
