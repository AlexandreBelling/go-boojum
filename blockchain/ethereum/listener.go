package ethereum

import (
	"context"
	"github.com/AlexandreBelling/go-boojum/blockchain"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"sync"
	"time"
)

// Listener is an implementation of blockchain.Listener
type Listener struct {
	ctx    context.Context
	client *ethclient.Client

	blockNumber   uint64
	refreshPeriod time.Duration

	Notifications []blockchain.Notification
	notifMut      sync.Mutex
}

// AddNotification registers a new notification in the notification list
func (l *Listener) AddNotification(notif blockchain.Notification) {
	l.notifMut.Lock()
	defer l.notifMut.Unlock()
	l.Notifications = append(l.Notifications, notif)
}

// Start runs the main routine of the Listener
func (l *Listener) Start(fromBlockNumber uint64) error {
	return nil
}

func (l *Listener) blockFetchingLoop() {
	ticker := time.NewTicker(l.refreshPeriod)

tickerLoop:
	for {
		select {
		case <-ticker.C:

			newBlockNumber, err := l.GetBlockNumber()
			if err != nil {
				// If client cannot be reached we resort to wait for the next tick
				continue tickerLoop
			}

			// If l.blockNumer == newBlocknumber this loop has no iteration
			for n := l.blockNumber + 1; n <= newBlockNumber; n++ {
				block, err := l.GetBlock(n)
				if err != nil {
					continue tickerLoop
				}

				l.ReadThroughNotify(block)
				l.blockNumber = n
			}

		case <-l.ctx.Done():
			ticker.Stop()
		}
	}
}

// GetBlockNumber returns the latest value of the blocknumber
//
// This is the value you want to use to fetch the latest known block
func (l *Listener) GetBlockNumber() (uint64, error) {
	headers, err := l.client.HeaderByNumber(l.ctx, nil)
	if err != nil {
		return 0, err
	}
	return headers.Number.Uint64(), nil
}

// GetBlock returns the block at given height
func (l *Listener) GetBlock(n uint64) (*types.Block, error) {
	var nBigInt *big.Int
	return l.client.BlockByNumber(l.ctx, nBigInt.SetUint64(n))
}

// ReadThroughNotify applies the Notifications on all the transaction of a given block
func (l *Listener) ReadThroughNotify(block *types.Block) {
	for _, ethtx := range block.Transactions() {
		tx := (*Transaction)(ethtx)
		l.notifMut.Lock()
		for _, notif := range l.Notifications {
			notif(tx)
		}
		l.notifMut.Unlock()
	}
}
