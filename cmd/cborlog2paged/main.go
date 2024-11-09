package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	cp "github.com/takanoriyanagitani/go-cbor2pages"
	util "github.com/takanoriyanagitani/go-cbor2pages/util"

	ap "github.com/takanoriyanagitani/go-cbor2pages/app/cbor2paged/log"
	lp "github.com/takanoriyanagitani/go-cbor2pages/log"

	aa "github.com/takanoriyanagitani/go-cbor2pages/cbor/arr2cbor/amacker"
	ca "github.com/takanoriyanagitani/go-cbor2pages/cbor/cbor2arr/amacker"
)

func EnvVarIoNew(key string) util.IO[string] {
	return func(_ context.Context) (string, error) {
		return os.Getenv(key), nil
	}
}

var padByteSource ap.PadByteSource = ap.PadByteSourceDefault

var pageSzSource ap.PageSizeSource = ap.PageSizeSource(util.ComposeIoErr(
	EnvVarIoNew("ENV_PAGE_SIZE4LOG"),
	lp.LogPageSizeFromStr,
))

var cfgSource ap.ConfigSource = ap.ConfigSourceNew(pageSzSource, padByteSource)

var padPage lp.PadPage = lp.PadPageDefault

var appSource ap.AppSource = cfgSource.ToAppSource(padPage)

var any2buf cp.CborSerializerBuf = aa.CborSerializerBufDefault
var any2raw cp.CborSerializer = any2buf.ToSerializer()

type IoConfig struct {
	io.Reader
	io.Writer
}

func (i IoConfig) ToCborMapIter() cp.CborMapIter {
	return ca.CborToArrNew(i.Reader).ToCborMapIter()
}

func (i IoConfig) ToOutputBuf() lp.OutputBuf {
	return lp.Writer{Writer: i.Writer}.ToOutputBuf()
}

func rdr2wtr(ctx context.Context, rdr io.Reader, wtr io.Writer) error {
	var br io.Reader = bufio.NewReader(rdr)
	icfg := IoConfig{
		Reader: br,
		Writer: wtr,
	}

	app, e := appSource(ctx)
	if nil != e {
		return e
	}

	return app.WriteMaps(
		ctx,
		icfg.ToCborMapIter(),
		icfg.ToOutputBuf(),
		any2raw,
	)
}

func stdin2stdout(ctx context.Context) error {
	return rdr2wtr(ctx, os.Stdin, os.Stdout)
}

func main() {
	e := stdin2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
