package device

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeviceECDSA(t *testing.T) {
	ctx := context.TODO()
	label := "TestDevice"
	algorithm := Algorithm("ECDSA")

	device, err := NewDevice(ctx, algorithm, label)

	assert.NoError(t, err, "Unexpected error for ECDSA algorithm")
	assert.NotNil(t, device, "Device should not be nil for ECDSA algorithm")
	assert.Equal(t, uint64(0), device.SignatureCounter, "Expected SignatureCounter to be 0 for ECDSA algorithm")
	assert.Equal(t, label, device.Label, "Unexpected Label for ECDSA algorithm")
	assert.NotEmpty(t, device.PublicKey, "Public key should not be empty for ECDSA algorithm")
	assert.NotEmpty(t, device.PrivateKey, "Private key should not be empty for ECDSA algorithm")
	assert.Equal(t, algorithm, device.Algorithm, "Unexpected Algorithm for ECDSA algorithm")
}

func TestNewDeviceRSA(t *testing.T) {
	ctx := context.TODO()
	label := "TestDevice"
	algorithm := Algorithm("RSA")

	device, err := NewDevice(ctx, algorithm, label)

	assert.NoError(t, err, "Unexpected error for RSA algorithm")
	assert.NotNil(t, device, "Device should not be nil for RSA algorithm")
	assert.Equal(t, uint64(0), device.SignatureCounter, "Expected SignatureCounter to be 0 for RSA algorithm")
	assert.Equal(t, label, device.Label, "Unexpected Label for RSA algorithm")
	assert.NotEmpty(t, device.PublicKey, "Public key should not be empty for RSA algorithm")
	assert.NotEmpty(t, device.PrivateKey, "Private key should not be empty for RSA algorithm")
	assert.Equal(t, algorithm, device.Algorithm, "Unexpected Algorithm for RSA algorithm")
}

func TestNewDeviceQuantic(t *testing.T) {
	ctx := context.TODO()
	label := "TestDevice"
	algorithm := Algorithm("QUANTIC")

	device, err := NewDevice(ctx, algorithm, label)

	assert.Error(t, err, "Expected error for QUANTIC algorithm")
	assert.Nil(t, device, "Device should be nil for QUANTIC algorithm")
}
