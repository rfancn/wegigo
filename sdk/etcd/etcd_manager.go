package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"context"
	"log"
	"encoding/json"
	"fmt"
)

type EtcdManager struct {
	cli *clientv3.Client
}

func NewEtcdManager(address string, port int) *EtcdManager {
	etcdServerUrl := fmt.Sprintf("http://%s:%d", address, port)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdServerUrl},
		DialTimeout: ETCD_CONNECT_TIMEOUT,
	})

	if err != nil {
		log.Println("NewEtcdManager(): Error connect etcd server:", err)
		return nil
	}

	return &EtcdManager{cli: cli}
}

func (m *EtcdManager) Close() {
	m.cli.Close()
}

/**
  * Get operation
 */
//GetResp key/value response from ETCD
func (m *EtcdManager) GetResp(key string, opts ...clientv3.OpOption) *clientv3.GetResponse {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_GET_TIMEOUT)
	resp, err := m.cli.Get(ctx, key, opts...)
	cancel()
	if err != nil {
		log.Println("EtcdManager Put(): Error read from etcd:", err)
		return nil
	}

	return resp
}

//GetResp key/value response with prefix from ETCD
func (m *EtcdManager) GetRespWithPrefix(key string, opts ...clientv3.OpOption) *clientv3.GetResponse {
	return m.GetResp(key, clientv3.WithPrefix())
}

//GetBytes: get value bytes from ETCD
func (m *EtcdManager) GetBytes(key string) []byte {
	resp := m.GetResp(key)
	if len(resp.Kvs) < 1 {
		return nil
	}

	//get the first kv item
	return resp.Kvs[0].Value
}

//GetBytesList: get value bytes list from ETCD
func (m *EtcdManager) GetBytesList(key string) [][]byte {
	resp := m.GetRespWithPrefix(key)

	list := make([][]byte, 0)
	for _, ev := range resp.Kvs {
		list = append(list, ev.Value)
	}

	return list
}

/**
  * Put operation
 */
func (m *EtcdManager) PutBytes(key string, value []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_PUT_TIMEOUT)
	_, err := m.cli.Put(ctx, key, string(value))
	cancel()
	if err != nil {
		log.Println("EtcdManager Put(): Error put to etcd:", err)
		return false
	}

	return true
}

//Put key/value to ETCD
func (m *EtcdManager) PutValue(key string, value interface{}) bool {
	bv, err := json.Marshal(value)
	if err != nil {
		log.Println("EtcdManager Put(): Error marshal value")
		return false
	}

	return m.PutBytes(key, bv)
}

/**
  * Watch operation
 */
//Watch for ETCD key changes
func (m *EtcdManager) Watch(key string, opts ...clientv3.OpOption) <-chan clientv3.WatchResponse {
	return m.cli.Watch(context.Background(), key, opts...)
}

func (m *EtcdManager) WatchWithPrefix(key string, opts ...clientv3.OpOption) <-chan clientv3.WatchResponse {
	return m.cli.Watch(context.Background(), key, clientv3.WithPrefix())
}

/**
  * Delete operation
 */
 //Delete operation always WithPrefix, which
func (m *EtcdManager) Delete(key string, opts ...clientv3.OpOption) bool {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_PUT_TIMEOUT)
	_, err := m.cli.Delete(ctx, key, opts...)
	cancel()
	if err != nil {
		log.Println("EtcdManager Delete(): Error delete key in etcd:", err)
		return false
	}

	return true
}

func (m *EtcdManager) DeleteWithPrefix(key string) bool {
	return m.Delete(key, clientv3.WithPrefix())
}

