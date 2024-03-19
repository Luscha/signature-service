package signer

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"math/big"
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

	// Encode the signature in ASN.1 DER format
	signatureBytes, err := asn1.Marshal(struct {
		R, S *big.Int
	}{r, s_})
	if err != nil {
		return nil, err
	}

	return signatureBytes, nil
}

func (s *SignerECDSA) Verify(ctx context.Context, publicKey, signature, data []byte) error {
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

	// Convert the parsed public key to ECDSA public key type
	ecPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to parse ECDSA public key")
	}

	// Parse the ASN.1 encoded signature
	var ecdsaSignature struct {
		R, S *big.Int
	}
	_, err = asn1.Unmarshal(signature, &ecdsaSignature)
	if err != nil {
		return err
	}

	// Compute the hash of the data
	hash := sha256.Sum256(data)

	// Verify the signature
	if !ecdsa.Verify(ecPubKey, hash[:], ecdsaSignature.R, ecdsaSignature.S) {
		return errors.New("ECDSA signature verification failed")
	}

	return nil
}
