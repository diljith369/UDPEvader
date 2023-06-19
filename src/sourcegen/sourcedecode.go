package sourcegen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

type Decoder struct {
	DecodedSource string
}

type ISourceBuilderFunctions interface {
	DecodeSource(encodedsource string)
}

func (oDecoder *Decoder) DecodeSource(encodedsource string) {
	oDecoder.DecodedSource = decodeanddecompress(encodedsource)
}

func checkandreturnerror(err error) string {
	return fmt.Sprintf("Error: %s", err.Error())
}

func decodeanddecompress(encodedsource string) string {
	dcodedval, err := base64.StdEncoding.DecodeString(encodedsource)
	if err != nil {
		return checkandreturnerror(err)
	}
	zipbuffer := bytes.NewReader(dcodedval)
	decompressedbytes, _ := gzip.NewReader(zipbuffer)
	decodedsource, err := io.ReadAll(decompressedbytes)
	if err != nil {
		return checkandreturnerror(err)
	}
	return string(decodedsource)
}
