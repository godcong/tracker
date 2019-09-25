package aria2shell

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type JSONRPC struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      string `json:"id"`
}

type RPC interface {
	Request(jsonrpc *JSONRPC) (e error)
}
type rpc struct {
	link string
}

func (r *rpc) Request(jsonrpc *JSONRPC) (e error) {
	cli := http.DefaultClient
	b, e := json.Marshal(jsonrpc)
	if e != nil {
		return e
	}
	resp, e := cli.Post(r.link, "application/json-rpc", bytes.NewReader(b))
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	all, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}
	log.Info(string(all))
	return nil
}

func NewRPC(url, port, path string) RPC {
	r := new(rpc)
	r.link = strings.Join([]string{url, ":", port, "/", path}, "")
	log.With("url", r.link).Info("")
	return r
}
