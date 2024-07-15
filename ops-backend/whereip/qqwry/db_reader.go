package qqwry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
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

	return r, nil
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

// 获取游标所在索引的IP区间
func (r *DBReader) CurrnetIPRange() (beginIP net.IP, endIP net.IP) {
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

func IsDBData(data []byte) bool {
	return true
}
