package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	v2 "github.com/lsongdev/id3-go/v2"
)

func main() {
	err := filepath.Walk("/Volumes/data/Music", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".mp3") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			log.Println(err)
			return nil
		}
		tag, err := v2.Read(f)
		if err != nil {
			return err
		}
		log.Println(tag.Title())

		f.Close()
		return nil
	})

	if err != nil {
		panic(err)
	}
}
