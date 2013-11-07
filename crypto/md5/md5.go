package md5

import omd5 "crypto/md5"
import "encoding/hex"

//direct sum a md5 value to byte
func Bytes(b []byte) []byte {
	hash := omd5.New()
	hash.Write(b)
	return hash.Sum(nil)
}

//direct sum a md5 value to hex
func Hex(b []byte) string {
	return hex.EncodeToString(Bytes(b))
}

func HexFromString(b string) string {
	return hex.EncodeToString(Bytes([]byte(b)))
}
