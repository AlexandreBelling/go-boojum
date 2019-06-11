package p2p

import (
	"sync"
	"time"
	"testing"
	"context"

	log "github.com/sirupsen/logrus"
)

const ( 
	nServerTest = 20
)

// Test the system
func Test(t *testing.T) {

	log.SetLevel(log.InfoLevel)
	servers := MakeServers(nServerTest)
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
		top := s.GetTopic(context.Background(), topicName)
		channel, err := top.Chan()
		if err != nil {
			t.Fail()
		}
		go consumeMessage(index, channel)
	}

	servers[0].Publish(topicName, []byte("Hello there"))
	wg.Wait()
}