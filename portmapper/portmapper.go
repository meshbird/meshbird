package portmapper

import "sync"

type PortMapper struct {
	mutex sync.Mutex
	pairs []int
}

func (pm *PortMapper) AddPair(localPort, publicPort int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.pairs = append(pm.pairs, []int{localPort, publicPort})
}

func (pm *PortMapper) RemovePair(localPort, publicPort int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if len(pm.pairs) == 0 {
		return
	}
	pairs := []int{}
	for _, pair := range pm.pairs {
		if pair[0] == localPort && pair[1] == publicPort {
			continue
		}
		pairs = append(pairs, pair)
	}
	pm.pairs = pairs
}

func (pm *PortMapper) Start() {
	go pm.run()
}

func (pm *PortMapper) Stop() {

}

func (pm PortMapper)
