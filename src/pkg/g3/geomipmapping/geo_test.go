package geomipmapping

import (
	"testing"
)

func TestPow(t *testing.T) {
	if pow2(0) != 1 {
		t.Error("pow2(0) != 1")
	}

	if pow2(1) != 2 {
		t.Error("pow2(1) != 2")
	}

	if pow2(3) != 8 {
		t.Error("pow2(3) != 8")
	}
}

func TestIndices(t *testing.T) {
	NewGeoMipMap()

	fmt.Println("", createPatchIndices(0, 2))
}
