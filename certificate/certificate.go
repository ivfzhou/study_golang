package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"software.sslmate.com/src/go-pkcs12"
)

/*
openssl genrsa -out ca_root_key.pem -passout pass:123456 8192
openssl req -x509 -new -key ca_root_key.pem -days 36500 -out ca_root_cert.pem -subj '/C=CN/ST=HuNan/L=ChangSha/CN=ivfzhou root certificate'
openssl genrsa -out ca_second_root_key.pem -passout pass:123456 8192
openssl req -new -key ca_second_root_key.pem -out ca_second_root_csr.pem -subj '/C=CN/ST=HuNan/L=ChangSha/CN=ivfzhou second root certificate'
openssl x509 -req -in ca_second_root_csr.pem -CA ca_root_cert.pem -CAkey ca_root_key.pem -days 35600 -out ca_second_root_signed_cert.pem
openssl genrsa -out my_key.pem -passout pass:123456 8192
openssl req -new -key my_key.pem -out my_csr.pem -subj '/C=CN/ST=HuNan/L=ChangSha/CN=ivfzhou certificate'
openssl x509 -req -in my_csr.pem -CA ca_second_root_signed_cert.pem -CAkey ca_second_root_key.pem -days 35600 -out my_cert.pem
openssl pkcs12 -export -inkey my_key.pem -in my_cert.pem -out my.pfx -passout pass:123456
*/

func main() {
	IsTheCertSign()
}

func ParsePfx() {
	bs, err := os.ReadFile(`certificate\my.pfx`)
	if err != nil {
		panic(err)
	}
	key, certificate, err := pkcs12.Decode(bs, "123456")
	if err != nil {
		panic(err)
	}
	fmt.Println(key)
	fmt.Println("NotBefore", certificate.NotBefore) // 2025-01-10 03:55:39 +0000 UTC
	fmt.Println("NotAfter", certificate.NotAfter)   // 2122-07-01 03:55:39 +0000 UTC
	fmt.Println("Issuer", certificate.Issuer)       // CN=ivfzhou second root certificate,L=ChangSha,ST=HuNan,C=CN
	fmt.Println("Subject", certificate.Subject)     // CN=ivfzhou certificate,L=ChangSha,ST=HuNan,C=CN
	fmt.Println("Raw", string(certificate.Raw))
}

func ParseCert() {
	bs, err := os.ReadFile(`certificate\my_cert.pem`)
	if err != nil {
		panic(err)
	}

	bs, err = CertPem2Der(bs)
	if err != nil {
		panic(err)
	}

	certificate, err := x509.ParseCertificate(bs)
	if err != nil {
		panic(err)
	}
	fmt.Println("NotBefore", certificate.NotBefore) // 2025-01-10 03:55:39 +0000 UTC
	fmt.Println("NotAfter", certificate.NotAfter)   // 2122-07-01 03:55:39 +0000 UTC
	fmt.Println("Issuer", certificate.Issuer)       // CN=ivfzhou second root certificate,L=ChangSha,ST=HuNan,C=CN
	fmt.Println("Subject", certificate.Subject)     // CN=ivfzhou certificate,L=ChangSha,ST=HuNan,C=CN
	fmt.Println("Raw", string(certificate.Raw))
}

func ParseCertChain() {
	bs, err := os.ReadFile(`certificate\my_cert_chain.pem`)
	if err != nil {
		panic(err)
	}

	lastCert, err := GetLastCert(bs)
	if err != nil {
		panic(err)
	}
	fmt.Println("NotBefore", lastCert.NotBefore) // 2025-01-10 03:55:39 +0000 UTC
	fmt.Println("NotAfter", lastCert.NotAfter)   // 2122-07-01 03:55:39 +0000 UTC
	fmt.Println("Issuer", lastCert.Issuer)       // CN=ivfzhou second root certificate,L=ChangSha,ST=HuNan,C=CN
	fmt.Println("Subject", lastCert.Subject)     // CN=ivfzhou certificate,L=ChangSha,ST=HuNan,C=CN
	fmt.Println("Raw", string(lastCert.Raw))
}

