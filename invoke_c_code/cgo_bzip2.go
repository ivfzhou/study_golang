//go:build linux

// Package bzip2 cgo 调用 C 函数，使用 libbzip2 库。
package bzip2

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -L/usr/lib -lbz2
#include <bzlib.h>
#include <stdlib.h>

bz_stream* bz2alloc() {
	return calloc(1, sizeof(bz_stream));
}

int bz2compress(bz_stream *s, int action, char *in, unsigned *inlen, char *out, unsigned *outlen);

void bzfree(bz_stream* s) {
	free(s);
}
*/
import "C"

import (
	"io"
	"sync"
	"unsafe"
)

type bzWriter struct {
	w      io.Writer
	stream *C.bz_stream
	outbuf [64 * 1024]byte
	lock   sync.Mutex
}

func NewBZip2Writer(out io.Writer) io.WriteCloser {
	const (
		blockSize   = 9
		verbosity   = 0
		workFactory = 30
	)

	w := &bzWriter{w: out, stream: C.bz2alloc()}
	C.BZ2_bzCompressInit(w.stream, blockSize, verbosity, workFactory)

	return w
}

func (bz *bzWriter) Write(data []byte) (int, error) {
	if bz == nil {
		panic("The Writer is nil")
	}

	var total int
	for len(data) > 0 {
		inLen, outLen := C.uint(len(data)), C.uint(cap(bz.outbuf))
		C.bz2compress(bz.stream, C.BZ_RUN, (*C.char)(unsafe.Pointer(&data[0])), &inLen, (*C.char)(unsafe.Pointer(&bz.outbuf)), &outLen)
		total += int(inLen)
		data = data[inLen:]
		if err := func() error {
			bz.lock.Lock()
			defer bz.lock.Unlock()
			_, err := bz.w.Write(bz.outbuf[:outLen])
			return err
		}(); err != nil {
			return total, err
		}
	}

	return total, nil
}

func (bz *bzWriter) Close() error {
	if bz == nil {
		panic("The Writer is nil")
	}

	defer func() {
		C.BZ2_bzCompressEnd(bz.stream)
		C.bzfree(bz.stream)
		bz.stream = nil
	}()

	for {
		inLen, outLen := C.uint(0), C.uint(cap(bz.outbuf))
		r := C.bz2compress(bz.stream, C.BZ_FINISH, nil, &inLen, (*C.char)(unsafe.Pointer(&bz.outbuf)), &outLen)
		if err := func() error {
			bz.lock.Lock()
			defer bz.lock.Unlock()
			_, err := bz.w.Write(bz.outbuf[:outLen])
			return err
		}(); err != nil {
			return err
		}
		if r == C.BZ_STREAM_END {
			return nil
		}
	}
}
