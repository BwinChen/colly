package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseTorrent(t *testing.T) {
	torrent, err := ParseTorrent("335a0ed729879050522ee43695feb2cc3063ec53.torrent")
	if err != nil {
		log.Fatalf("Error parsing .torrent file: %v", err)
	}
	log.Println(*torrent)
}

func TestIndexTorrent(t *testing.T) {
	torrent, err := ParseTorrent("3732a1a6bf349fef05af0c1c529531b2226073f4.torrent")
	if err != nil {
		log.Fatalf("Error parsing .torrent file: %v", err)
	}
	id, err := IndexTorrent(torrent)
	if err != nil {
		log.Fatalf("Error indexing torrent: %v", err)
	}
	log.Printf("ES Torrent id: %s\n", id)
}

func TestSearchByInfoHash(t *testing.T) {
	h, err := SearchByInfoHash("3732a1a6bf349fef05af0c1c529531b2226073f4")
	if err != nil {
		log.Fatalf("Error searching by infohash: %v", err)
	}
	log.Println(h)
}

func TestDeleteByInfoHash(t *testing.T) {
	h, err := DeleteByInfoHash("3732a1a6bf349fef05af0c1c529531b2226073f4")
	if err != nil {
		log.Fatalf("Error deleting by infohash: %v", err)
	}
	log.Println(h)
}

func TestWalk(t *testing.T) {
	root := "D:\\McAfee\\Desktop\\torrents"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("WalkFunc error: %v\n", err)
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".torrent") {
			// 删除文件
			defer func(path string) {
				err := os.Remove(path)
				if err != nil {
					log.Printf("Remove Error: %v\n", err)
					return
				}
				_ = os.Remove(filepath.Dir(path))
			}(path)
			// 解析torrent
			torrent, err := ParseTorrent(path)
			if err != nil {
				log.Printf("ParseTorrent error: %v\n", err)
				return nil
			}
			// es去重
			hit, err := SearchByInfoHash(torrent.InfoHash)
			if err != nil {
				log.Printf("SearchByInfoHash error: %v\n", err)
				return nil
			}
			if hit > 0 {
				return nil
			}
			// 索引torrent
			id, err := IndexTorrent(torrent)
			if err != nil {
				log.Printf("IndexTorrent error: %v\n", err)
				return nil
			}
			log.Printf("ES Torrent id: %s\n", id)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Walk error: %v\n", err)
	}
}
