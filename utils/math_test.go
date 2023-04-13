package utils_test

import (
	"testing"

	"github.com/ditrit/badaas/utils"
)

func TestIsAnINT(t *testing.T) {
	data := map[float64]bool{
		6.32:      false,
		-45.3:     false,
		-45.0:     true,
		0.0:       true,
		1.0000001: false,
	}
	for k, v := range data {
		res := utils.IsAnInt(k)
		if res != v {
			t.Errorf("expected %v got %v for %f", v, res, k)
		}
	}
}
