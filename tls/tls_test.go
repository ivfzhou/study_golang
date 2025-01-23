package tls_test

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"os"
	"testing"
)

/*
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=IvfzhouRootCA" -days 365 -out ca.crt

openssl genrsa -out server.key 2048
openssl req -new -key server.key -config csr.ini -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -extensions v3_ext -extfile csr.ini -out server.crt

openssl genrsa -out client.key 2048
openssl req -new -key client.key -config csr.ini -out client.csr
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -extensions v3_ext -extfile csr.ini -out client.crt
openssl pkcs12 -export -in client.crt -inkey client.key -out client.pfx
cat client.key client.crt > client.pem2

解析 pfx
openssl pkcs12 -in .pfx -clcert -nokeys -out .crt
openssl pkcs12 -in .pfx -nocerts -nodes -out .pem
*/

func TestServer(t *testing.T) {
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		t.Fatal(err)
	}
	rootCAPool := x509.NewCertPool()
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	rootCAPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAPool,
	}
	listen, err := tls.Listen("tcp", ":8080", config)
	if err != nil {
		t.Fatal(err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			t.Error(err)
			continue
		}
		go func() {
			defer conn.Close()
			reader := bufio.NewReader(conn)
			str, err := reader.ReadString('\n')
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(str)
			_, err = conn.Write([]byte("hello\n"))
			if err != nil {
				t.Error(err)
			}
		}()
	}
}

func TestClient(t *testing.T) {
	certificate, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		t.Fatal(err)
	}
	rootCAPool := x509.NewCertPool()
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		t.Fatal(err)
	}
	rootCAPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   "127.0.0.1",
		RootCAs:      rootCAPool,
	}
	conn, err := tls.Dial("tcp", "127.0.0.1:8080", config)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("hello\n"))
	if err != nil {
		t.Error(err)
	}
	buf := make([]byte, 100)
	n, err := conn.Read(buf)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(buf[:n]))
}
