package semigraph

import "unicode/utf8"

// unicode_octant_map_entry defines a mapping for octant bits.
type unicodeOctantMapEntry struct {
	OctantBits uint8
	Data       uint8
}

// unicodeOctantMap contains predefined mappings for specific octant bits.
var unicodeOctantMap = []unicodeOctantMapEntry{
	{0x00, 0x00 /* U+00A0 */},
	{0x01, 0xA8 /* U+1CEA8 */},
	{0x02, 0xAB /* U+1CEAB */},
	{0x03, 0xC2 /* U+1FB82 */},
	{0x05, 0x98 /* U+2598 */},
	{0x0A, 0x9D /* U+259D */},
	{0x0F, 0x80 /* U+2580 */},
	{0x14, 0xE6 /* U+1FBE6 */},
	{0x28, 0xE7 /* U+1FBE7 */},
	{0x3F, 0xC5 /* U+1FB85 */},
	{0x40, 0xA3 /* U+1CEA3 */},
	{0x50, 0x96 /* U+2596 */},
	{0x55, 0x8C /* U+258C */},
	{0x5A, 0x9E /* U+259E */},
	{0x5F, 0x9B /* U+259B */},
	{0x80, 0xA0 /* U+1CEA0 */},
	{0xA0, 0x97 /* U+2597 */},
	{0xA5, 0x9A /* U+259A */},
	{0xAA, 0x90 /* U+2590 */},
	{0xAF, 0x9C /* U+259C */},
	{0xC0, 0x82 /* U+2582 */},
	{0xF0, 0x84 /* U+2584 */},
	{0xF5, 0x99 /* U+2599 */},
	{0xFA, 0x9F /* U+259F */},
	{0xFC, 0x86 /* U+2586 */},
	{0xFF, 0x88 /* U+2588 */},
}

// findUnicodeOctantMapData performs a binary search on unicodeOctantMap.
// Returns the data byte if found, otherwise returns -(index where it would be inserted).
func findUnicodeOctantMapData(octantBits uint8) int {
	first := 0
	last := len(unicodeOctantMap)
	for first < last {
		i := (first + last) / 2
		if octantBits == unicodeOctantMap[i].OctantBits {
			return int(unicodeOctantMap[i].Data)
		}
		if octantBits > unicodeOctantMap[i].OctantBits {
			first = i + 1
		} else {
			last = i
		}
	}
	return -first
}

// octantBitsToUTF8 converts octant bits to a UTF-8 byte slice.
// It returns the slice and its length.
func octantBitsToUTF8(octantBits uint8) rune {
	// Base UTF-8 sequences as translated from C
	var utf8Sequences = [6][4]byte{
		{0xF0, 0x9C, 0xB4, 0x80}, // For U+1CD00 base range
		{0xC2, 0xA0, 0x00, 0x00}, // For U+00A0
		{0xE2, 0x96, 0x80, 0x00}, // For U+2580 base range
		{0xF0, 0x9C, 0xBA, 0xA0}, // For U+1CEA0 base range
		{0xF0, 0x9F, 0xAE, 0x80}, // For U+1FB80 base range
		{0xF0, 0x9F, 0xAF, 0xA6}, // For U+1FBE6 base range
	}

	data := findUnicodeOctantMapData(octantBits)
	var seqStart []byte
	var codeOffset uint8

	if data < 0 {
		seqStart = utf8Sequences[0][:]
		codeOffset = uint8(octantBits) + uint8(data) // data is negative, so this effectively subtracts
	} else if data == 0x00 {
		seqStart = utf8Sequences[1][:]
		codeOffset = 0
	} else {
		seqStart = utf8Sequences[2+((data>>5)&0x3)][:]
		codeOffset = uint8(data) & 0x1F
	}

	// Determine length based on the first byte, as in the C code
	var length int
	if (seqStart[0] >> 4) > 0xE {
		length = 4
	} else if (seqStart[0] >> 4) == 0xE {
		length = 3
	} else {
		length = 2
	}

	outBuf := make([]byte, 4) // Max 4 bytes for UTF-8
	copy(outBuf, seqStart)

	// Apply offsets
	if length >= 1 {
		outBuf[length-1] |= (codeOffset & 0x3F)
	}
	if length >= 2 {
		outBuf[length-2] |= (codeOffset >> 6)
	}
	// Note: The C code sets out_buf[4] = 0x00, which is for null termination.
	// In Go, we return a slice of the correct length, so no null termination needed.

	out, _ := utf8.DecodeRune(outBuf)
	return out
}
