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

```
goos: darwin
goarch: arm64
pkg: github.com/PlakarLabs/go-cdc-chunkers/tests
Benchmark_FastCDC-8           15          75999261 ns/op        1766.04 MB/s         14327 chunks         131256 B/op          4 allocs/op
Benchmark_UltraCDC-8          12          88462517 ns/op        1517.23 MB/s          3945 chunks         131256 B/op          4 allocs/op
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