package core

type Config struct {

}

type State struct {
	ListenPort int
}

type Node struct {
	config *Config
	state *State
}

func NewNode(cfg *Config) *Node {
	n := new(Node)
	return n
}

func (n *Node) Run() error {
	return nil
}

func (n *Node) Stop() error {
	return nil
}