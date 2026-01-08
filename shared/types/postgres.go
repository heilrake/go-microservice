package types

type PostgresConfig struct {
	DSN         string
	MaxConns    int32
	MinConns    int32
	MaxIdleTime string // example: "5m"
}
