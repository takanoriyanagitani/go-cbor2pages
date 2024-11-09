package arr2cbor

import (
	"bytes"

	fc "github.com/fxamacker/cbor/v2"

	cp "github.com/takanoriyanagitani/go-cbor2pages"
)

func ArrToCborToBuf(
	arr cp.CborAny,
	buf *bytes.Buffer,
) error {
	return fc.MarshalToBuffer(arr, buf)
}

var CborSerializerBufDefault cp.CborSerializerBuf = ArrToCborToBuf
