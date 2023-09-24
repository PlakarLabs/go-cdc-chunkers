package tests

import (
	"bytes"
	"crypto/sha256"
	"io"
	"math/rand"
	"testing"

	chunkers "github.com/PlakarLabs/go-cdc-chunkers"
	_ "github.com/PlakarLabs/go-cdc-chunkers/chunkers/fastcdc"
	_ "github.com/PlakarLabs/go-cdc-chunkers/chunkers/ultracdc"
)

const (
	datalen = 128 << 20
)

var rb, _ = io.ReadAll(io.LimitReader(rand.New(rand.NewSource(0)), datalen))

func Test_FastCDC(t *testing.T) {
	r := bytes.NewReader(rb)

	hasher := sha256.New()
	hasher.Write(rb)
	sum1 := hasher.Sum(nil)

	hasher.Reset()

	chunker, err := chunkers.NewChunker("fastcdc", r)
	if err != nil {
		t.Fatalf(`chunker error: %s`, err)
	}
	for err := error(nil); err == nil; {
		chunk, err := chunker.Next()
		if err != nil && err != io.EOF {
			t.Fatalf(`chunker error: %s`, err)
		}
		if len(chunk) < int(chunker.MinSize()) && err != io.EOF {
			t.Fatalf(`chunker return a chunk below MinSize before last chunk: %s`, err)
		}
		if len(chunk) > int(chunker.MaxSize()) {
			t.Fatalf(`chunker return a chunk above MaxSize`)
		}
		hasher.Write(chunk)
		if err == io.EOF {
			break
		}
	}
	sum2 := hasher.Sum(nil)

	if !bytes.Equal(sum1, sum2) {
		t.Fatalf(`chunker produces incorrect output`)
	}
}

func Test_UltraCDC(t *testing.T) {
	r := bytes.NewReader(rb)

	hasher := sha256.New()
	hasher.Write(rb)
	sum1 := hasher.Sum(nil)

	hasher.Reset()

	chunker, err := chunkers.NewChunker("ultracdc", r)
	if err != nil {
		t.Fatalf(`chunker error: %s`, err)
	}
	for err := error(nil); err == nil; {
		chunk, err := chunker.Next()
		if err != nil && err != io.EOF {
			t.Fatalf(`chunker error: %s`, err)
		}
		if len(chunk) < int(chunker.MinSize()) && err != io.EOF {
			t.Fatalf(`chunker return a chunk below MinSize before last chunk: %s`, err)
		}
		if len(chunk) > int(chunker.MaxSize()) {
			t.Fatalf(`chunker return a chunk above MaxSize`)
		}
		hasher.Write(chunk)
		if err == io.EOF {
			break
		}
	}
	sum2 := hasher.Sum(nil)

	if !bytes.Equal(sum1, sum2) {
		t.Fatalf(`chunker produces incorrect output`)
	}
}

func Benchmark_FastCDC(b *testing.B) {
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	nchunks := 0
	for i := 0; i < b.N; i++ {
		chunker, err := chunkers.NewChunker("fastcdc", r)
		if err != nil {
			b.Fatalf(`chunker error: %s`, err)
		}
		for err := error(nil); err == nil; {
			_, err = chunker.Next()
			nchunks++
		}
		r.Reset(rb)
	}
	b.ReportMetric(float64(nchunks)/float64(b.N), "chunks")
}

func Benchmark_UltraCDC(b *testing.B) {
	r := bytes.NewReader(rb)
	b.SetBytes(int64(r.Len()))
	b.ResetTimer()
	nchunks := 0
	for i := 0; i < b.N; i++ {
		chunker, err := chunkers.NewChunker("ultracdc", r)
		if err != nil {
			b.Fatalf(`chunker error: %s`, err)
		}
		for err := error(nil); err == nil; {
			_, err = chunker.Next()
			nchunks++
		}
		r.Reset(rb)
	}
	b.ReportMetric(float64(nchunks)/float64(b.N), "chunks")
}
