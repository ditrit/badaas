package validator_test

import (
	"testing"

	"github.com/ditrit/badaas/validator"
)

func TestValidEmail(t *testing.T) {
	testCases := []struct {
		desc string
		in   string
		out  bool
	}{
		{
			desc: "missing @",
			in:   "bob.bobemail.com",
			out:  false,
		}, {
			desc: "missing domain",
			in:   "bob.bob@",
			out:  false,
		}, {
			desc: "odd chars",
			in:   "bob.bob%@email.com",
			out:  true,
		}, {
			desc: "valid",
			in:   "bob.bob@email.com",
			out:  true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if res := validator.ValidEmail(tC.in); res != tC.out {
				t.Errorf("got=%v, expected=%v", res, tC.out)
			}
		})
	}
}
