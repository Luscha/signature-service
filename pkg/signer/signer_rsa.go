package signer

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	// Sign the hash using RSA-PSS with SHA-256
	signature, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, hash[:], nil)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (s *SignerRSA) Verify(ctx context.Context, publicKey, signature, data []byte) error {
	// Decode the PEM-encoded public key
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return errors.New("failed to decode PEM block containing public key")
	}

	// Parse the public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	// Convert the parsed public key to RSA public key type
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("failed to parse RSA public key")
	}

	// Hash the data using SHA-256
	hashed := sha256.Sum256(data)

	// Verify the signature using RSA-PSS with SHA-256
	err = rsa.VerifyPSS(rsaPubKey, crypto.SHA256, hashed[:], signature, nil)
	if err != nil {
		return fmt.Errorf("RSA-PSS signature verification failed: %v", err)
	}

	return nil
}
