package cmderrors

import (
	"fmt"
	"os"
)

func FailErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
