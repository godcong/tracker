package aria2shell

import (
	"strings"
)

type rpc struct {
	link string
}

func NewRPC(url, port, path string) {
	r := new(rpc)
	r.link = strings.Join([]string{url, ":", port, "/", path}, "")
	log.With("url", r.link).Info("")
}
