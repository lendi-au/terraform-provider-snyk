package snyk

import (
	"testing"
)

// Test provider structure - runs the TF internal validation function to ensure provider structure works.
func TestProvider(t *testing.T) {
	provider := Provider("test")()

	if err := provider.InternalValidate(); err != nil {
		t.Fatal(err)
	}
}
