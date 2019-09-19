package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	from := flag.String("from", "video.txt", "set the from file path")
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
		_, _ = os.Create(string(line) + ".mp4")
	}
	fmt.Println("done")

}
