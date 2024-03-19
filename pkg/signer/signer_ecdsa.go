package signer

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type SignerECDSA struct{}

func NewSignerECDSA() *SignerECDSA {
	return &SignerECDSA{}
}

func (s *SignerECDSA) Sign(ctx context.Context, privateKey []byte, data []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Compute the hash of the data
	hash := sha256.Sum256(data)
	// Sign the hash using the ECDSA private key
	r, s_, err := ecdsa.Sign(rand.Reader, key, hash[:])
	if err != nil {
		return nil, err
	}

	// Concatenate R and S into a single byte slice
	rBytes, sBytes := r.Bytes(), s_.Bytes()
	signature := make([]byte, 64)
	copy(signature[32-len(rBytes):32], rBytes)
	copy(signature[64-len(sBytes):], sBytes)

	return signature, nil
}
