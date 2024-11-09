package log2paged

import (
	"context"
	"errors"
	"iter"

	cp "github.com/takanoriyanagitani/go-cbor2pages"
	util "github.com/takanoriyanagitani/go-cbor2pages/util"
	itools "github.com/takanoriyanagitani/go-cbor2pages/util/iter"

	lp "github.com/takanoriyanagitani/go-cbor2pages/log"
)

type App struct {
	lp.PadPage
	lp.LogPageSize
	cp.PadByte
}

func (a App) WriteAll(
	ctx context.Context,
	docs cp.CborRawDocuments,
	out lp.OutputBuf,
) error {
	return a.PadPage.WriteAll(
		ctx,
		docs,
		out,
		a.LogPageSize,
		a.PadByte,
	)
}

type Config struct {
	lp.LogPageSize
	cp.PadByte
}

func (c Config) ToApp(pad lp.PadPage) App {
	return App{
		PadPage:     pad,
		LogPageSize: c.LogPageSize,
		PadByte:     c.PadByte,
	}
}

type PageSizeSource util.IO[lp.LogPageSize]
type PadByteSource util.IO[cp.PadByte]

func PadByteSourceNewStatic(pb cp.PadByte) PadByteSource {
	return func(_ context.Context) (cp.PadByte, error) {
		return pb, nil
	}
}

var PadByteSourceDefault PadByteSource = PadByteSourceNewStatic(
	cp.PadByteDefault,
)

type ConfigSource util.IO[Config]

func ConfigSourceNew(page PageSizeSource, pad PadByteSource) ConfigSource {
	return func(ctx context.Context) (Config, error) {
		psz, epage := page(ctx)
		pdb, epadb := pad(ctx)
		return Config{
			LogPageSize: psz,
			PadByte:     pdb,
		}, errors.Join(epage, epadb)
	}
}

type AppSource util.IO[App]

func (c ConfigSource) ToAppSource(pad lp.PadPage) AppSource {
	return AppSource(util.ComposeIo(
		util.IO[Config](c),
		func(cfg Config) App { return cfg.ToApp(pad) },
	))
}

func (a App) WriteMaps(
	ctx context.Context,
	maps cp.CborMapIter,
	out lp.OutputBuf,
	ser cp.CborSerializer,
) error {
	var anys cp.CborAnyIter = cp.CborAnyIter(itools.Map[cp.CborMap, cp.CborAny](
		iter.Seq[cp.CborMap](maps),
		func(m cp.CborMap) cp.CborAny { return m },
	))
	var raws cp.CborRawDocuments = anys.ToRawDocuments(ser)
	return a.WriteAll(
		ctx,
		raws,
		out,
	)
}
