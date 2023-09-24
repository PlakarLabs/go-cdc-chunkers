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

package ultracdc

import (
	"errors"
	"math/bits"
	"unsafe"

	chunkers "github.com/PlakarLabs/go-cdc-chunkers"
)

func init() {
	chunkers.Register("ultracdc", newUltraCDC)
}

var errMinSize = errors.New("MinSize is required and must be 64B <= MinSize <= 1GB")
var errMaxSize = errors.New("MaxSize is required and must be 64B <= MaxSize <= 1GB")

type ChunkerOpts struct {
	minSize uint32
	maxSize uint32
}

func (o *ChunkerOpts) MinSize() uint32 {
	return o.minSize
}

func (o *ChunkerOpts) MaxSize() uint32 {
	return o.maxSize
}

func (o *ChunkerOpts) NormalSize() uint32 {
	return o.minSize + 8*1024
}

func (o *ChunkerOpts) Validate() error {
	if o.minSize < 64 || o.minSize > 1024*1024*1024 {
		return errMinSize
	}
	if o.maxSize < 64 || o.maxSize > 1024*1024*1024 {
		return errMaxSize
	}
	return nil
}

type UltraCDC struct {
}

func newUltraCDC() chunkers.ChunkerImplementation {
	return &UltraCDC{}
}

func (c *UltraCDC) DefaultOptions() chunkers.ChunkerOpts {
	return &ChunkerOpts{
		minSize: 2 * 1024,
		maxSize: 64 * 1024,
	}
}

func (c *UltraCDC) Algorithm(options chunkers.ChunkerOpts, data []byte, n uint32) uint32 {
	src := (*uint64)(unsafe.Pointer(&data[0]))

	const (
		Pattern uint64 = 0xAAAAAAAAAAAAAAAA
		MaskS   uint64 = 0x2F
		MaskL   uint64 = 0x2C
		LEST    uint32 = 64
	)
	MinSize := options.MinSize()
	MaxSize := options.MaxSize()
	NormalSize := options.NormalSize()

	i := MinSize
	cnt := uint32(0)
	mask := MaskS

	switch {
	case n <= MinSize:
		return n
	case n >= MaxSize:
		n = MaxSize
	case n <= NormalSize:
		NormalSize = n
	}

	outBufWin := (*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(src)) + uintptr(i)))
	dist := uint64(bits.OnesCount64(*outBufWin ^ Pattern))
	i += 8

	for i < n {
		if i == NormalSize {
			mask = MaskL
		}

		inBufWin := (*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(src)) + uintptr(i)))
		if (*outBufWin ^ *inBufWin) == 0 {
			cnt++
			if cnt == LEST {
				return i + 8
			}
			i += 8
			continue
		}

		cnt = 0
		for j := 0; j < 8; j++ {
			if (dist & mask) == 0 {
				return i + 8
			}
			inByte := *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(inBufWin)) + uintptr(j)))
			outByte := *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(outBufWin)) + uintptr(j)))
			dist = dist + uint64(ultraCDC_hammingDistanceTable[outByte][inByte])
		}
		outBufWin = inBufWin
		i += 8
	}

	return n
}
