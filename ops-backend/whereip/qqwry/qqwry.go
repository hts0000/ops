package qqwry

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	downloadURL = `https://gh-release.zu1k.com/HMBSbige/qqwry/qqwry.dat`

	dbFilename = `qqwry.dat`
)

func getDBFile() {
	client := http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout:   5 * time.Second,
			IdleConnTimeout:       10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 20 * time.Second,
			Proxy:                 http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error code: %v\n", resp.StatusCode)
		panic("http status error")
	}
	defer resp.Body.Close()

	fp, err := os.OpenFile(dbFilename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("create file error: %v\n", err)
		panic("create file error")
	}
	defer fp.Close()

	nn, err := io.Copy(fp, resp.Body)
	if err != nil {
		fmt.Printf("write file error: %v, written byte: %v\n", err, nn)
		panic("write file error")
	}

	fmt.Println("get qqwry db file success")
}

func parseDBFile() {
	fp, err := os.Open(dbFilename)
	if err != nil {
		fmt.Printf("open db file error: %v\n", err)
		panic("open db file error")
	}
	defer fp.Close()

	info, err := fp.Stat()
	if err != nil {
		fmt.Printf("get db file stat info error: %v\n", err)
		panic("get db file stat info error")
	}

	fmt.Printf("db file size is: %v\n", info.Size())

	data := make([]byte, info.Size())
	_, err = fp.Read(data)
	if err != nil {
		fmt.Printf("read head error: %v\n", err)
		panic("read head error")
	}

	fmt.Printf("data size: %v\n", len(data))

	first := binary.LittleEndian.Uint32(data[:4])
	last := binary.LittleEndian.Uint32(data[4:])
	fmt.Printf("first index: %d, last index: %d\n", first, last)

	for first <= last {
		index := data[first : first+7]
		ip := index[:4]
		offset := byte3ToUint32(index[4:])
		fmt.Printf("first: %v, index: %v, ip: %v, offset: %v\n", first, len(index), ip, offset)
		first += 7
	}
}

func checkDBFile(n int) bool {
	return n == 8
}

func byte3ToUint32(b []byte) uint32 {
	_ = b[2] // early bounds check to guarantee safety of writes below
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}
