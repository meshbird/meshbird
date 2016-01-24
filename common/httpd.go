package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

type HttpService struct {
	BaseService

	localnode *LocalNode
	iface     *InterfaceService
	logger    *log.Logger
}

type Response struct {
	IfaceName   string `json:"iface"`
	LocalIPAddr string `json:"local_ip_addr"`
}

func (hs *HttpService) Name() string {
	return "http-service"
}

func (hs *HttpService) Init(ln *LocalNode) (err error) {
	hs.logger = log.New(os.Stderr, "[httpd] ", log.LstdFlags)
	hs.iface = ln.Service("iface").(*InterfaceService)
	hs.localnode = ln
	return nil
}

func (hs *HttpService) Run() error {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		iName := hs.iface.instance.Name()
		ipAddr := hs.localnode.State().PrivateIP.String()
		fmt.Println(iName, ipAddr)
		resp := Response{iName, ipAddr}
		data, err := json.Marshal(resp)
		fmt.Println(hs.iface.instance.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	http.ListenAndServe(":15080", nil)
	return nil
}
