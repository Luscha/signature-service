package workspace

import (
	"context"
	"fmt"

	"signature.service/pkg/logger"
	"signature.service/pkg/signature"
	"signature.service/pkg/signer"
)

var (
	ErrDeviceNotFound = fmt.Errorf("device not found")
)

func (w *Workspace) Sign(ctx context.Context, deviceId string, data []byte) (*signature.Signature, error) {
	device, err := w.storage.GetDevice(ctx, deviceId)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	signer, err := signer.NewSigner(device)
	if err != nil {
		logger.GetLogger(ctx).WithError(err).Error("could not instantiate signer")
		return nil, err
	}

	signature, err := signer.Sign(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).WithError(err).Error("could not sign data")
		return nil, err
	}

	err = w.storage.PostSignature(ctx, signature)
	if err != nil {
		logger.GetLogger(ctx).WithError(err).Error("could not save signature")
		return nil, err
	}

	return signature, nil
}

func (w *Workspace) ListSignatures(ctx context.Context) ([]*signature.Signature, error) {
	signatures, err := w.storage.ListSignature(ctx)
	if err != nil {
		return nil, err
	}

	return signatures, nil
}
