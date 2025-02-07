package util

import (
	"hash/crc32"
	"strconv"
)

var table = crc32.MakeTable(crc32.IEEE)

func Checksum(s string) string {
	return strconv.FormatUint(uint64(crc32.Checksum([]byte(s), table)), 10)
}
