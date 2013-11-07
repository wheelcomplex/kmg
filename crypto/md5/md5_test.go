package md5

import (
	"encoding/hex"
	"kmg/kmpTest"
	"strings"
	"testing"
)

func TestBytes(t *testing.T) {
	expect_bytes, _ := hex.DecodeString("d41d8cd98f00b204e9800998ecf8427e")
	kmpTest.Assert(t, Bytes([]byte("")), expect_bytes)
	kmpTest.Assert(t, Bytes([]byte("")), expect_bytes)
	expect_bytes, _ = hex.DecodeString("9e107d9d372bb6826bd81d3542a419d6")
	kmpTest.Assert(t, Bytes([]byte("The quick brown fox jumps over the lazy dog")),
		expect_bytes)

}
func TestHex(t *testing.T) {
	kmpTest.Assert(t, Hex([]byte("")), "d41d8cd98f00b204e9800998ecf8427e")
	kmpTest.Assert(t, Hex([]byte("The quick brown fox jumps over the lazy dog")),
		"9e107d9d372bb6826bd81d3542a419d6")
}
func TestHexFromString(t *testing.T) {
	kmpTest.Assert(t, HexFromString(strings.Repeat("1", 10000)), "b223cca8b360eae4e49568512e2de29f")
}
