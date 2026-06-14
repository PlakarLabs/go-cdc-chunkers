# Conformance vectors

`vectors.json` is the **cross-implementation contract** for the chunkers in this
repository. A port of any algorithm — in any language — can regenerate the exact
same inputs from the recipe below, run its own chunker, and assert that its
chunk boundaries match `chunk_lengths`. If they do, the two implementations are
byte-for-byte interoperable and their chunks will deduplicate against each other.

These vectors are produced and checked by `tests/vectors_test.go`. Regenerate
them after an *intentional* behaviour change with:

```sh
go test ./tests/ -run TestVectors -update
```

Without `-update`, the test asserts the current code still reproduces every
vector. A drift means a code change altered cut points — which is a breaking
change for deduplication and must be deliberate.

## File format

`vectors.json` is a JSON object keyed by `"<algorithm>|<size-profile>|<recipe>"`.
Each value is one vector:

| field           | meaning |
|-----------------|---------|
| `algorithm`     | registered algorithm name (e.g. `fastcdc`, `ultracdc-v1.0.0`) |
| `input_recipe`  | how to reconstruct the input bytes (see below) |
| `input_len`     | input length in bytes |
| `input_sha256`  | SHA-256 of the input bytes, hex — confirm your reconstruction before comparing cuts |
| `opts`          | `{min_size, normal_size, max_size, key_hex}` — `key_hex` absent means no key |
| `chunk_lengths` | the full ordered list of chunk sizes; **this is the contract** |
| `cuts_hash`     | SHA-256 over the little-endian uint64 cut-length sequence (a one-shot equality check for `chunk_lengths`) |
| `content_hash`  | SHA-256 of the reconstructed content; equals `input_sha256` for a lossless chunker |

## Input recipes

Inputs are described by recipe rather than stored inline, so the file stays
small and ASCII. Reconstruct `input_len` bytes as follows:

- **`zeros`** — `input_len` bytes, all `0x00`.
- **`repeat-plakar`** — the ASCII bytes of `"plakar"` repeated; byte *i* is
  `"plakar"[i % 6]`. Highly compressible; exercises the failure-to-cut tail.
- **`random-seed0`** — Go's `math/rand` stream seeded with `0`, via
  `rand.New(rand.NewSource(0)).Read(buf)`. This is the Go PRNG's exact byte
  stream. A port in another language **cannot** reproduce it bit-for-bit; for
  those, treat `input_sha256` as authoritative — ship an equivalent
  incompressible corpus of your own and pin its hash, or skip `random-seed0`
  and rely on `zeros` / `repeat-plakar`, which are trivially portable.

## Algorithm variants and known spec divergences

The repo registers both *legacy* (boundary-stable) and *spec-faithful*
(`-v1.x.x`) variants. The vectors capture **both**, so a port can target
whichever it needs. The divergences below are deliberate and tested:

- **`ultracdc`** returns the cut at `i + j` (the exact matching byte inside the
  8-byte window). The UltraCDC paper (Zhou et al., IPCCC 2022, Algorithm 1)
  returns `i + 8` (the window's right edge). They are **not** byte-compatible.
  Use **`ultracdc-v1.0.0`**, which returns `i + 8` and matches the paper
  exactly, for interoperability with reference UltraCDC implementations.

- **`jc`** / **`jc-v1.0.0`** include a `n <= NormalSize: return n` early-return
  that the JC paper (Jin et al., IEEE TPDS 2023, Algorithm 1) does not have.
  The core loop is byte-perfect; the divergence only affects a final partial
  segment at EOF, which may be returned uncut where the paper would split it.
  Use **`jc-v1.1.0`** for the paper-faithful behaviour (legacy masks, no
  early-return). Note `jc-v1.0.0` differs from legacy `jc` only in using
  *generated* masks (same one-bit counts, different bit positions — which the
  paper states does not affect dedup).

- **`fastcdc`**, **`fastcdc-v1.0.0`**, **`kfastcdc`** are fully faithful to the
  FastCDC paper (Xia et al., USENIX ATC 2016, Algorithm 1). `-v1.0.0` uses
  generated masks; `kfastcdc` derives its Gear table from `key_hex` (blake3).

- **`fixed-v1.0.0`** is a fixed-size chunker (cuts every `normal_size` bytes);
  included as a control.

See the papers under `papers/` in the repository root for the canonical
algorithm definitions.
