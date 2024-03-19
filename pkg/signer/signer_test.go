package signer

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	device_pkg "signature.service/pkg/device"
)

func TestNewSignerKnownAlgorithmECDSA(t *testing.T) {
	device := &device_pkg.Device{Algorithm: device_pkg.ECDSA}
	signer, err := NewSigner(device)

	assert.NoError(t, err, "Unexpected error for known algorithm ECDSA")
	assert.NotNil(t, signer, "Signer should not be nil for known algorithm ECDSA")
}

func TestNewSignerKnownAlgorithmRSA(t *testing.T) {
	device := &device_pkg.Device{Algorithm: device_pkg.RSA}
	signer, err := NewSigner(device)

	assert.NoError(t, err, "Unexpected error for known algorithm RSA")
	assert.NotNil(t, signer, "Signer should not be nil for known algorithm RSA")
}

func TestNewSignerUnknownAlgorithm(t *testing.T) {
	device := &device_pkg.Device{Algorithm: "UNKNOWN"}
	signer, err := NewSigner(device)

	assert.Error(t, err, "Expected error for unknown algorithm")
	assert.Nil(t, signer, "Signer should be nil for unknown algorithm")
}

func TestGenerateSecuredDataToBeSignedCounterNotZero(t *testing.T) {
	device := &device_pkg.Device{
		Id:               "testID",
		SignatureCounter: 10,
		LastSignature:    "testLastSignature",
	}
	signer := &Signer{device: device}

	data := []byte("testData")

	expectedSecuredData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, device.LastSignature)
	expectedSecuredDataBytes := []byte(expectedSecuredData)

	actualSecuredData := signer.generateSecuredDataToBeSigned(context.Background(), data)

	assert.Equal(t, expectedSecuredDataBytes, actualSecuredData, "Secured data does not match for counter not zero")
}

func TestGenerateSecuredDataToBeSignedCounterZero(t *testing.T) {
	device := &device_pkg.Device{
		Id:               "testID",
		SignatureCounter: 0,
	}
	signer := &Signer{device: device}

	data := []byte("testData")

	expectedLastSignature := base64.StdEncoding.EncodeToString([]byte(device.Id))
	expectedSecuredData := fmt.Sprintf("0_%s_%s", data, expectedLastSignature)
	expectedSecuredDataBytes := []byte(expectedSecuredData)

	actualSecuredData := signer.generateSecuredDataToBeSigned(context.Background(), data)

	assert.Equal(t, expectedSecuredDataBytes, actualSecuredData, "Secured data does not match for counter zero")
}
