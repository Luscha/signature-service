package signer

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	device_pkg "signature.service/pkg/device"
)

func TestSignerSign(t *testing.T) {
	ctx := context.TODO()
	// Create a mock device
	mockDevice, _ := device_pkg.NewDevice(ctx, device_pkg.ECDSA, "")
	signer, _ := NewSigner(mockDevice)

	// Data to be signed
	data := []byte("test_data")

	// Get initial values of SignatureCounter and LastSignature
	initialCounter := mockDevice.SignatureCounter
	initialLastSignature := mockDevice.LastSignature

	signedDataExpected := signer.generateSecuredDataToBeSigned(ctx, data)

	// Call the Sign method
	signature, err := signer.Sign(ctx, data)

	// Verify the result
	assert.NoError(t, err, "Unexpected error during signing")
	assert.NotNil(t, signature, "Signature should not be nil")

	// Verify changes in the device's state
	assert.Equal(t, initialCounter+1, mockDevice.SignatureCounter, "SignatureCounter should be incremented")
	assert.NotEqual(t, initialLastSignature, mockDevice.LastSignature, "LastSignature should be updated")
	assert.Equal(t, signature.Signature, mockDevice.LastSignature, "LastSignature should be updated to signature")

	// Decode the signature and verify its components
	_, err = base64.StdEncoding.DecodeString(signature.Signature)
	assert.NoError(t, err, "Error decoding base64 signature")

	assert.Equal(t, mockDevice.Id, signature.Device, "Device ID in signature should match")
	assert.Equal(t, data, signature.Data, "Data in signature should match")
	assert.Equal(t, signedDataExpected, signature.SignedData, "SignedData in signature should match")
}

func TestSignerSignECDSA(t *testing.T) {
	ctx := context.TODO()
	label := "TestDevice"
	algorithm := device_pkg.ECDSA

	// Create a new ECDSA device
	device, err := device_pkg.NewDevice(ctx, algorithm, label)
	assert.NoError(t, err, "Failed to create ECDSA device")
	assert.NotNil(t, device, "Device should not be nil")

	// Create a signer
	signer, err := NewSigner(device)
	assert.NoError(t, err, "Failed to create signer")
	assert.NotNil(t, signer, "Signer should not be nil")

	// Sign some data
	data := []byte("test_data")
	signature, err := signer.Sign(ctx, data)
	assert.NoError(t, err, "Failed to sign data")
	assert.NotNil(t, signature, "Signature should not be nil")

	// Verify the signature using the device's public key
	err = signer.Verify(ctx, signature)
	assert.NoError(t, err, "Failed to verify ECDSA signature")
}

func TestSignerSignRSA(t *testing.T) {
	ctx := context.TODO()
	label := "TestDevice"
	algorithm := device_pkg.RSA

	// Create a new RSA device
	device, err := device_pkg.NewDevice(ctx, algorithm, label)
	assert.NoError(t, err, "Failed to create RSA device")
	assert.NotNil(t, device, "Device should not be nil")

	// Create a signer
	signer, err := NewSigner(device)
	assert.NoError(t, err, "Failed to create signer")
	assert.NotNil(t, signer, "Signer should not be nil")

	// Sign some data
	data := []byte("test_data")
	signature, err := signer.Sign(ctx, data)
	assert.NoError(t, err, "Failed to sign data")
	assert.NotNil(t, signature, "Signature should not be nil")

	// Verify the signature using the device's public key
	err = signer.Verify(ctx, signature)
	assert.NoError(t, err, "Failed to verify RSA signature")
}
