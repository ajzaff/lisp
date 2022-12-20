package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"io"
	"math/rand"
	"testing"

	"github.com/ajzaff/lisp/print"
	"github.com/ajzaff/lisp/x/fuzzutil"
)

type CompressWriter interface {
	io.Writer
	Reset(io.Writer)
	Flush() error
}

func benchmarkCompress(b *testing.B, zw CompressWriter) {
	g := fuzzutil.NewGenerator(rand.New(rand.NewSource(1337)))
	g.ConsMaxDepth = 10
	g.ConsWeight = 5

	printBuf := bytes.NewBuffer(make([]byte, 0, 256))
	compressBuf := bytes.NewBuffer(make([]byte, 0, 256))
	for i := 0; i < b.N; i++ {
		printBuf.Reset()
		compressBuf.Reset()
		zw.Reset(compressBuf)

		v := g.Next()

		print.StdPrinter(printBuf).Print(v)
		origLen := printBuf.Len()

		if _, err := io.Copy(zw, printBuf); err != nil {
			b.Error(err)
		}
		zw.Flush()

		if compressLen := compressBuf.Len(); compressLen > 0 {
			compressRatio := float64(origLen) / float64(compressLen) / float64(b.N)
			b.ReportMetric(compressRatio, "CompressRatio/op")
		}
	}
}

func BenchmarkCompressFlate(b *testing.B) {
	w, err := flate.NewWriter(nil, flate.DefaultCompression)
	if err != nil {
		b.Fail()
	}
	defer w.Close()
	benchmarkCompress(b, w)
}

func BenchmarkCompressGZip(b *testing.B) {
	w, err := gzip.NewWriterLevel(nil, gzip.DefaultCompression)
	if err != nil {
		b.Fail()
	}
	defer w.Close()
	benchmarkCompress(b, w)
}

type wrapLZWWriter struct {
	*lzw.Writer
}

func (w wrapLZWWriter) Flush() error {
	return nil
}

func (w wrapLZWWriter) Reset(dst io.Writer) {
	w.Writer.Reset(dst, lzw.LSB, 8)
}

func BenchmarkCompressLZW(b *testing.B) {
	w := lzw.NewWriter(nil, lzw.LSB, 8)
	defer w.Close()
	benchmarkCompress(b, wrapLZWWriter{w.(*lzw.Writer)})
}

func BenchmarkCompressZLib(b *testing.B) {
	w, err := zlib.NewWriterLevel(nil, zlib.DefaultCompression)
	if err != nil {
		b.Fail()
	}
	defer w.Close()
	benchmarkCompress(b, w)
}
