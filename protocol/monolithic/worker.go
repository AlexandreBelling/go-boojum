package monolithic

import(
	"github.com/AlexandreBelling/go-boojum/aggregator"
)

// Worker is responsible to perform aggregation
type Worker struct {
	boo *aggregator.Boojum
}

// StartConsuming in the pendings queue until it receives done
func (w *Worker) StartConsuming(pendings chan Tree, done chan bool) {
	for {
		select {
		case <- done:
			return
		case job := <- pendings:

			// Scheduler ensures that whenever job is received
			// left and right are already assigned
			job.payloadChan <- w.boo.AggregateTrees(
				job.left.payload,
				job.right.payload,
			)	
		}
	}
}