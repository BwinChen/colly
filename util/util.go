package util

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0"
	IP        = "192.168.0.6"
)

func Deadline(duration string) int64 {
	n := time.Now()
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatal(err)
	}
	n = n.Add(d)
	return n.Unix()
}

func ConvertSize(s string) (int64, error) {
	s = strings.ToLower(s)
	ss := strings.Split(s, " ")
	f, err := strconv.ParseFloat(ss[0], 64)
	if err != nil {
		log.Println(err)
		return 0, errors.New("文件大小转换错误")
	}
	if ss[1] == "bytes" || ss[1] == "b" {
		return int64(f), nil
	} else if ss[1] == "kib" || ss[1] == "kb" {
		return int64(f * 1024), nil
	} else if ss[1] == "mib" || ss[1] == "mb" {
		return int64(f * 1024 * 1024), nil
	} else if ss[1] == "gib" || ss[1] == "gb" {
		return int64(f * 1024 * 1024 * 1024), nil
	} else {
		return int64(f * 1024 * 1024 * 1024 * 1024), nil
	}
}

func InfoHash(magnet string) string {
	return magnet[strings.Index(magnet, "btih:")+5 : strings.Index(magnet, "&dn=")]
}
