package curator

import (
	"path"
	"sort"
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

type Locker interface {
	Acquire() error
	Release() error
}

type Mutex struct {
	client   *ZookeeperClient
	basePath string
	lockName string
	lockPath string
	aclv     []zk.ACL
}

var _ Locker = &Mutex{}

func NewMutex(client *ZookeeperClient, basePath string, aclv []zk.ACL) *Mutex {
	return &Mutex{
		client:   client,
		basePath: basePath,
		lockName: "lock-",
		aclv:     aclv,
	}
}

func createTheLock(client *ZookeeperClient, lockPath string, data []byte, aclv []zk.ACL) (string, error) {
	var nodePath string
	var err error

	for i := 0; i < 2; i++ {
		nodePath, err = client.CreateProtectedEphemeralSequential(lockPath, data, aclv)
		if err == nil {
			break
		} else if err == zk.ErrNoNode {
			_, err = CreateAll(client, path.Dir(lockPath), []byte{}, 0, aclv)
			if err != nil {
				if err != zk.ErrNodeExists {
					break
				}
				err = nil
			}
		}
	}
	return nodePath, err
}

func (m *Mutex) Acquire() error {
	for {
		nodePath, err := createTheLock(m.client, path.Join(m.basePath, m.lockName), []byte{}, m.aclv)
		if err != nil {
			return err
		}

		err = m.internalLockLoop(nodePath)
		if err == nil {
			m.lockPath = nodePath
			break

		} else if err != zk.ErrNoNode {
			return err
		}
	}

	return nil
}

type mutexSortChildren struct {
	lockName string
	children []string
}

func (m *mutexSortChildren) Len() int {
	return len(m.children)
}

func (m *mutexSortChildren) Less(i, j int) bool {
	lhs := strings.Split(m.children[i], m.lockName)
	rhs := strings.Split(m.children[j], m.lockName)
	return lhs[len(lhs)-1] < rhs[len(rhs)-1]
}

func (m *mutexSortChildren) Swap(i, j int) {
	m.children[i], m.children[j] = m.children[j], m.children[i]
}

func (m *Mutex) internalLockLoop(nodePath string) error {
	deleteNode := true
	defer func() {
		if deleteNode {
			m.client.Delete(path.Join(m.basePath, m.lockPath), -1)
		}
	}()

	seq := path.Base(nodePath)

	sorter := &mutexSortChildren{}

	for {
		children, _, err := m.client.Children(m.basePath)
		if err != nil {
			return err
		}

		sorter.lockName = m.lockName
		sorter.children = children
		sort.Sort(sorter)

		index := -1
		for i, v := range children {
			if v == seq {
				index = i
				break
			}
		}

		// not find our seq node, maybe the node had been removed when session expired
		if index < 0 {
			return zk.ErrNoNode
		}

		if index < 1 {
			// got the lock
			deleteNode = false
			break
		} else {
			previousSeq := children[index-1]
			exist, _, ch, err := m.client.ExistsW(path.Join(m.basePath, previousSeq))
			if err != nil {
				return err
			}

			if !exist {
				continue
			}

			ev := <-ch
			if ev.Err != nil {
				return ev.Err
			}
		}
	}

	return nil
}

func (m *Mutex) Release() error {
	if m.lockPath == "" {
		return zk.ErrNotLocked
	}

	err := m.client.Delete(m.lockPath, -1)
	m.lockPath = ""
	return err
}
