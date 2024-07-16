package qqwry

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	downloadURL = `https://gh-release.zu1k.com/HMBSbige/qqwry/qqwry.dat`

	dbFilename = `qqwry.dat`
)

func GetDBFile() {
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

func ParseDBFile() {
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

	decoder := simplifiedchinese.GBK.NewDecoder()
	db, err := NewDBReader(data, decoder)
	if err != nil {
		panic(err)
	}

	fmt.Printf("first index: %v, last index: %v, db version: %v\n", db.FirstIndex(), db.LastIndex(), db.Version())

	cases := []struct {
		name string
		ip   net.IP
	}{
		{
			name: "normal case",
			ip:   net.IPv4(71, 91, 219, 120),
		},
		{
			name: "last index",
			ip:   net.IPv4(255, 255, 255, 100),
		},
		{
			name: "first index",
			ip:   net.IPv4(0, 0, 0, 1),
		},
	}

	for _, cc := range cases {
		fmt.Printf("case: %v, ip: %v, record: %v\n", cc.name, cc.ip, db.FindRecord(cc.ip))
	}

	for ; db.HasNextIndex(); db.NextIndex() {
		beginIP, endIP := db.CurrnetIPRange()
		fmt.Printf("currnet index: %v, begin ip: %v, end ip: %v, mod: %v, offset: %v, part1: %v, part2: %v\n",
			db.CurrentIndex(), beginIP, endIP, db.CurrentMode(), db.CurrentOffset(), db.CurrnetPart1(), db.CurrentPart2())
	}
}
