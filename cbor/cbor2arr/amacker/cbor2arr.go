package cbor2arr

import (
	"io"

	fc "github.com/fxamacker/cbor/v2"

	cp "github.com/takanoriyanagitani/go-cbor2pages"
)

type CborToArr struct {
	*fc.Decoder
}

func (c CborToArr) ToCborArrayIter() cp.CborArrayIter {
	return func(yield func(cp.CborArray) bool) {
		var buf []any
		var err error = nil

		for {
			clear(buf)
			buf = buf[:0]

			err = c.Decoder.Decode(&buf)
			if nil != err {
				return
			}

			if !yield(buf) {
				return
			}
		}
	}
}

func (c CborToArr) ToCborAnyIter() cp.CborAnyIter {
	return func(yield func(cp.CborAny) bool) {
		var buf any
		var err error = nil

		for {
			buf = nil

			err = c.Decoder.Decode(&buf)
			if nil != err {
				return
			}

			if !yield(buf) {
				return
			}
		}
	}
}

func (c CborToArr) ToCborMapIter() cp.CborMapIter {
	return func(yield func(cp.CborMap) bool) {
		var buf map[string]any
		var err error = nil

		for {
			clear(buf)

			err = c.Decoder.Decode(&buf)
			if nil != err {
				return
			}

			if !yield(buf) {
				return
			}
		}
	}
}

func CborToArrNew(rdr io.Reader) CborToArr {
	return CborToArr{Decoder: fc.NewDecoder(rdr)}
}
