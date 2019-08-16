// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.package spec

package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"net"
)

const (
	rsaKeySize = 2048
)

// generate a new private key
func newPrivateKey(size int) (*rsa.PrivateKey, error) {
	// TODO: support more key types
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// convert private key to PEM format
func convertKeyToPEM(blockType string, dataBytes *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:    blockType,
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PrivateKey(dataBytes),
		},
	)
}

func NewCSR(commonName string, hostList []string, IPList []string) ([]byte, []byte, error) {
	// TODO: option to use an exist private key
	privKey, err := newPrivateKey(rsaKeySize)
	if err != nil {
		return nil, nil, err
	}

	var ipAddrList []net.IP
	for _, ip := range IPList {
		ipAddr := net.ParseIP(ip)
		ipAddrList = append(ipAddrList, ipAddr)
	}

	// set CSR attributes
	csrTemplate := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization:       []string{"PingCAP"},
			OrganizationalUnit: []string{"TiDB Operator"},
			CommonName:         commonName,
		},
		DNSNames:    hostList,
		IPAddresses: ipAddrList,
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privKey)
	if err != nil {
		return nil, nil, err
	}

	return csr, convertKeyToPEM("RSA PRIVATE KEY", privKey), nil
}
