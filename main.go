package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
)

func main() {
	src := flag.String("src", "tracker.txt", "set the file source")
	target := flag.String("target", "ouput.txt", "set the output filename")
	flag.Parse()

	file, e := os.Open(*src)
	if e != nil {
		panic(e)
	}
	defer file.Close()
	results := make(map[string][]byte)

	bio := bufio.NewReader(file)
output:
	for {
		line, _, e := bio.ReadLine()
		if e != nil {
			if e == io.EOF {
				break output
			}
			panic(e)
		}
		line = bytes.TrimSpace(line)
		if filterProtocol(string(line)) {
			results[string(line)] = nil
		}
	}

	outFile, e := os.OpenFile(*target, os.O_CREATE|os.O_SYNC|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if e != nil {
		panic(e)
	}
	defer outFile.Close()
	for key := range results {
		outFile.WriteString(key)
		outFile.WriteString("\n")
	}

}

var protocals = []string{"http", "udp", "wss", "tcp"}

func filterProtocol(line string) bool {
	for _, v := range protocals {
		if strings.Index(line, v) == 0 {
			return true
		}
	}
	return false
}
