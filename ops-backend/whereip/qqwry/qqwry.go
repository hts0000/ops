package qqwry

import (
	"bytes"
	"encoding/binary"
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

type MODE byte

const (
	R_MODE1 MODE = 0x01
	R_MODE2 MODE = 0x02
)

type Decoder interface {
	String(string) (string, error)
}

type DBReader struct {
	data []byte // db file data

	head     []byte // db head
	firstIdx uint32 // first index absolute position
	lastIdx  uint32 // last index absolute position

	cursor uint32 // index cursor

	version uint32

	decoder Decoder
}

func NewDBReader(data []byte, decoder Decoder) *DBReader {
	if !IsDBData(data) {
		return nil
	}
	firstIdx := binary.LittleEndian.Uint32(data[:4])
	lastIdx := binary.LittleEndian.Uint32(data[4:])
	return &DBReader{
		data:     data,
		head:     data[0:8],
		firstIdx: firstIdx,
		lastIdx:  lastIdx,
		cursor:   firstIdx,
		decoder:  decoder,
	}
}

func (r *DBReader) FirstIndex() uint32 {
	return r.firstIdx
}

func (r *DBReader) LastIndex() uint32 {
	return r.lastIdx
}

// 获取游标所在索引的绝对定位
func (r *DBReader) CurrentIndex() uint32 {
	return r.cursor
}

// 获取游标所在索引的偏移量
func (r *DBReader) CurrentOffset() uint32 {
	// +4 跳过begin ip
	return r.ReadOffset(r.cursor + 4)
}

// 获取游标所在索引的记录模式
func (r *DBReader) CurrentMode() MODE {
	// +4 跳过end ip
	return r.ReadMode(r.CurrentOffset() + 4)
}

func (r *DBReader) CurrnetIPRange() (net.IP, net.IP) {
	return r.ReadIP(r.cursor), r.ReadIP(r.CurrentOffset())
}

// 获取游标所在索引的第1部分记录
func (r *DBReader) CurrnetPart1() string {
	// +4 跳过end ip
	pos := r.CurrentOffset() + 4
	mod := r.CurrentMode()
	for mod == R_MODE1 || mod == R_MODE2 {
		// +1 跳过mode
		pos = r.ReadOffset(pos + 1)
		mod = r.ReadMode(pos)
	}
	// 获取第1部分记录，读取到0为止
	i := bytes.IndexByte(r.data[pos:], 0)
	part1 := string(r.data[pos : pos+uint32(i)])

	// 转换成GBK编码
	part1, err := r.decoder.String(part1)
	if err != nil {
		return ""
	}

	return part1
}

// 获取游标所在索引的第2部分记录
func (r *DBReader) CurrentPart2() string {
	// +4 跳过end ip
	pos := r.CurrentOffset() + 4
	mod := r.CurrentMode()
	for mod == R_MODE1 {
		// +1 跳过mode
		pos = r.ReadOffset(pos + 1)
		mod = r.ReadMode(pos)
	}
	if mod == R_MODE2 {
		// +1 跳过mode
		// +3 跳过part1的offset
		pos = pos + 1 + 3
	} else {
		// 跳过第1部分记录
		// +1 跳过结束符0
		pos = pos + uint32(bytes.IndexByte(r.data[pos:], 0)+1)
	}

	// 获取第2部分记录，读取到0为止
	i := bytes.IndexByte(r.data[pos:], 0)
	part2 := string(r.data[pos : pos+uint32(i)])

	// 转换成GBK编码
	part2, err := r.decoder.String(part2)
	if err != nil {
		return ""
	}

	return part2
}

func (r *DBReader) HasNextIndex() bool {
	return r.cursor <= r.lastIdx
}

func (r *DBReader) NextIndex() uint32 {
	v := r.cursor
	r.cursor += 7
	return v
}

func (r *DBReader) ReadMode(position uint32) MODE {
	return MODE(r.data[position])
}

func (r *DBReader) ReadOffset(position uint32) uint32 {
	return uint32(r.data[position]) | uint32(r.data[position+1])<<8 | uint32(r.data[position+2])<<16
}

func (r *DBReader) ReadIP(position uint32) net.IP {
	b := r.data[position : position+4]
	return net.IPv4(b[3], b[2], b[1], b[0])
}

func (r *DBReader) ResetCursor() {
	r.cursor = r.firstIdx
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

	j, c := 0, 10
	decoder := simplifiedchinese.GBK.NewDecoder()
	db := NewDBReader(data, decoder)
	for ; db.HasNextIndex() && j < c; db.NextIndex() {
		beginIP, endIP := db.CurrnetIPRange()
		fmt.Printf("currnet index: %v, begin ip: %v, end ip: %v, mod: %v, offset: %v, part1: %v, part2: %v\n",
			db.CurrentIndex(), beginIP, endIP, db.CurrentMode(), db.CurrentOffset(), db.CurrnetPart1(), db.CurrentPart2())

		j++
	}

	first := binary.LittleEndian.Uint32(data[:4])
	last := binary.LittleEndian.Uint32(data[4:])
	fmt.Printf("first index: %d, last index: %d\n", first, last)

	var findRecordOffset func(countryIndex, cityIndex uint32) (uint32, uint32)
	findRecordOffset = func(countryIndex, cityIndex uint32) (uint32, uint32) {
		mod := data[countryIndex]
		fmt.Printf("mod: %v, countryOffset: %v, cityOffset: %v\n", mod, countryIndex, cityIndex)
		switch mod {
		case 0x01:
			countryIndex = LittleEndianByte3ToUint32(data[countryIndex+1 : countryIndex+4])
			countryIndex, cityIndex = findRecordOffset(countryIndex, countryIndex)
		case 0x02:
			countryIndex = LittleEndianByte3ToUint32(data[countryIndex+1 : countryIndex+4])
			cityIndex += 4
			countryIndex, _ = findRecordOffset(countryIndex, cityIndex)
		default:
			// +1 skip country end char 0
			i := bytes.IndexByte(data[countryIndex:], 0) + 1
			cityIndex = countryIndex + uint32(i)
		}
		return countryIndex, cityIndex
	}

	i, cnt := 0, 10
	for first <= last && i < cnt {
		index := data[first : first+7]
		offset := LittleEndianByte3ToUint32(index[4:])
		beginIP := index[:4]
		endIP := data[offset : offset+4]
		mod := data[offset+4]
		// offset +4 skip the end ip
		countryIndex, cityIndex := findRecordOffset(offset+4, offset+4)
		countryOffset, cityOffset := countryIndex+uint32(bytes.IndexByte(data[countryIndex:], 0)), cityIndex+uint32(bytes.IndexByte(data[cityIndex:], 0))
		country, city := string(data[countryIndex:countryOffset]), string(data[cityIndex:cityOffset])
		country, err = decoder.String(country)
		if err != nil {
			panic(err)
		}
		city, err := decoder.String(city)
		if err != nil {
			panic(err)
		}
		// country, city := string(countryBytes), string(cityBytes)

		fmt.Printf("first: %v, index: %v, begin ip: %v, end ip: %v, mod: %v, offset: %v, countryIndex: %v, countryOffset: %v, cityIndex: %v, cityOffset: %v, country: %v, city: %v\n", first, len(index), beginIP, endIP, mod, offset, countryIndex, countryOffset, cityIndex, cityOffset, country, city)
		// fmt.Printf("first: %v, index: %v, begin ip: %v, end ip: %v, mod: %v, offset: %v, countryIndex: %v, cityIndex: %v\n", first, len(index), beginIP, endIP, mod, offset, countryIndex, cityIndex)
		// fmt.Printf("first: %v, index: %v, begin ip: %v, end ip: %v, mod: %v, offset: %v, countryIndex: %v, cityIndex: %v, j: %v, k: %v\n", first, len(index), beginIP, endIP, mod, offset, countryIndex, cityIndex, j, k)
		// fmt.Printf("first: %v, index: %v, begin ip: %v, end ip: %v, mod: %v, offset: %v\n", first, len(index), beginIP, endIP, mod, offset)
		first += 7
		i++
	}
}

func IsDBData(data []byte) bool {
	return true
}

func LittleEndianByte3ToUint32(b []byte) uint32 {
	_ = b[2] // early bounds check to guarantee safety of writes below
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}
