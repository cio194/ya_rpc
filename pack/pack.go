package pack

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"net"
)

func NewPack(dataLen uint32) []byte {
	// pack: dataLen data
	pack := make([]byte, 4+dataLen)
	data := pack
	Putuint32(&data, dataLen)
	return pack
}

func ParsePack(conn net.Conn) ([]byte, error) {
	// pack: dataLen data

	// 获取数据长度
	var lenBuf [4]byte
	_, err := io.ReadFull(conn, lenBuf[:])
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(lenBuf[:])

	// 获取数据
	data := make([]byte, dataLen)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func Getuint32(dp *[]byte) (uint32, error) {
	const msg = "get uint32"
	return get32(dp, msg)
}

func Putuint32(dp *[]byte, u uint32) {
	put32(dp, u)
}

func Getfloat64(dp *[]byte) (float64, error) {
	const msg = "get float64"
	u, err := get64(dp, msg)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(u), nil
}

func Putfloat64(dp *[]byte, f float64) {
	put64(dp, math.Float64bits(f))
}

func Getstring(dp *[]byte) (string, error) {
	const msg = "get string"
	strLen, err := get32(dp, msg)
	if err != nil {
		return "", err
	}
	return getStr(dp, strLen, msg)
}

func Putstring(dp *[]byte, s string) {
	put32(dp, uint32(len(s)))
	putStr(dp, s)
}

func Geterror(dp *[]byte) (error, error) {
	const msg = "get error_msg"
	// 获取长度
	errLen, err := get32(dp, msg)
	if err != nil {
		return nil, err
	}
	// 获取错误消息
	if errLen == 0 {
		return nil, nil
	}
	errMsg, err := getStr(dp, errLen, msg)
	if err != nil {
		return nil, err
	}
	return errors.New(errMsg), nil
}

func Puterror(dp *[]byte, err error) {
	// 写入长度
	put32(dp, ErrLen(err))
	// 写入错误消息
	if err == nil {
		return
	}
	putStr(dp, err.Error())
}

func get32(dp *[]byte, msg string) (uint32, error) {
	data := *dp
	if len(data) < 4 {
		return 0, errors.New(msg)
	}
	*dp = data[4:]
	return binary.BigEndian.Uint32(data[:4]), nil
}

func put32(dp *[]byte, u uint32) {
	data := *dp
	binary.BigEndian.PutUint32(data, u)
	*dp = data[4:]
}

func get64(dp *[]byte, msg string) (uint64, error) {
	data := *dp
	if len(data) < 8 {
		return 0, errors.New(msg)
	}
	*dp = data[8:]
	return binary.BigEndian.Uint64(data[:8]), nil
}

func put64(dp *[]byte, u uint64) {
	data := *dp
	binary.BigEndian.PutUint64(data, u)
	*dp = data[8:]
}

func getStr(dp *[]byte, strLen uint32, msg string) (string, error) {
	data := *dp
	if uint32(len(data)) < strLen {
		return "", errors.New(msg)
	}
	*dp = data[strLen:]
	return string(data[:strLen]), nil
}

func putStr(dp *[]byte, s string) {
	data := *dp
	copy(data, s)
	*dp = data[len(s):]
}

func ErrLen(err error) uint32 {
	if err == nil {
		return 0
	} else {
		return uint32(len(err.Error()))
	}
}
