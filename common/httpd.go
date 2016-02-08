package common

import (
	"encoding/json"
	"net/http"
	_ "net/http/pprof"
)

type HttpService struct {
	BaseService

	localnode *LocalNode
	iface     *InterfaceService
}

type Response struct {
	IfaceName   string      `json:"iface"`
	LocalIPAddr string      `json:"local_ip_addr"`
	Peers       interface{} `json:"net_peers"`
}

func (hs *HttpService) Name() string {
	return "http-service"
}

func (hs *HttpService) Init(ln *LocalNode) (err error) {
	hs.localnode = ln
	hs.iface = ln.Service("iface").(*InterfaceService)
	return nil
}

func (hs *HttpService) Run() error {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		iName := hs.iface.instance.Name()
		ipAddr := hs.localnode.State().PrivateIP().String()
		peers := hs.iface.netTable.PeerAddresses()
		resp := Response{iName, ipAddr, peers}
		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			return
		}
	})
	http.ListenAndServe(":15080", nil)
	return nil
}
