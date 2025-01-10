package main

import (
	"bytes"
	"net"
)

type interceptorReader struct {
	net.Conn
	bs bytes.Buffer
}

func (r *interceptorReader) Read(p []byte) (n int, err error) {
	defer func() {
		r.bs.Write(p)
	}()
	return r.Conn.Read(p)
}

func (r *interceptorReader) Bytes() []byte {
	return r.bs.Bytes()
}

type interceptorWriter struct {
	net.Conn
	bs bytes.Buffer
}

func (w *interceptorWriter) Write(p []byte) (n int, err error) {
	w.bs.Write(p)
	return w.Conn.Write(p)
}

func (w *interceptorWriter) Bytes() []byte {
	return w.bs.Bytes()
}
