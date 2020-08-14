package ffi

/*
#cgo LDFLAGS: -lmmr -lsqlite3 -ldl
#include <inttypes.h>

extern int32_t run(int64_t th);
extern char* proof(uint64_t last_leaf, uint64_t member);
*/
import "C"

func RunMMR(th int64) {
	C.run((C.int64_t)(th))
}

func ProofLeaves(last_leaf uint64, member uint64) string {
	return C.GoString(C.proof((C.uint64_t)(last_leaf), (C.uint64_t)(member)))
}
