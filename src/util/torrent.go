package util

import (
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"log"
	"os"
	"time"
)

func ParseTorrent(path string) (*Torrent, error) {
	file, err := os.Open(path)
	if err != nil {
		return &Torrent{}, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("Error closing .torrent file: %v", err)
		}
	}(file)
	var mi *metainfo.MetaInfo
	mi, err = metainfo.Load(file)
	if err != nil {
		return &Torrent{}, err
	}
	var t Torrent
	t.InfoHash = mi.HashInfoBytes().String()
	if mi.CreationDate > 0 {
		t.CreationDate = time.Unix(mi.CreationDate, 0).Format("2006-01-02 15:04:05")
	} else {
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Printf("Stat error: %v", err)
		}
		t.CreationDate = fileInfo.ModTime().Format("2006-01-02 15:04:05")
	}
	var info map[string]interface{}
	err = bencode.Unmarshal(mi.InfoBytes, &info)
	if err != nil {
		return &Torrent{}, err
	}
	t.Name = info["name"].(string)
	v, e := info["files"]
	if e {
		// 多文件torrent
		files := v.([]interface{})
		for _, file := range files {
			var f File
			f.Length = file.(map[string]interface{})["length"].(int64)
			t.Length += f.Length
			//f.Path = "/"
			path := file.(map[string]interface{})["path"].([]interface{})
			for i, p := range path {
				f.Path += p.(string)
				if i != len(path)-1 {
					f.Path += "/"
				}
			}
			t.Files = append(t.Files, f)
		}
		t.FileNumber = len(files)
	} else {
		// 单文件torrent
		var f File
		f.Length = info["length"].(int64)
		f.Path = info["name"].(string)
		t.Length = f.Length
		t.Files = append(t.Files, f)
		t.FileNumber = 1
	}
	return &t, nil
}
