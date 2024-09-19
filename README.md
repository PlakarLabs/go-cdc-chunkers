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
goarch: amd64
pkg: github.com/PlakarLabs/go-cdc-chunkers/tests
cpu: VirtualApple @ 2.50GHz
Benchmark_Restic_Rabin_Next-8                          1        1749383125 ns/op         613.78 MB/s          1301 chunks
Benchmark_Askeladdk_FastCDC_Copy-8                     2         513506770 ns/op        2091.00 MB/s        105327 chunks
Benchmark_Jotfs_FastCDC_Next-8                         3         434035306 ns/op        2473.86 MB/s          1725 chunks
Benchmark_Tigerwill90_FastCDC_Split-8                  3         344989056 ns/op        3112.39 MB/s          2013 chunks
Benchmark_Mhofmann_FastCDC_Next-8                      2         516671625 ns/op        2078.19 MB/s          1718 chunks
Benchmark_PlakarLabs_FastCDC_Copy-8                    8         138843406 ns/op        7733.47 MB/s          3647 chunks
Benchmark_PlakarLabs_FastCDC_Split-8                   8         131869604 ns/op        8142.45 MB/s          3647 chunks
Benchmark_PlakarLabs_FastCDC_Next-8                    8         131754844 ns/op        8149.54 MB/s          3647 chunks
Benchmark_PlakarLabs_UltraCDC_Copy-8                  15          75377942 ns/op        14244.78 MB/s         4096 chunks
Benchmark_PlakarLabs_UltraCDC_Split-8                 15          79355653 ns/op        13530.75 MB/s         4096 chunks
Benchmark_PlakarLabs_UltraCDC_Next-8                  15          74150153 ns/op        14480.64 MB/s         4096 chunks
Benchmark_PlakarLabs_JC_Copy-8                        14          79943033 ns/op        13431.34 MB/s         4033 chunks
Benchmark_PlakarLabs_JC_Split-8                       14          78178872 ns/op        13734.42 MB/s         4033 chunks
Benchmark_PlakarLabs_JC_Next-8                        14          78148342 ns/op        13739.79 MB/s         4033 chunks
PASS
ok      github.com/PlakarLabs/go-cdc-chunkers/tests     75.089s
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
  - [Xiaozhong Jin, Haikun Liu, Chencheng Ye, Xiaofei Liao, Hai Jin and Yu Zhang "Accelerating Content-Defined Chunking for Data Deduplication Based on Speculative Jump" IEEE TRANSACTIONS ON PARALLEL AND DISTRIBUTED SYSTEMS, VOL. 34, NO. 9, SEPTEMBER 2023](https://ieeexplore.ieee.org/stamp/stamp.jsp?tp=&arnumber=10168293)
  