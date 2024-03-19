package device

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

type KeyGeneratorECDSA struct{}

func NewKeyGeneratorECDSA() *KeyGeneratorECDSA {
	return &KeyGeneratorECDSA{}
}

func (kg *KeyGeneratorECDSA) Generate(ctx context.Context) ([]byte, []byte, error) {
	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Encode private key to PEM format
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}
