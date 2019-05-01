package p1

import (
	"fmt"
	"reflect"
)

// Hex -> ASCII
//

// [1,6,1] -> [17, 97]
func compact_encode(hex_array []uint8) []uint8 {
	term := uint8(0)
	// If hex_array has terminator, ignore the last uint8
	if hasTerm(hex_array) {
		term = 1
		hex_array = hex_array[:len(hex_array)-1]
	}

	oddlen := uint8(len(hex_array) % 2)
	flags := 2*term + oddlen
	// hexarray now has an even length whose first nibble is the flags.

	res := make([]uint8, len(hex_array)/2+1)

	//fmt.Printf("hex_array len = %d \n", len(hex_array))
	//fmt.Printf("reslen = %d \n", len(res))

	var encodeArray []uint8
	if oddlen == 1 { // if hex_array len is odd
		encodeArray = append([]uint8{flags}, hex_array...)
	} else {
		encodeArray = append([]uint8{flags, 0}, hex_array...)
	}

	//fmt.Printf("encode_array = %d \n", encodeArray)
	decodeNibbles(encodeArray, res)
	return res
}

func decodeNibbles(nibbles []uint8, res []uint8) {
	//fmt.Printf("nibbles = %d ", nibbles)
	//fmt.Printf("len(res) = %d ", len(res))
	for bi, ni := 0, 0; ni < len(nibbles); bi, ni = bi+1, ni+2 {
		//fmt.Printf("ni = %d bi = %d \n", ni, bi)
		//fmt.Println(nibbles[ni])
		//fmt.Println(nibbles[ni + 1])
		res[bi] = nibbles[ni]<<4 | nibbles[ni+1]
	}
}

// If Leaf, ignore 16 at the end
// ASCII code array -> HexArray
// [17, 97] -> [1, 1, 6, 1] -> [1, 6, 1]
func compact_decode(compact []uint8) []uint8 {
	base := keybytesToHex(compact)
	// delete terminator flag
	//fmt.Println(base)
	// Do not need 16 no matter what
	//if base[0] < 2 {
	//	//	base = base[:len(base)-1]
	//	//}

	// apply odd flag
	//fmt.Println("curBase = ", base)
	chop := 2 - base[0]&1 // if base[0] = 1, 3, chop = 1, if base[1] = 2 chop = 2

	//if

	return base[chop:]
}

func keybytesToHex(str []uint8) []uint8 {
	l := len(str) * 2 // + 1
	var nibbles = make([]uint8, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	//nibbles[l-1] = 16
	return nibbles
}

func Test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})))
	fmt.Println(compact_encode([]uint8{7, 0, 16}))
}

// hasTerm returns whether a hex key has the terminator flag.
// What is terminator flag?
func hasTerm(s []uint8) bool {
	return len(s) > 0 && s[len(s)-1] == 16
}
