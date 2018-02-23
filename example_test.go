package curator_test

import (
	"log"
	"time"

	"github.com/eahydra/go-curator"
)

func ExampleNewZookeeperClient() {
	servers := "192.168.27.216:2181,192.168.27.217:2181,192.168.27.218:2181"
	ensemble := curator.NewFixedEnsembleProvider(servers)
	sessionTimeout := 3 * time.Second
	connectionTimeout := 1 * time.Second
	retryPolicy := curator.NewRetryForever(500 * time.Millisecond)
	client, err := curator.NewZookeeperClient(curator.DefaultZookeeperFactory, ensemble, sessionTimeout, connectionTimeout, retryPolicy, true)
	if err != nil {
		log.Println("failed to curator.NewZookeeperClient, err:", err)
		return
	}
	if err := client.Start(); err != nil {
		log.Println("failed to client.Start, err:", err)
		return
	}
	defer client.Close()

	children, stat, err := client.Children("/zookeeper")
	if err != nil {
		log.Println("failed to client.Children, err:", err)
		return
	}

	log.Printf("stat:%+v", stat)

	for _, child := range children {
		log.Println("child:", child)
	}
}

func ExampleZookeeperClientBuilder() {
	servers := "192.168.27.216:2181,192.168.27.217:2181,192.168.27.218:2181"
	client, err := curator.NewZookeeperClientBuidler().
		WithEnsembleProvider(curator.NewFixedEnsembleProvider(servers)).
		WithRetryPolicy(curator.NewRetryForever(100 * time.Millisecond)).
		Build()
	if err != nil {
		log.Println("failed to curator.NewZookeeperClient, err:", err)
		return
	}
	if err := client.Start(); err != nil {
		log.Println("failed to client.Start, err:", err)
		return
	}
	defer client.Close()

	children, stat, err := client.Children("/zookeeper")
	if err != nil {
		log.Println("failed to client.Children, err:", err)
		return
	}

	log.Printf("stat:%+v", stat)

	for _, child := range children {
		log.Println("child:", child)
	}
}
