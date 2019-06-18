package common

type Config struct {
	Key              string `default:""`
	HostAddr         string
	PublicAddrs      []string `arg:"-p,separate"`
	BindAddrs        []string `default:"" arg:"-b,separate"`
	SeedAddrs        []string `default:"" arg:"-s,separate"`
	TransportThreads int      `default:"1"`
	Ip               string   `default:"10.237.0.1/16" arg:"required"`
	Mtu              int      `default:"9000"`
	Verbose          int      `default:"0"`
}
