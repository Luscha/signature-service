package signer

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"signature.service/pkg/device"
	"signature.service/pkg/signature"
)

type signatureWrapper struct {
	sig     *signature.Signature
	counter uint64
}

func TestConcurrentSigning(t *testing.T) {
	// Create a device
	device, err := device.NewDevice(context.Background(), device.ECDSA, "TestDevice")
	if err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	var wg sync.WaitGroup
	numJobs := 100
	numWorkers := 20

	// Channel to receive sign jobs
	signJobs := make(chan struct{}, numJobs*numWorkers)

	// Channel to receive signatures
	signatures := make(chan *signatureWrapper, numJobs*numWorkers)

	// Send sign jobs to workers
	for i := 0; i < numWorkers*numJobs; i++ {
		wg.Add(1)
		signJobs <- struct{}{}
	}

	// Worker function to perform signing
	worker := func() {
		// Create a signer with the device
		signer, err := NewSigner(device)
		if err != nil {
			t.Fatalf("failed to create signer: %v", err)
		}
		for range signJobs {
			sig, err := signer.Sign(context.Background(), []byte("data to sign"))
			if err != nil {
				t.Errorf("signing error: %v", err)
				wg.Done()
				continue
			}
			parts := strings.Split(string(sig.SignedData), "_")
			if len(parts) < 2 {
				t.Errorf("invalid SignedData format: %s", sig.SignedData)
				continue
			}
			counter, _ := strconv.ParseUint(parts[0], 10, 64)
			signatures <- &signatureWrapper{
				sig:     sig,
				counter: counter,
			}
			wg.Done()
		}
	}

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	wg.Wait()
	close(signJobs)
	close(signatures)

	// Collect signatures
	collectedSignatures := make(map[uint64]struct{})
	for sig := range signatures {
		collectedSignatures[sig.counter] = struct{}{}
	}

	for counter := uint64(0); counter < uint64(numJobs*numWorkers); counter++ {
		if _, ok := collectedSignatures[counter]; !ok {
			t.Errorf("gap detected in signature counters: counter %d missing", counter)
		}
	}

	// Check if the device has the correct last signature
	assert.Equal(t, device.SignatureCounter, uint64(numJobs*numWorkers), "incorrect signautre counter")
}
