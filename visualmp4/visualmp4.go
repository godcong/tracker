package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	from := flag.String("from", "video.txt", "set the from file path")
	output := flag.String("output", "video", "set the file output path")
	flag.Parse()

	file, e := os.Open(*from)
	if e != nil {
		return
	}

	reader := bufio.NewReader(file)
ReadEnd:
	for {
		line, _, e := reader.ReadLine()
		if e == io.EOF {
			break ReadEnd
		}
		_, _ = os.Create(filepath.Join(*output, string(line)) + ".mp4")
	}
	fmt.Println("done")

}
