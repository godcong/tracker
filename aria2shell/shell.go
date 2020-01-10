package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/zyxar/argo/rpc"
)

func main() {
	rpcSecre := flag.String("secret", "", "set --rpc-secret for aria2c")
	rpcURI := flag.String("uri", "http://localhost:6800/jsonrpc", "set rpc address")
	filepath := flag.String("path", ".", "get torrent files")
	move := flag.String("to", "./success", "move file when success")
	flag.Parse()


	if e != nil {
		log.Fatal(e)
	}
	files := getFiles(*filepath)

	for _, file := range files {
		log.Println("add file:", file)
		gid, e := client.AddTorrent(file)
		if e != nil {
			log.Println(e)
			continue
		}
		log.Println("to", *move)
		log.Println("gid", gid)
	}
	log.Println("all was done")
}

func getFiles(path string) (files []string) {
	info, e := os.Stat(path)
	if e != nil {
		return nil
	}

	if info.IsDir() {
		file, e := os.Open(path)
		if e != nil {
			return nil
		}
		defer file.Close()
		names, e := file.Readdirnames(-1)
		if e != nil {
			return nil
		}
		var fullPath string
		for _, name := range names {
			fullPath = filepath.Join(path, name)
			if filepath.Ext(fullPath) != ".torrent" {
				log.Println("path", fullPath, "skip")
				continue
			}
			files = append(files, fullPath)
		}
	}

	return files
}
