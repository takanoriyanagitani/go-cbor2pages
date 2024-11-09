package log

import (
	"bytes"
	"context"
	"errors"
	"io"

	cp "github.com/takanoriyanagitani/go-cbor2pages"
)

var (
	ErrNoSpace         error = errors.New("no more space to add a CBOR log")
	ErrInvalidPageSize error = errors.New("invalid page size")
)

//go:generate stringer -type=LogPageSize line.go
type LogPageSize uint16

const (
	LogPageXs LogPageSize = 8
	LogPageSm LogPageSize = 64
	LogPageMd LogPageSize = 512
	LogPageLg LogPageSize = 4096
	LogPageXl LogPageSize = 32768
)

var LogPageSizeMap map[string]LogPageSize = map[string]LogPageSize{
	"Xs": LogPageXs,
	"Sm": LogPageSm,
	"Md": LogPageMd,
	"Lg": LogPageLg,
	"Xl": LogPageXl,
}

func LogPageSizeFromStr(s string) (LogPageSize, error) {
	sz, found := LogPageSizeMap[s]
	switch found {
	case true:
		return sz, nil
	default:
		return LogPageXs, ErrInvalidPageSize
	}
}

type PadPage func(*bytes.Buffer, LogPageSize, cp.PadByte)

func PadPageSimple(page *bytes.Buffer, sz LogPageSize, pad cp.PadByte) {
	var current uint16 = uint16(page.Len())
	if uint16(sz) <= current {
		return
	}

	page.Grow(int(sz))

	var padCount uint16 = uint16(sz) - current
	for i := 0; i < int(padCount); i++ {
		_ = page.WriteByte(byte(pad)) // always nil error or panic
	}
}

var PadPageDefault PadPage = PadPageSimple

func AddSerializedCbor(
	page *bytes.Buffer,
	sz LogPageSize,
	doc cp.CborDocumentRaw,
) error {
	var current uint32 = uint32(page.Len())
	var next uint32 = current + uint32(len(doc))
	if uint32(sz) < next {
		return ErrNoSpace
	}

	_, _ = page.Write(doc) // always nil error or panic
	return nil             // TODO
}

type OutputBuf func(context.Context, *bytes.Buffer) error

type Writer struct {
	io.Writer
}

func (w Writer) ToOutputBuf() OutputBuf {
	return func(_ context.Context, buf *bytes.Buffer) error {
		_, e := buf.WriteTo(w.Writer)
		return e
	}
}

func (p PadPage) WriteAll(
	ctx context.Context,
	docs cp.CborRawDocuments,
	output OutputBuf,
	sz LogPageSize,
	pad cp.PadByte,
) error {
	var buf bytes.Buffer
	var err error = nil

	for doc := range docs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = AddSerializedCbor(&buf, sz, doc)
		if nil != err {
			if errors.Is(err, ErrNoSpace) {
				PadPageSimple(&buf, sz, pad)
				e := output(ctx, &buf)
				if nil != e {
					return e
				}
				buf.Reset()
				err = AddSerializedCbor(&buf, sz, doc)
				if nil != err {
					return err // the doc may be too big for the page size
				}
			} else {
				return err // unexpected error
			}
		}
	}

	return output(ctx, &buf)
}

func (p PadPage) WriteAllDefault(
	ctx context.Context,
	docs cp.CborRawDocuments,
	output OutputBuf,
) error {
	return p.WriteAll(
		ctx,
		docs,
		output,
		LogPageLg,
		cp.PadByteDefault,
	)
}
