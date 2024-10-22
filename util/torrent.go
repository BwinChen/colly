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
			log.Fatalf("Error closing .torrent file: %v", err)
		}
	}(file)
	var mi *metainfo.MetaInfo
	mi, err = metainfo.Load(file)
	if err != nil {
		return &Torrent{}, err
	}
	var torrent Torrent
	torrent.InfoHash = mi.HashInfoBytes().String()
	torrent.CreationDate = time.Unix(mi.CreationDate, 0).Format("2006-01-02 15:04:05")
	var info map[string]interface{}
	err = bencode.Unmarshal(mi.InfoBytes, &info)
	if err != nil {
		return &Torrent{}, err
	}
	torrent.Name = info["name"].(string)
	files := info["files"].([]interface{})
	for _, file := range files {
		var f File
		f.Length = file.(map[string]interface{})["length"].(int64)
		torrent.Length += f.Length
		//f.Path = "/"
		path := file.(map[string]interface{})["path"].([]interface{})
		for i, p := range path {
			f.Path += p.(string)
			if i != len(path)-1 {
				f.Path += "/"
			}
		}
		torrent.Files = append(torrent.Files, f)
	}
	torrent.FileNumber = len(files)
	return &torrent, nil
}
