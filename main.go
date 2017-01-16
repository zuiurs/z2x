package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()

	vmlinuz, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	defer vmlinuz.Close()

	var size int64
	if info, err := vmlinuz.Stat(); err == nil {
		// exit for file which 1GB or more sizefile
		if info.Size() > 1e9 {
			fmt.Fprintf(os.Stderr, "too large file: %d\n", size)
			os.Exit(1)
		} else {
			size = info.Size()
		}
	}

	gzMagic := []byte{0x1f, 0x8b, 0x08, 0x00}
	hdr := 0
	offset := -1

	data, err := ioutil.ReadAll(vmlinuz)
	for i, v := range data {
		if hdr == len(gzMagic) {
			offset = i - len(gzMagic)
			break
		}

		if v == gzMagic[hdr] {
			hdr++
		} else {
			hdr = 0
		}
	}

	if offset < 0 {
		fmt.Fprint(os.Stderr, "this is not vmlinuz format\n")
		os.Exit(1)
	}

	vmlinux, err := os.Create("vmlinux")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	defer vmlinux.Close()

	buf := bytes.NewBuffer(data[offset:])
	r, err := gzip.NewReader(buf)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	io.Copy(vmlinux, r)
}
