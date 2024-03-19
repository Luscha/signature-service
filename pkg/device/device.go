package device

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"signature.service/pkg/logger"
)

var mu sync.Mutex

type Algorithm string

const (
	ECDSA Algorithm = "ECDSA"
	RSA   Algorithm = "RSA"
)

type Device struct {
	Id               string    `json:"id"`
	Label            string    `json:"label"`
	SignatureCounter uint64    `json:"signature_counter"`
	LastSignature    string    `json:"last_signature"`
	Algorithm        Algorithm `json:"algorithm"`
	PublicKey        string    `json:"public_key"`
	PrivateKey       string    `json:"-"`
}

func (d *Device) Acquire() {
	mu.Lock()
}

func (d *Device) Release() {
	mu.Unlock()
}

func (d *Device) IncCounter() {
	d.SignatureCounter++
}

func (d *Device) SetLastSignature(base64Signature string) {
	d.LastSignature = base64Signature
}

func NewDevice(ctx context.Context, algorithm Algorithm, label string) (*Device, error) {
	var keyGen KeyGenerator
	switch algorithm {
	case ECDSA:
		{
			keyGen = NewKeyGeneratorECDSA()
		}
	case RSA:
		{
			keyGen = NewKeyGeneratorRSA()
		}
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algorithm)
	}

	privateKey, publicKey, err := keyGen.Generate(ctx)
	if err != nil {
		logger.GetLogger(ctx).WithError(err).Error("could not generate key pair")
		return nil, fmt.Errorf("could not generate key pair: %s", err.Error())
	}

	return &Device{
		Id:               uuid.New().String(),
		Algorithm:        algorithm,
		Label:            label,
		PublicKey:        string(publicKey),
		PrivateKey:       string(privateKey),
		SignatureCounter: 0,
		LastSignature:    "",
	}, nil
}
