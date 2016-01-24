package common

import (
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

type HttpService struct {
	BaseService

	localnode *LocalNode
	netTable  *NetTable
	iface     *InterfaceService
	logger    *log.Logger
}
type Response struct {
	ifaceName   string `json:"iface"`
	localIPAddr string `json:"local_ip_addr"`
}

func (hs *HttpService) Name() string {
	return "http-service"
}

func (hs *HttpService) Init(ln *LocalNode) (err error) {
	hs.logger = log.New(os.Stderr, "[httpd] ", log.LstdFlags)
	hs.iface = ln.Service("iface").(*InterfaceService)
	hs.localnode = ln
	hs.netTable = ln.Service("net-table").(*NetTable)
	return nil
}

func (hs *HttpService) Run() error {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		resp := Response{hs.iface.instance.Name(), hs.localnode.State().PrivateIP.String()}
		data, err := json.Marshal(resp)
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
