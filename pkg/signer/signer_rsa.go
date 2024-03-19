package signer

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type SignerRSA struct{}

func NewSignerRSA() *SignerRSA {
	return &SignerRSA{}
}

func (s *SignerRSA) Sign(ctx context.Context, privateKey []byte, data []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}
