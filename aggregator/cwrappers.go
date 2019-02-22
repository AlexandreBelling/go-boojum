package boojum

// #cgo CFLAGS: -I${SRCDIR}/c-boojum/src -I${SRCDIR}/c-boojum/depends
// #cgo LDFLAGS: -L${SRCDIR}/compiled -lboojum -lstdc++ -lff -lgomp -lsnark -lgmp -lprocps -lm -lcrypto -lgmpxx
// #include <./c-boojum/src/libboojum.h>
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
	treePtr := (*unsafe.Pointer)(unsafe.Pointer(&(**treeBuffer)[0]))
	C.make_example_proof(
		treePtr,
	)
}

func proveAggregation(
	leftBuffer *[]byte,
	rightBuffer *[]byte,
	outputBuffer **[]byte,
) {
	outputPtr := (*unsafe.Pointer)(unsafe.Pointer(&(**outputBuffer)[0]))

	C.prove_aggregation(
		unsafe.Pointer(&(*leftBuffer)[0]),
		unsafe.Pointer(&(*rightBuffer)[0]),
		outputPtr,
	)
}

func verify(treeBuffer *[]byte) (bool) {
	valid := C.verify(
		unsafe.Pointer(&(*treeBuffer)[0]),
	)
	return bool(valid)
}
