package layered

import (
	net "github.com/AlexandreBelling/go-boojum/network"
)

// Participant ..
type Participant struct {
	Network 		net.Network
	NetworkIn 		chan []byte
}