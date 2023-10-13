# go-cdc-chunkers

## WARNING:
This is a work in progress.

Feel free to join [my discord server](https://discord.com/invite/YC6j4rbvSk) or start discussions at Github.


## Overview
`go-cdc-chunkers` is a Golang package designed to provide unified access to multiple Content-Defined Chunking (CDC) algorithms.
With a simple and intuitive interface, users can effortlessly chunk data using their preferred CDC algorithm.

## Use-cases
Content-Defined Chunking (CDC) algorithms are used in data deduplication and backup systems to break up data into smaller chunks based on their content, rather than their size or location. This allows for more efficient storage and transfer of data, as identical chunks can be stored or transferred only once. CDC algorithms are useful because they can identify and isolate changes in data, making it easier to track and manage changes over time. Additionally, CDC algorithms can be optimized for performance, allowing for faster and more efficient processing of large amounts of data.


## Features
- Unified interface for multiple CDC algorithms.
- Supported algorithms: fastcdc, ultracdc.
- Efficient and optimized for performance.
- Comprehensive error handling.

## Installation
```sh
go get github.com/PlakarLabs/go-cdc-chunkers
```


## Usage
Here's a basic example of how to use the package:

```go
    chunker, err := chunkers.NewChunker("fastcdc", rd)   // or ultracdc
    if err != nil {
        log.Fatal(err)
    }

    offset := 0
    for {
        chunk, err := chunker.Next()
        if err != nil && err != io.EOF {
            log.Fatal(err)
        }

        chunkLen := len(chunk)
        fmt.Println(offset, chunkLen)

        if err == io.EOF {
            // no more chunks to read
            break
        }
        offset += chunkLen
    }
```

## Benchmarks
Performances is a key feature in CDC, `go-cdc-chunkers` strives at optimizing its implementation of CDC algorithms,
finding the proper balance in usability, CPU-usage and memory-usage.

The following benchmark shows the performances of chunking 1GB of random data,
with a minimum chunk size of 256KB and a maximum chunk size of 1MB,
for multiple implementations available as well as multiple methods of consumption of the chunks:

```
goos: darwin
goarch: arm64
pkg: github.com/PlakarLabs/go-cdc-chunkers/tests
Benchmark_Restic_Rabin_Next-8                          1        1926270208 ns/op         557.42 MB/s          1286 chunks        8922128 B/op         11 allocs/op
Benchmark_Askeladdk_FastCDC_Copy-8                     2         686870770 ns/op        1563.24 MB/s        105327 chunks        1048592 B/op          1 allocs/op
Benchmark_Jotfs_FastCDC_Next-8                         3         473524972 ns/op        2267.55 MB/s          1725 chunks        2097264 B/op          2 allocs/op
Benchmark_Tigerwill90_FastCDC_Split-8                  3         395610639 ns/op        2714.14 MB/s          2013 chunks        2097328 B/op          3 allocs/op
Benchmark_Mhofmann_FastCDC_Next-8                      2         592342938 ns/op        1812.70 MB/s          1718 chunks        1048688 B/op          2 allocs/op
Benchmark_PlakarLabs_FastCDC_Copy-8                    9         155380773 ns/op        6910.39 MB/s          3647 chunks        2097318 B/op          3 allocs/op
Benchmark_PlakarLabs_FastCDC_Split-8                   9         120039241 ns/op        8944.92 MB/s          3647 chunks        2097314 B/op          3 allocs/op
Benchmark_PlakarLabs_FastCDC_Next-8                    9         138454935 ns/op        7755.17 MB/s          3647 chunks        2097314 B/op          3 allocs/op
Benchmark_PlakarLabs_UltraCDC_Copy-8                  20          52055098 ns/op        20627.03 MB/s         4096 chunks        2097314 B/op          3 allocs/op
Benchmark_PlakarLabs_UltraCDC_Split-8                 24          48617734 ns/op        22085.39 MB/s         4096 chunks        2097313 B/op          3 allocs/op
Benchmark_PlakarLabs_UltraCDC_Next-8                  24          48486024 ns/op        22145.39 MB/s         4096 chunks        2097312 B/op          3 allocs/op
PASS
ok      github.com/PlakarLabs/go-cdc-chunkers/tests     28.622s
```

## Contributing
We welcome contributions!
If you have a feature request, bug report, or wish to contribute code, please open an issue or pull request.

## Support
If you find `go-cdc-chunkers` useful, please consider supporting its development by [sponsoring the project on GitHub](https://github.com/sponsors/poolpOrg).
Your support helps ensure the project's continued maintenance and improvement.


## License
This project is licensed under the ISC License. See the [LICENSE.md](LICENSE.md) file for details.


## Reference

  - [Xia, Wen, et al. "Fastcdc: a fast and efficient content-defined chunking approach for data deduplication." 2016 USENIX Annual Technical Conference](https://www.usenix.org/system/files/conference/atc16/atc16-paper-xia.pdf)
  - [Zhou, Wang, Xia, Zhang "UltraCDC:A Fast and Stable Content-Defined Chunking Algorithm for Deduplication-based Backup Storage Systems" 2022 IEEE](https://ieeexplore.ieee.org/document/9894295)