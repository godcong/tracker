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
	src := flag.String("src", "tracker_source.txt", "set the file source")
	target := flag.String("target", "tracker.txt", "set the output filename")
	aria := flag.Bool("aria", true, "set the bool to open aria format")
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

	os.Remove(*target)

	outFile, e := os.OpenFile(*target, os.O_CREATE|os.O_SYNC|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if e != nil {
		panic(e)
	}
	defer outFile.Close()
	for key := range results {
		outFile.WriteString(key)
		if *aria {
			outFile.WriteString(",")
		} else {
			outFile.WriteString("\n")
		}
	}

}

var supportProtocols = []string{"http", "udp", "wss", "tcp"}

func filterProtocol(line string) bool {
	for _, v := range supportProtocols {
		if strings.Index(line, v) == 0 {
			return true
		}
	}
	return false
}
