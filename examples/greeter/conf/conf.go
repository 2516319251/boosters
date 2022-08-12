package conf

type Config struct {
	Server *Server `yaml:"server"`
}

type Server struct {
	Http *Http `yaml:"http"`
	Grpc *Grpc `yaml:"grpc"`
}

type Http struct {
	Addr string `yaml:"addr"`
}

type Grpc struct {
	Addr string `yaml:"addr"`
}
