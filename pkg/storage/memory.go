package storage

import (
	"context"
	"fmt"

	"signature.service/pkg/device"
	"signature.service/pkg/signature"
)

type MemoryStorage struct {
	devices    []*device.Device
	signatures []*signature.Signature
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		devices:    make([]*device.Device, 0),
		signatures: make([]*signature.Signature, 0),
	}
}

func (m *MemoryStorage) PostDevice(ctx context.Context, device *device.Device) error {
	m.devices = append(m.devices, device)
	return nil
}

func (m *MemoryStorage) PutDevice(ctx context.Context, device *device.Device) error {
	for _, d := range m.devices {
		if d.Id == device.Id {
			d.LastSignature = device.LastSignature
			d.SignatureCounter = device.SignatureCounter
			return nil
		}
	}

	return fmt.Errorf("device not found")
}

func (m *MemoryStorage) GetDevice(ctx context.Context, id string) (*device.Device, error) {
	for _, d := range m.devices {
		if d.Id == id {
			return d, nil
		}
	}

	return nil, fmt.Errorf("device not found")
}

func (m *MemoryStorage) ListDevices(ctx context.Context) ([]*device.Device, error) {
	return m.devices, nil
}

func (m *MemoryStorage) PostSignature(ctx context.Context, sign *signature.Signature) error {
	m.signatures = append(m.signatures, sign)
	return nil
}

func (m *MemoryStorage) ListSignature(ctx context.Context) ([]*signature.Signature, error) {
	return m.signatures, nil
}
