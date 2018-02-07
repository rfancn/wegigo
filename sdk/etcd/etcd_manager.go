package etcd

import (
	"context"
	"log"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
)

type EtcdManager struct {
	cli *clientv3.Client
}

func NewEtcdManager(url string) (*EtcdManager, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: ETCD_CONNECT_TIMEOUT,
	})

	if err != nil {
		return nil, err
	}

	return &EtcdManager{cli: cli}, nil
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

//GetValue: get value bytes from ETCD
func (m *EtcdManager) GetValue(key string) []byte {
	resp := m.GetResp(key)
	if len(resp.Kvs) < 1 {
		return nil
	}

	//get the first kv item
	return resp.Kvs[0].Value
}

//GetValue: get value bytes slice from ETCD
func (m *EtcdManager) GetValueList(key string) [][]byte {
	resp := m.GetRespWithPrefix(key)

	vList := make([][]byte, 0)
	for _, ev := range resp.Kvs {
		vList = append(vList, ev.Value)
	}
	return vList
}

//GetItems: get item list from ETCD
func (m *EtcdManager) GetItemList(key string) []map[string]string {
	resp := m.GetRespWithPrefix(key)

	itemList := make([]map[string]string, 0)
	for _, ev := range resp.Kvs {
		item := make(map[string]string)
		item[string(ev.Key)] = string(ev.Value)
		itemList = append(itemList, item)
	}

	return itemList
}

/**
  * Put operation
 */
func (m *EtcdManager) Put(key string, value string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_PUT_TIMEOUT)
	_, err := m.cli.Put(ctx, key, string(value))
	cancel()
	if err != nil {
		log.Println("EtcdManager Put(): Error put to etcd:", err)
		return false
	}

	return true
}

func (m *EtcdManager) PutValueBytes(key string, value []byte) bool {
	return m.Put(key, string(value))
}

//Put key/value to ETCD
func (m *EtcdManager) PutValueAny(key string, value interface{}) bool {
	bv, err := json.Marshal(value)
	if err != nil {
		log.Println("EtcdManager Put(): Error marshal value")
		return false
	}

	return m.PutValueBytes(key, bv)
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

/**
  transaction put
 */
func (m *EtcdManager) TxnPut(key string, value string) bool {
	//Read stored key value
	v := m.GetResp(key).Kvs[0]

	kvc := clientv3.NewKV(m.cli)
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_PUT_TIMEOUT)
	//new transaction
	resp, err := kvc.Txn(ctx).
		//if modification revision equals to what we get
		If(clientv3.Compare(clientv3.ModRevision(key), "=", v.ModRevision)).
		//then put value
		Then(clientv3.OpPut(key, value)).
		Commit()
	cancel()
	if err != nil {
		log.Println("EtcdManager TxnPut(): Error commit transaction:", err)
		return false
	}

	return resp.Succeeded
}

func (m *EtcdManager) TxnPutValueBytes(key string, value []byte) bool {
	return m.TxnPut(key, string(value))
}

func (m *EtcdManager) TxnPutValueAny(key string, value interface{}) bool {
	bv, err := json.Marshal(value)
	if err != nil {
		log.Println("EtcdManager TxnPutValueAny(): Error marshal value")
		return false
	}

	return m.TxnPutValueBytes(key, bv)
}

