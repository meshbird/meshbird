package common

type Service interface {
	Init(*LocalNode) error
	Name() string
	Run() error
	Stop()
}