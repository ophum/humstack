package entity

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	Gigabyte = 1024 * 1024 * 1024
	Megabyte = 1024 * 1024
	Kilobyte = 1024
)

type ByteUnit struct {
	bytes int64
}

func NewByteUnit(bytes int64) *ByteUnit {
	return &ByteUnit{bytes}
}

func ParseByteUnit(v string) (*ByteUnit, error) {
	if b, err := regexp.Match("^[0-9]*G$", []byte(v)); err != nil {
		return nil, err
	} else if b {
		n, err := strconv.ParseInt(v[:len(v)-1], 10, 64)
		return &ByteUnit{n * Gigabyte}, err
	}

	if b, err := regexp.Match("^[0-9]*M$", []byte(v)); err != nil {
		return nil, err
	} else if b {
		n, err := strconv.ParseInt(v[:len(v)-1], 10, 64)
		return &ByteUnit{n * Megabyte}, err
	}

	if b, err := regexp.Match("^[0-9]*K$", []byte(v)); err != nil {
		return nil, err
	} else if b {
		n, err := strconv.ParseInt(v[:len(v)-1], 10, 64)
		return &ByteUnit{n * Kilobyte}, err
	}

	if b, err := regexp.Match("^[0-9]*$", []byte(v)); err != nil {
		return nil, err
	} else if b {
		n, err := strconv.ParseInt(v, 10, 64)
		return &ByteUnit{n}, err
	}

	return nil, errors.New("invalid argument")
}

func (s ByteUnit) Int64() int64 {
	return s.bytes
}

func (s ByteUnit) String() string {
	if s.bytes%Gigabyte == 0 {
		return fmt.Sprintf("%dG", s.bytes/Gigabyte)
	} else if s.bytes%Megabyte == 0 {
		return fmt.Sprintf("%dM", s.bytes/Megabyte)
	} else if s.bytes%Kilobyte == 0 {
		return fmt.Sprintf("%dK", s.bytes/Kilobyte)
	}
	return fmt.Sprint(s.bytes)
}
