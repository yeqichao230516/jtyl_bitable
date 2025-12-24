package service

import (
	"errors"
	"fmt"
	"strings"
)

const (
	keyRCP = uint32(0xDEADBEEF)
	keyRCS = uint32(0xCAFEBABE)
	modP   = 1_000_003
)

// ---------------- 内部工具 ----------------
func f(num, key uint32) uint32 {
	const prime = 0x9E3779B9
	return (num*prime + key) ^ (key >> 16)
}

func feistel(x, key uint32) uint32 {
	l, r := x/1000, x%1000
	for i := uint32(0); i < 4; i++ {
		l, r = r, l^f(r, key+i)
	}
	return r*1000 + l
}

func transform(num uint32, key uint32) uint32 {
	x := feistel(num, key) % modP
	if x >= 1_000_000 { // 把 [1M..1M+2] 再压回 0–999999
		x -= 1_000_000
	}
	return x
}

func parse(code string) (brand, line string, num uint32, err error) {
	parts := strings.Split(code, "-")
	if len(parts) != 3 {
		return "", "", 0, errors.New("invalid format")
	}
	var n uint32
	if _, e := fmt.Sscanf(parts[2], "%06d", &n); e != nil {
		return "", "", 0, e
	}
	return parts[0], parts[1], n, nil
}

// EncodeRM 输入 "RM-FL-000001" 返回 "RCP-FL-xxxxxx" 和 "RCS-FL-xxxxxx"
func EncodeRM(code string) (rcp, rcs string, err error) {
	brand, line, num, err := parse(code)
	if err != nil || brand != "RM" {
		return "", "", errors.New("input must be RM-**-******")
	}
	rcpNum := transform(num, keyRCP)
	rcsNum := transform(num, keyRCS)
	return fmt.Sprintf("RCP-%s-%06d", line, rcpNum),
		fmt.Sprintf("RCS-%s-%06d", line, rcsNum), nil
}

// DecodeRCP 输入 "RCP-BZ-399304" 返回 "RM-BZ-000001"
func DecodeRCP(code string) (string, error) {
	brand, line, num, err := parse(code)
	if err != nil || brand != "RCP" {
		return "", errors.New("input must be RCP-**-******")
	}
	orig := transform(num, keyRCP)
	return fmt.Sprintf("RM-%s-%06d", line, orig), nil
}

// DecodeRCS 同理
func DecodeRCS(code string) (string, error) {
	brand, line, num, err := parse(code)
	if err != nil || brand != "RCS" {
		return "", errors.New("input must be RCS-**-******")
	}
	orig := transform(num, keyRCS)
	return fmt.Sprintf("RM-%s-%06d", line, orig), nil
}
