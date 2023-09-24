# go-cdc-chunkers

go-cdc-chunkers is a Golang package that provides a unified access to multiple CDC chunking algorithms.

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


## Reference

  - [Xia, Wen, et al. "Fastcdc: a fast and efficient content-defined chunking approach for data deduplication." 2016 USENIX Annual Technical Conference](https://www.usenix.org/system/files/conference/atc16/atc16-paper-xia.pdf)
  - [Zhou, Wang, Xia, Zhang "UltraCDC:A Fast and Stable Content-Defined Chunking Algorithm for Deduplication-based Backup Storage Systems" 2022 IEEE](https://ieeexplore.ieee.org/document/9894295)