package blockchain

// Client is an interface to Ethereum
type Client interface {
	Sender
	Listener
}

// Sender is an abstract interface
// for object able to send transaction on a blockchain
//
// It is also required to contain the signature logic
type Sender interface {
	SendTransaction(Transaction) error
}

// Listener is an abstract interface
// for object able to forward a stream of mined transaction
type Listener interface {
	// Notify will add a Notification in the list
	// Every notification will be triggered
	AddNotification(Notification)
}

// A Notification is a function that can be triggered
//
// A notification is typically very quick to execute. If you want to apply
// complex and long processing for listened transaction. We strongly
// that you leave those processing to another go-routine.
type Notification func(Transaction)

// Transaction is a generic interface carried by a blockchain client
type Transaction interface {
	// General purpose getter, they have a sense
	// whatever the blockchain
	GetTo() string
	GetData() []byte
}
