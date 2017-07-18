package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNonce(t *testing.T) {
	filename := "rumpelstilzchen"
	contents, err := getNonce(filename)
	if err == nil {
		t.Fatalf("ReadFile %s: error expected, none found", filename)
	}

	filename = "ec2_test.go"
	contents, err = getNonce(filename)
	if err != nil {
		t.Fatalf("ReadFile %s: %v", filename, err)
	}

	assert.EqualValues(t, 436, len(contents))
}
