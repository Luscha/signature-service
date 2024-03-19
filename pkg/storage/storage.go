package storage

import (
	"context"

	"signature.service/pkg/device"
	"signature.service/pkg/signature"
)

type Storage interface {
	PostDevice(ctx context.Context, device *device.Device) error
	PutDevice(ctx context.Context, device *device.Device) error
	GetDevice(ctx context.Context, id string) (*device.Device, error)
	ListDevices(ctx context.Context) ([]*device.Device, error)
	PostSignature(ctx context.Context, sign *signature.Signature) error
	ListSignature(ctx context.Context) ([]*signature.Signature, error)
}
