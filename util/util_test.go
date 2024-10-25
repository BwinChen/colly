package util

import (
	"log"
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

func TestSAdd(t *testing.T) {
	r, err := SAdd("4192944")
	if err != nil {
		return
	}
	log.Println(r)
}

func TestSIsMember(t *testing.T) {
	r, err := SIsMember("4192944")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println(r)
}

func TestSRem(t *testing.T) {
	r, err := SRem("4192944")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println(r)
}
