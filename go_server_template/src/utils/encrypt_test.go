package utils

import (
	"fmt"
	"testing"
)

func TestDecrypt(t *testing.T) {
	desc, _ := ConfigDecrypt("1c5ccfc8e831d3606d307d0e8d701ff212ebc468352b097ef30db5cf14437bae")
	fmt.Println("========  " + desc)
}
