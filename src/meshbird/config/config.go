package config

type Config struct {
	Key              string `default:"hello-world"`
	SeedAddrs        string `default:"dc1/10.0.0.1/16,dc2/10.0.0.2"`
	LocalAddr        string `default:"10.0.0.1"`
	LocalPrivateAddr string `default:"192.168.0.1"`
	Dc               string `default:"dc1"`
	TransportThreads int    `default:"1"`
	Ip               string `default:"10.237.0.1/16"`
	Mtu              int    `default:"9000"`
	Verbose          int    `default:"0"`
}
