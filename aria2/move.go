package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	from := flag.String("from", "./", "from path")
	to := flag.String("to", "", "to path")
	remove := flag.Bool("remove", false, "remove if file is the same and exist")
	flag.Parse()

	sFrom, e := filepath.Abs(*from)
	if e != nil {
		panic(e)
	}

	files := getFiles(sFrom)
	log.Println("from:", *from)
	log.Println("to:", *to)
	if *to == "" {
		panic(e)
	}
	sTo, e := filepath.Abs(*to)
	if e != nil {
		panic(e)
	}
	//os.Mkdir(sTo, os.ModePerm)
	for _, file := range files {
		if !checkFinish(file) {
			log.Println(file, "unfinished")
			continue
		}

		info, e := os.Stat(file)
		if e != nil {
			log.Fatal(e)
		}
		_, toFile := filepath.Split(file)
		log.Println(toFile, "move")
		if info.IsDir() {
			toPath := filepath.Join(sTo, toFile)
			//_ = os.MkdirAll(toPath, os.ModePerm)
			moves := getFiles(file)
			for _, subFile := range moves {
				_, toSubFile := filepath.Split(subFile)
				e = moveFile(subFile, toPath, toSubFile, *remove)
				if e != nil {
					fmt.Println(e)
					continue
				}
			}
		} else {
			e = moveFile(file, sTo, toFile, *remove)
			if e != nil {
				fmt.Println(e)
				continue
			}
		}
	}
	log.Println("all finished")
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
			if filepath.Ext(fullPath) == ".aria2" || filepath.Ext(fullPath) == ".torrent" {
				log.Printf("%s skip\n", fullPath)
				continue
			}
			files = append(files, fullPath)
		}
	}

	return files
}
func moveFile(sourcePath, toPath, destFile string, remove bool) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	//ignore error
	_ = os.MkdirAll(toPath, os.ModePerm)
	dest := filepath.Join(toPath, destFile)
	info, err := os.Stat(dest)
	if !os.IsNotExist(err) {
		log.Println(dest, "exist")
		if remove {
			sourceInfo, err := inputFile.Stat()
			if err != nil {
				return err
			}
			if info.Size() == sourceInfo.Size() {
				log.Println("remove:", sourcePath)
				return os.Remove(sourcePath)
			}
			log.Println("skip remove:", sourcePath)
		}

		return nil
	}

	err = os.Rename(sourcePath, dest)
	if err != nil {
		log.Println("not same disk:", sourcePath, dest)
		log.Println(err)

	} else {
		log.Println("same disk:", sourcePath, dest)
		return nil
	}
	outputFile, err := os.Create(dest)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

func checkFinish(path string) bool {
	_, e := os.Open(path + ".aria2")
	if !os.IsNotExist(e) {
		return false
	}
	return true
}
