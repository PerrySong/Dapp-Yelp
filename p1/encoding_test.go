package p1

import (
	"fmt"
	"testing"
)

func TestTest_compact_encode(t *testing.T) {
	fmt.Println(compact_encode([]uint8{6, 1, 16}))
	fmt.Println(compact_decode([]uint8{32, 97}))
}
