package curator

import (
	"path"

	"github.com/samuel/go-zookeeper/zk"
)

func CreateAll(client *ZookeeperClient, nodePath string, value []byte, flags int32, aclv []zk.ACL) (string, error) {
	if exists, _, err := client.Exists(nodePath); exists && err == nil {
		return nodePath, zk.ErrNodeExists
	}

	i := len(nodePath)
	for i > 0 && nodePath[i-1] == '/' {
		i--
	}

	j := i
	for j > 0 && nodePath[j-1] != '/' {
		j--
	}

	if j > 1 {
		// Create parent
		if _, err := CreateAll(client, nodePath[0:j-1], []byte{}, 0, aclv); err != nil {
			if err != zk.ErrNodeExists {
				return "", err
			}
		}
	}

	return client.Create(nodePath, value, flags, aclv)
}

func DeleteAll(client *ZookeeperClient, node string) error {
	children, stat, err := client.Children(node)
	if err != nil {
		if err != zk.ErrNoNode {
			return err
		}
		return nil
	}

	for _, child := range children {
		if err = DeleteAll(client, path.Join(node, child)); err != nil {
			return err
		}
	}

	return client.Delete(node, stat.Version)
}
