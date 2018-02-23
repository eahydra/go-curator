package curator

type FixedEnsembleProvider string

func NewFixedEnsembleProvider(connString string) FixedEnsembleProvider {
	return FixedEnsembleProvider(connString)
}

func (FixedEnsembleProvider) Start() error {
	return nil
}

func (FixedEnsembleProvider) Close() error {
	return nil
}

func (f FixedEnsembleProvider) GetConnectionString() string {
	return string(f)
}
