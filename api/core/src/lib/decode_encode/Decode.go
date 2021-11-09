package decode_encode

import (
	"encoding/base64"
	"io"
	"strings"
)

func DecodeBase64ToReader(bs64 string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding,strings.NewReader(bs64))
}


func DecodeBase64ToBytes(bs64 string) ([]byte ,error){
	return io.ReadAll(DecodeBase64ToReader(bs64))
}



