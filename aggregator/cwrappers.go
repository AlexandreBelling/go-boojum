package aggregator

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
func makeExampleProof(treeBuffer **byte) {

	freeMeAfter := unsafe.Pointer(*treeBuffer)
	defer C.free(freeMeAfter)

	var treePtr *byte
	treePtr = nil
	treePtrC := (*unsafe.Pointer)(unsafe.Pointer(&treePtr))
	
	C.make_example_proof(
		treePtrC,
	)

	*treeBuffer = treePtr
}

func proveAggregation(
	leftBuffer *byte,
	rightBuffer *byte,
	outputBuffer **byte,
) {

	freeMeAfter := unsafe.Pointer(*outputBuffer)
	defer C.free(freeMeAfter)

	var outputPtr *byte
	outputPtr = nil
	outputPtrC := (*unsafe.Pointer)(unsafe.Pointer(&outputPtr))

	C.prove_aggregation(
		unsafe.Pointer(&(*leftBuffer)),
		unsafe.Pointer(&(*rightBuffer)),
		(*unsafe.Pointer)(outputPtrC),
	)

	*outputBuffer = outputPtr
}

func verify(treeBuffer *byte) (bool) {
	valid := C.verify(
		unsafe.Pointer(&(*treeBuffer)),
	)
	return bool(valid)
}

func memFree(treeBuffer *byte) {
	// No need to call the internal memfree function for that
	C.free(unsafe.Pointer(treeBuffer))
}
