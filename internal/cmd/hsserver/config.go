package hsserver

type config struct {
	http struct {
		port int
	}
	db struct {
		dsn string
	}
}

type grpcConfig struct {
	config
}

type httpConfig struct {
	config
}
