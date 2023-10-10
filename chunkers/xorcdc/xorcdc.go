/*
 * Copyright (c) 2021 Gilles Chehade <gilles@poolp.org>
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

package xorcdc

import (
	"encoding/binary"
	"errors"
	"unsafe"

	chunkers "github.com/PlakarLabs/go-cdc-chunkers"
)

func init() {
	chunkers.Register("xorcdc", newXorCDC)
}

var errNormalSize = errors.New("NormalSize is required and must be 64B <= NormalSize <= 1GB")
var errMinSize = errors.New("MinSize is required and must be 64B <= MinSize <= 1GB && MinSize < NormalSize")
var errMaxSize = errors.New("MaxSize is required and must be 64B <= MaxSize <= 1GB && MaxSize > NormalSize")

type XorCDC struct {
}

func newXorCDC() chunkers.ChunkerImplementation {
	return &XorCDC{}
}

func (c *XorCDC) DefaultOptions() *chunkers.ChunkerOpts {
	return &chunkers.ChunkerOpts{
		MinSize:    2 * 1024,
		MaxSize:    64 * 1024,
		NormalSize: 8 * 1024,
	}
}

func (c *XorCDC) Validate(options *chunkers.ChunkerOpts) error {
	if options.NormalSize == 0 || options.NormalSize < 64 || options.NormalSize > 1024*1024*1024 {
		return errNormalSize
	}
	if options.MinSize < 64 || options.MinSize > 1024*1024*1024 || options.MinSize >= options.NormalSize {
		return errMinSize
	}
	if options.MaxSize < 64 || options.MaxSize > 1024*1024*1024 || options.MaxSize <= options.NormalSize {
		return errMaxSize
	}
	return nil
}

func countSetBits(num uint64) int {
	count := 0
	for num > 0 {
		num &= (num - 1) // Brian Kernighan’s Algorithm
		count++
	}
	return count
}

func (c *XorCDC) Algorithm(options *chunkers.ChunkerOpts, data []byte, n int) int {
	MinSize := options.MinSize
	MaxSize := options.MaxSize
	NormalSize := options.NormalSize

	const (
		bitPattern = 0x55555555
		threshold  = 16
		windowSize = 64
	)

	switch {
	case n <= MinSize:
		return n
	case n >= MaxSize:
		n = MaxSize
	case n <= NormalSize:
		NormalSize = n
	}

	for i := 0; i < n; {
		end := i + windowSize
		if end > n {
			end = n
		}

		var word uint64
		// Ensure there are enough bytes to fill a 64-bit word
		if end-i >= 8 {
			word = binary.BigEndian.Uint64(data[i : i+8])
		} else {
			// If less than 8 bytes remain, handle it appropriately
			// e.g., by copying the remaining bytes into word
			copy((*[8]byte)(unsafe.Pointer(&word))[:], data[i:end])
		}

		// Apply XOR with the bit pattern and count the set bits
		xorResult := word ^ bitPattern
		setBitsCount := countSetBits(xorResult)

		// Check if the count of set bits exceeds the threshold
		if setBitsCount >= threshold {
			chunkSize := end - i
			if chunkSize < MinSize && n-i >= MinSize {
				return i + MinSize // Ensure chunk is at least minSize
			}
			return end // Return the ending index of the window as the cutpoint offset
		}

		i += windowSize // Slide the window
	}

	// If reached here, return maxSize or the remaining data size, whichever is smaller
	if n >= MaxSize {
		return MaxSize
	}
	return n
}