func EqualCert() {
	bs, err := os.ReadFile(`certificate\my.pfx`)
	if err != nil {
		panic(err)
	}
	_, certificate, err := pkcs12.Decode(bs, "123456")
	if err != nil {
		panic(err)
	}

	bs, err = os.ReadFile(`certificate\my_cert.pem`)
	if err != nil {
		panic(err)
	}

	c1, err := CertDer2Pem(certificate.Raw)
	if err != nil {
		panic(err)
	}

	b, err := EqualCertPublicKey(c1, bs)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}

func IsTheCertSign() {
	signer, err := os.ReadFile(`certificate\ca_second_root_signed_cert.pem`)
	if err != nil {
		panic(err)
	}
	signer, err = CertPem2Der(signer)
	if err != nil {
		panic(err)
	}
	sginerCert, err := x509.ParseCertificate(signer)
	if err != nil {
		panic(err)
	}

	signed, err := os.ReadFile(`certificate\my_cert.pem`)
	if err != nil {
		panic(err)
	}
	signed, err = CertPem2Der(signed)
	if err != nil {
		panic(err)
	}
	signedCert, err := x509.ParseCertificate(signed)
	if err != nil {
		panic(err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(sginerCert)
	opts := x509.VerifyOptions{
		// DNSName: signedCert.DNSNames[0],
		Roots: roots,
	}

	if _, err = signedCert.Verify(opts); err != nil {
		fmt.Println("Failed to verify the certificate", err)
	} else {
		fmt.Println("The certificate is issued by the CA")
	}
}

func EqualCertPublicKey(c1pem, c2pem []byte) (bool, error) {
	cert1, err := GetLastCert(c1pem)
	if err != nil {
		return false, err
	}
	cert2, err := GetLastCert(c2pem)
	if err != nil {
		return false, err
	}
	pub1, err := publicKeyToString(cert1.PublicKey)
	if err != nil {
		return false, err
	}
	pub2, err := publicKeyToString(cert2.PublicKey)
	if err != nil {
		return false, err
	}
	return pub1 == pub2, nil
}

func publicKeyToString(pub any) (string, error) {
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return base64.StdEncoding.EncodeToString(pub.N.Bytes()), nil
	case *ecdsa.PublicKey:
		return base64.StdEncoding.EncodeToString(append(pub.X.Bytes(), pub.Y.Bytes()...)), nil
	default:
		return "", fmt.Errorf("unsupported public key type")
	}
}

func CertPem2Der(src []byte) ([]byte, error) {
	block, rest := pem.Decode(src)
	if len(rest) > 0 {
		return nil, fmt.Errorf("has reset")
	}
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("not found cert")
	}
	return block.Bytes, nil
}

func CertDer2Pem(src []byte) ([]byte, error) {
	p := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: src,
	}
	buf := &bytes.Buffer{}
	err := pem.Encode(buf, p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func KeyDer2Pem(src []byte) ([]byte, error) {
	p := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: src,
	}
	buf := &bytes.Buffer{}
	err := pem.Encode(buf, p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func KeyPem2Der(src, passwd []byte) ([]byte, error) {
	block, rest := pem.Decode(src)
	if len(rest) > 0 {
		return nil, fmt.Errorf("has reset")
	}
	if block == nil {
		return nil, fmt.Errorf("not found key")
	}
	if block.Type == "PRIVATE KEY" {
		return block.Bytes, nil
	}
	if block.Type == "ENCRYPTED PRIVATE KEY" {
		pemBlock, err := x509.DecryptPEMBlock(block, passwd)
		if err != nil {
			return nil, err
		}
		return pemBlock, nil
	}
	return nil, fmt.Errorf("not found key")
}

func CsrPem2Der(src []byte) ([]byte, error) {
	block, rest := pem.Decode(src)
	if len(rest) > 0 {
		return nil, fmt.Errorf("has reset")
	}
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("not found cert req")
	}
	return block.Bytes, nil
}

func CsrDer2Pem(src []byte) ([]byte, error) {
	p := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: src,
	}
	buf := &bytes.Buffer{}
	err := pem.Encode(buf, p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetLastCert(crtPem []byte) (*x509.Certificate, error) {
	var certs []*x509.Certificate
	var block *pem.Block
	for {
		block, crtPem = pem.Decode(crtPem)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			certs = append(certs, cert)
		}
	}
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certs found")
	}
	return certs[len(certs)-1], nil
}
