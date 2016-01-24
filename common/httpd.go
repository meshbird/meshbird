package common

import (
	"io"
	"log"
	"net/http"
	"os"
	//	"encoding/json"
)

type HttpService struct {
	BaseService

	localnode *LocalNode
	netTable  *NetTable
	logger    *log.Logger
}

func (hs *HttpService) Name() string {
	return "http-service"
}

func (hs *HttpService) Init(ln *LocalNode) (err error) {
	hs.logger = log.New(os.Stderr, "[httpd] ", log.LstdFlags)
	hs.localnode = ln
	hs.netTable = ln.Service("net-table").(*NetTable)
	return nil
}

func (hs *HttpService) Run() error {
	http.HandleFunc("/stat", getStats)
	http.ListenAndServe(":15080", nil)
	return nil
}

func getStats(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello mazafaka")
}
