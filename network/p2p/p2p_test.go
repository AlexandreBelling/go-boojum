package p2p

import (
	"fmt"
	"sync"
	"time"
	"testing"
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/libp2p/go-libp2p-core/peer"

	bnetwork "github.com/AlexandreBelling/go-boojum/network"

)

const ( 
	nServerTest = 20
)

// Test the system
func Test(t *testing.T) {

	log.SetLevel(log.InfoLevel)

	servers := make([]Server, nServerTest)
	wlp := bnetwork.NewMockWhiteListProvider()

	for i := 0; i<nServerTest; i++ {

		addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", 9000 + i)
		s, _ := DefaultServer(addr, wlp)
		servers[i] = *s

		pi := peer.AddrInfo{
			ID:		s.Host.ID(),
			Addrs:	s.Host.Addrs(),
		}

		marshalled, _ := pi.MarshalJSON()
		wlp.Add(marshalled)
	}

	for _, s := range servers {
		s.Start()
	}
	time.Sleep(time.Duration(5) * time.Second)
	
	topicName := "test"

	var wg sync.WaitGroup
	consumeMessage := func(i int, in <- chan []byte) {
		defer wg.Done()
		ctx, can := context.WithTimeout(context.Background(), time.Duration(10) * time.Second)
		defer can()
		
		select {
		case <- ctx.Done():
			t.Fail()
			return
		case msg := <- in:
			log.Infof("%v received from %v", string(msg), i)
			return
		}
	}

	wg.Add(nServerTest)
	for index, s := range servers {
		top, err := s.GetTopic(context.Background(), topicName)
		if err != nil {
			t.Fail()
		}

		go consumeMessage(index, top.Chan())
	}

	servers[0].Publish(topicName, []byte("Hello there"))
	wg.Wait()
}