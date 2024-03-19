package workspace

import (
	"context"

	"signature.service/pkg/device"
)

func (w *Workspace) CreateDevice(ctx context.Context, algorithm device.Algorithm, label string) (*device.Device, error) {
	device, err := device.NewDevice(ctx, algorithm, label)
	if err != nil {
		return nil, err
	}

	err = w.storage.PostDevice(ctx, device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (w *Workspace) ListDevices(ctx context.Context) ([]*device.Device, error) {
	devices, err := w.storage.ListDevices(ctx)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (w *Workspace) GetDevice(ctx context.Context, id string) (*device.Device, error) {
	device, err := w.storage.GetDevice(ctx, id)
	if err != nil {
		return nil, err
	}

	return device, nil
}
