# limits

[![Build Status](https://travis-ci.org/gozips/limits.svg?branch=master)](https://travis-ci.org/gozips/limits)
[![GoDoc](https://godoc.org/github.com/gozips/limits?status.svg)](http://godoc.org/github.com/gozips/limits)

Apply limits to gozip sources

## Examples

Limit number of

    zips := gozips.New(limits.Count(3, sources.FS))
    zips.Add("file1.txt")
    zips.Add("file2.txt")
    zips.Add("file3.txt")
    zips.Add("file4.txt") // => This will return as an exceeded limit error
    zips.Add("file5.txt") // => This will return as an exceeded limit error

---

Limit read size

    zips := gozips.New(limits.Size(256, sources.FS))
    zips.Add("1024.txt") // => will read up to 256 bytes
    zips.Add("1024.txt") // => will read up to 256 bytes
    zips.Add("128.txt")  // => will read to EOF

## License

MIT
