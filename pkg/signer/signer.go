package signer

import (
	"context"
	"encoding/base64"
	"fmt"

	device_pkg "signature.service/pkg/device"
	"signature.service/pkg/signature"
)

type SignerImpl interface {
	Sign(ctx context.Context, privateKey []byte, data []byte) ([]byte, error)
}

type Signer struct {
	device *device_pkg.Device
	impl   SignerImpl
}

func NewSigner(device *device_pkg.Device) (*Signer, error) {
	var impl SignerImpl
	switch device.Algorithm {
	case device_pkg.ECDSA:
		{
			impl = NewSignerECDSA()
		}
	case device_pkg.RSA:
		{
			impl = NewSignerRSA()
		}
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", device.Algorithm)
	}

	return &Signer{
		device: device,
		impl:   impl,
	}, nil
}

func (s *Signer) Sign(ctx context.Context, data []byte) (*signature.Signature, error) {
	s.device.Acquire()
	defer s.device.Release()

	safeDataToBeSigned := s.generateSecuredDataToBeSigned(ctx, data)
	signedData, err := s.impl.Sign(ctx, []byte(s.device.PrivateKey), safeDataToBeSigned)

	if err != nil {
		return nil, err
	}

	base64SignedData := base64.StdEncoding.EncodeToString(signedData)
	s.device.IncCounter()
	s.device.SetLastSignature(base64SignedData)
	return &signature.Signature{
		Device:     s.device.Id,
		Data:       data,
		SignedData: safeDataToBeSigned,
		Signature:  base64SignedData,
	}, nil
}

func (s *Signer) generateSecuredDataToBeSigned(ctx context.Context, data []byte) []byte {
	var lastSignature string
	if s.device.SignatureCounter == 0 {
		// Use base64-encoded device ID as the last signature
		lastSignature = base64.StdEncoding.EncodeToString([]byte(s.device.Id))
	} else {
		lastSignature = s.device.LastSignature
	}

	// Construct secured_data_to_be_signed string
	securedDataToBeSigned := fmt.Sprintf("%s_%s_%s", fmt.Sprintf("%d", s.device.SignatureCounter), data, lastSignature)

	return []byte(securedDataToBeSigned)
}
