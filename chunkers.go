package chunkers

/*
 * Copyright (c) 2023 Gilles Chehade <gilles@poolp.org>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

import (
	"bufio"
	"errors"
	"io"
)

type ChunkerOpts interface {
	MinSize() uint32
	MaxSize() uint32
	NormalSize() uint32
	Validate() error
}

type ChunkerImplementation interface {
	DefaultOptions() ChunkerOpts
	Algorithm(ChunkerOpts, []byte, uint32) uint32
}

type Chunker struct {
	rd             *bufio.Reader
	options        ChunkerOpts
	implementation ChunkerImplementation
}

func (c *Chunker) MinSize() uint32 {
	return c.options.MinSize()
}

func (c *Chunker) MaxSize() uint32 {
	return c.options.MaxSize()
}

func (c *Chunker) NormalSize() uint32 {
	return c.options.NormalSize()
}

var chunkers map[string]func() ChunkerImplementation = make(map[string]func() ChunkerImplementation)

func Register(name string, implementation func() ChunkerImplementation) error {
	if _, exists := chunkers[name]; exists {
		return errors.New("algorithm already registered")
	}
	chunkers[name] = implementation
	return nil
}

func NewChunker(algorithm string, reader io.Reader) (*Chunker, error) {
	var implementationAllocator func() ChunkerImplementation

	implementationAllocator, exists := chunkers[algorithm]
	if !exists {
		return nil, errors.New("unknown algorithm")
	}

	chunker := &Chunker{}
	chunker.implementation = implementationAllocator()
	chunker.options = chunker.implementation.DefaultOptions()
	err := chunker.options.Validate()
	if err != nil {
		return nil, err
	}
	chunker.rd = bufio.NewReaderSize(reader, int(chunker.implementation.DefaultOptions().MaxSize())*2)
	return chunker, nil
}

func (chunker *Chunker) Next() ([]byte, error) {
	data, err := chunker.rd.Peek(int(chunker.options.MaxSize()))
	if err != nil && err != io.EOF && err != bufio.ErrBufferFull {
		return nil, err
	}

	n := len(data)
	if n == 0 {
		return nil, io.EOF
	}

	cutpoint := chunker.implementation.Algorithm(chunker.options, data[:n], uint32(n))
	if _, err = chunker.rd.Discard(int(cutpoint)); err != nil {
		return nil, err
	}

	return data[:cutpoint], nil
}
