package boojum

// #cgo CFLAGS: -I${SRCDIR}/src -I${SRCDIR}/depends
// #cgo LDFLAGS: -L${SRCDIR}/objects -lsnark-aggregation -lstdc++ -lff -lgomp -lsnark -lgmp -lprocps -lm -lcrypto -lgmpxx
// #include <./src/libboojum.h>
import "C"
import (
	"sync"
	"unsafe"
)

var onceInit sync.Once

func initialize() {
	onceInit.Do(func() {
		C.initialize()
	})
}

func runGenerators(dir string) {
	C.run_generators(C.CString(dir))
}

// Assign an example_tree to 
func makeExampleProof(treeBuffer **[]byte) {
	C.make_example_proof(
		&unsafe.Pointer(&(**tree_buffer[0])),
	)
}

func proveAggregation(
	leftProof *[]byte,
	rightProof *[]byte,
	outputProof **[]byte,
) {
	C.prove_aggregation(
		unsafe.Pointer(&(*leftProof[0])),
		unsafe.Pointer(&(*rightProof[0])),
		&unsafe.Pointer(&(**outputProof[0])),
	)
}

func verify(treeBuff *[]byte) bool {
	return C.verify(
		unsafe.Pointer(&(*treeBuff[0])),
	)
}
