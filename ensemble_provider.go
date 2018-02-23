package curator

type EnsembleProvider interface {
	Start() error
	Close() error
	GetConnectionString() string
}
