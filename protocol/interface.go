package protocol

import (
	"github.com/AlexandreBelling/go-boojum/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type MemberProvider interface {
	GetMembers() []ID
}

type DefaultMembersProvider struct {
	WLP 	network.WhiteListProvider
}

// GetMembers return the list of all the members
func (d *DefaultMembersProvider) GetMembers() []ID {
	peers, _ := d.WLP.GetPeers()
	members := make([]ID, len(peers))
	for index, peer := range peers {
		copy(members[index][:], peer)
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
