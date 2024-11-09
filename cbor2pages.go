package cbor2pages

import (
	"bytes"
	"context"
	"errors"
	"iter"
)

type CborDocumentRaw []byte

type CborRawDocuments iter.Seq[CborDocumentRaw]

type CborArray []any
type CborArrayIter iter.Seq[CborArray]

type CborMap map[string]any
type CborMapIter iter.Seq[CborMap]

type CborAny any
type CborAnyIter iter.Seq[CborAny]

type CborSerializer func(CborAny) (CborDocumentRaw, error)
type CborSerializerBuf func(CborAny, *bytes.Buffer) error

func (b CborSerializerBuf) ToSerializer() CborSerializer {
	var buf bytes.Buffer
	var err error = nil

	return func(a CborAny) (CborDocumentRaw, error) {
		buf.Reset()
		err = b(a, &buf)
		return buf.Bytes(), err
	}
}

func (a CborAnyIter) ToRawDocuments(ser CborSerializer) CborRawDocuments {
	return func(yield func(CborDocumentRaw) bool) {
		for item := range a {
			serialized, e := ser(item)
			if nil != e {
				return
			}

			if !yield(serialized) {
				return
			}
		}
	}
}

type PageSize uint32

type PadByte uint8

const (
	PadByteDefault PadByte = 0xa0 // an empty map map[string]any
	PadByteArray PadByte = 0x80 // an empty array []
)

const (
	PageSizeXs PageSize = 256
	PageSizeSm PageSize = 1024
	PageSizeMd PageSize = 4096
	PageSizeLg PageSize = 16384
	PageSizeXl PageSize = 65536
)

const (
	PageSizeTiny   PageSize = 16
	PageSizeSmall  PageSize = 256
	PageSizeNormal PageSize = 4096
	PageSizeLarge  PageSize = 65536
	PageSizeLARGE  PageSize = 1048576
	PageSizeHuge   PageSize = 16777216
	PageSizeHUGE   PageSize = 268435456
)

var (
	ErrNoSpace error = errors.New("no spage to add the row")
)

type Row []any

type SerializedRow []byte

type AddRow func(context.Context, *bytes.Buffer, SerializedRow, PageSize) error
