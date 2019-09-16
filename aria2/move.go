package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	from := flag.String("from", "./", "from path")
	to := flag.String("to", "", "to path")
	flag.Parse()

	sFrom, e := filepath.Abs(*from)
	if e != nil {
		panic(e)
	}

	files := getFiles(sFrom)
	fmt.Println("file:", files)
	fmt.Println("from:", *from)
	fmt.Println("to:", *to)
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
			fmt.Println("unfinished:", file)
			continue
		}

		info, e := os.Stat(file)
		if e != nil {
			fmt.Println("error:", e)
			continue
		}
		_, toFile := filepath.Split(file)
		//last := len(toFile) - 1
		fmt.Println("to:", toFile)
		if info.IsDir() {
			toPath := filepath.Join(sTo, toFile)
			_ = os.MkdirAll(toPath, os.ModePerm)
			moves := getFiles(file)
			for _, subFile := range moves {
				_, toSubFile := filepath.Split(subFile)
				e = moveFile(subFile, filepath.Join(toPath, toSubFile))
				if e != nil {
					fmt.Println(e)
				}
			}
		} else {
			e = moveFile(file, filepath.Join(sTo, toFile))
			if e != nil {
				fmt.Println(e)
			}
		}
	}

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
				fmt.Printf("skip[%s]\n", fullPath)
				continue
			}
			files = append(files, fullPath)
		}
	}

	return files
}
func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
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
