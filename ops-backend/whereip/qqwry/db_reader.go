package qqwry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

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
	cursor   uint32 // index cursor
	version  string // db file version

	decoder Decoder
}

func NewDBReader(data []byte, decoder Decoder) (*DBReader, error) {
	if !IsDBData(data) {
		return nil, errors.New("invalid db file")
	}

	firstIdx := binary.LittleEndian.Uint32(data[:4])
	lastIdx := binary.LittleEndian.Uint32(data[4:])

	r := &DBReader{
		data:     data,
		head:     data[0:8],
		firstIdx: firstIdx,
		lastIdx:  lastIdx,
		cursor:   firstIdx,
		decoder:  decoder,
	}

	// +4 跳过begin ip
	r.version = r.ReadPart1(r.ReadOffset(lastIdx+4)) + r.ReadPart2(r.ReadOffset(lastIdx+4))

	return r, nil
}

func (r *DBReader) FirstIndex() uint32 {
	return r.firstIdx
}

func (r *DBReader) LastIndex() uint32 {
	return r.lastIdx
}

func (r *DBReader) Version() string {
	return r.version
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

// 获取游标所在索引的IP区间
func (r *DBReader) CurrnetIPRange() (beginIP net.IP, endIP net.IP) {
	return r.ReadIP(r.cursor), r.ReadIP(r.CurrentOffset())
}

// 获取游标所在索引的第1部分记录
func (r *DBReader) CurrnetPart1() string {
	pos := r.CurrentOffset()
	return r.ReadPart1(pos)
}

// 获取游标所在索引的第2部分记录
func (r *DBReader) CurrentPart2() string {
	pos := r.CurrentOffset()
	return r.ReadPart2(pos)
}

func (r *DBReader) HasNextIndex() bool {
	return r.cursor <= r.lastIdx
}

func (r *DBReader) NextIndex() uint32 {
	v := r.cursor
	r.cursor += 7
	return v
}

func (r *DBReader) ReadMode(pos uint32) MODE {
	return MODE(r.data[pos])
}

func (r *DBReader) ReadOffset(pos uint32) uint32 {
	return uint32(r.data[pos]) | uint32(r.data[pos+1])<<8 | uint32(r.data[pos+2])<<16
}

func (r *DBReader) ReadIP(pos uint32) net.IP {
	b := r.data[pos : pos+4]
	return net.IPv4(b[3], b[2], b[1], b[0])
}

func (r *DBReader) ReadPart1(pos uint32) string {
	// +4 跳过end ip
	pos += 4
	mod := r.ReadMode(pos)
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

	return strings.Trim(part1, " ")
}

func (r *DBReader) ReadPart2(pos uint32) string {
	// +4 跳过end ip
	pos += 4
	mod := r.ReadMode(pos)
	for mod == R_MODE1 {
		// +1 跳过mode
		pos = r.ReadOffset(pos + 1)
		mod = r.ReadMode(pos)
	}

	// 判断是否还需要跳转
	redirect := func(pos uint32) uint32 {
		mod := r.ReadMode(pos)
		if mod == R_MODE1 || mod == R_MODE2 {
			// +1 跳过mode
			// 拿到再次跳转后的绝对定位
			pos = r.ReadOffset(pos + 1)
		}
		return pos
	}

	if mod == R_MODE2 {
		// +3 跳过part1的offset
		pos = pos + 3 + 1
		pos = redirect(pos)
	} else {
		// 跳过第1部分记录
		// +1 跳过结束符0
		pos = pos + uint32(bytes.IndexByte(r.data[pos:], 0)) + 1
		pos = redirect(pos)
	}

	// 获取第2部分记录，读取到0为止
	i := bytes.IndexByte(r.data[pos:], 0)
	part2 := string(r.data[pos : pos+uint32(i)])

	// 转换成GBK编码
	part2, err := r.decoder.String(part2)
	if err != nil {
		return ""
	}

	return strings.Trim(part2, " ")
}

func (r *DBReader) FindRecord(ip net.IP) string {
	start, end := r.FirstIndex(), r.LastIndex()
	target := binary.LittleEndian.Uint32(ip[12:])

	idxLen := uint32(4 + 3)

	// 在索引区寻找最后一个小于等于target的位置，返回索引位置
	search := func(start, end, target uint32) uint32 {
		for start+idxLen <= end {
			// 找到[start:end]中间的那条index
			mid := (end-start)/idxLen/2*idxLen + start

			// 读取中间index的ip，将其转换成uint32
			ip := r.ReadIP(mid)
			ipc := binary.BigEndian.Uint32(ip[12:])

			if ipc > target {
				end = mid - idxLen
			} else {
				start = mid + idxLen
			}
		}
		return start
	}

	targetIdx := search(start, end, target)
	if targetIdx >= r.LastIndex() {
		fmt.Printf("cannot find ip: %v record\n", ip)
		return ""
	}

	pos := r.ReadOffset(targetIdx + 4)
	return r.ReadPart1(pos) + r.ReadPart2(pos)
}

func (r *DBReader) ResetCursor() {
	r.cursor = r.firstIdx
}

func IsDBData(data []byte) bool {
	return true
}
