package curator

import "time"

type ZookeeperClientBuilder struct {
	factory           ZookeeperFactory
	ensemble          EnsembleProvider
	sessionTimeout    time.Duration
	connectionTimeout time.Duration
	retryPolicy       RetryPolicy
	canBeReadOnly     bool
}

func NewZookeeperClientBuidler() *ZookeeperClientBuilder {
	return &ZookeeperClientBuilder{
		factory:           DefaultZookeeperFactory,
		sessionTimeout:    DefaultSessionTimeout,
		connectionTimeout: DefaultConnectionTimeout,
		retryPolicy:       NewExponentialBackoffRetry(29, 500*time.Millisecond, 5*time.Second),
		canBeReadOnly:     true,
	}
}

func (b *ZookeeperClientBuilder) WithZookeeperFactory(factory ZookeeperFactory) *ZookeeperClientBuilder {
	b.factory = factory
	return b
}

func (b *ZookeeperClientBuilder) WithEnsembleProvider(ensemble EnsembleProvider) *ZookeeperClientBuilder {
	b.ensemble = ensemble
	return b
}

func (b *ZookeeperClientBuilder) WithSessionTimeout(timeout time.Duration) *ZookeeperClientBuilder {
	b.sessionTimeout = timeout
	return b
}

func (b *ZookeeperClientBuilder) WithConnectionTimeout(timeout time.Duration) *ZookeeperClientBuilder {
	b.connectionTimeout = timeout
	return b
}

func (b *ZookeeperClientBuilder) WithRetryPolicy(policy RetryPolicy) *ZookeeperClientBuilder {
	b.retryPolicy = policy
	return b
}

func (b *ZookeeperClientBuilder) WithCanBeReadOnly(canBeReadOnly bool) *ZookeeperClientBuilder {
	b.canBeReadOnly = canBeReadOnly
	return b
}

func (b *ZookeeperClientBuilder) Build() (*ZookeeperClient, error) {
	return NewZookeeperClient(b.factory, b.ensemble, b.sessionTimeout, b.connectionTimeout, b.retryPolicy, b.canBeReadOnly)
}
