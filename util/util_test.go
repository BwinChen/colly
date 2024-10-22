package util

import (
	"fmt"
	"log"
	"testing"
)

func TestParseTorrent(t *testing.T) {
	torrent, err := ParseTorrent("3732a1a6bf349fef05af0c1c529531b2226073f4.torrent")
	if err != nil {
		log.Fatalf("Error parsing .torrent file: %v", err)
	}
	fmt.Println(*torrent)
}

func TestIndexTorrent(t *testing.T) {
	torrent, err := ParseTorrent("3732a1a6bf349fef05af0c1c529531b2226073f4.torrent")
	if err != nil {
		log.Fatalf("Error parsing .torrent file: %v", err)
	}
	err = IndexTorrent(*torrent)
	if err != nil {
		log.Fatalf("Error indexing torrent: %v", err)
	}
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
