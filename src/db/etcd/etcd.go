package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"xlog"

	"github.com/coreos/etcd/clientv3"
)

var GlobalEtcd *Etcd

type Etcd struct {
	RawClient *clientv3.Client
}

func NewEtcd(userName, password string, hosts []string) *Etcd {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: hosts,
		Username:  userName,
		Password:  password,
	})

	if err != nil {
		xlog.Logger.Info(fmt.Sprintf("hosts: %v  err:%v", hosts, err.Error()))
		return nil
	}

	etcd := new(Etcd)
	etcd.RawClient = cli

	return etcd
}

var (
	ErrKeyNotFound   = fmt.Errorf("key not found")
	ErrMarshalJson   = fmt.Errorf("unable to marshal")
	ErrUnmarshalJson = fmt.Errorf("unable to unmarshal")
	ErrEtcd          = fmt.Errorf("etcd error")
)

// PutMarshal
func (c *Etcd) PutMarshal(key string, value interface{}, ttlSecond int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		xlog.Logger.Error(err.Error())
		return fmt.Errorf("%v:%v", ErrMarshalJson, err)
	}

	if ttlSecond > 0 {
		result, err := c.RawClient.Grant(context.Background(), ttlSecond)
		if err != nil {
			xlog.Logger.Error(err)
			return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
		}
		xlog.Logger.Debug("LeaseID ", fmt.Sprintf("%x", result.ID), ttlSecond)

		_, err = c.RawClient.Put(context.Background(), key, string(data), clientv3.WithLease(result.ID))
		if err != nil {
			xlog.Logger.Error(err.Error())
			return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
		}
		return nil
	}

	_, err = c.RawClient.Put(context.Background(), key, string(data))
	if err != nil {
		xlog.Logger.Error(err.Error())
		return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
	}

	return nil
}

// GetUnmarshal
func (c *Etcd) GetUnmarshal(key string, value interface{}, opts ...clientv3.OpOption) error {
	result, err := c.RawClient.Get(context.Background(), key, opts...)
	if err != nil {
		xlog.Logger.Error(err.Error())
		return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
	}

	if result.Count == 0 {
		xlog.Logger.Warn("key ", key)
		return ErrKeyNotFound
	}

	if err := json.Unmarshal(result.Kvs[0].Value, value); err != nil {
		xlog.Logger.Error(err.Error())
		return fmt.Errorf("%v:%v", ErrUnmarshalJson, err.Error())
	}

	return nil
}

// ListUnmarshal
func (c *Etcd) ListUnmarshal(key string, valueArr interface{}, opts ...clientv3.OpOption) error {
	val := reflect.ValueOf(valueArr)
	ind := reflect.Indirect(val)
	slice := ind

	opts = append(opts, clientv3.WithPrefix())

	result, err := c.RawClient.Get(context.Background(), key, opts...)
	if err != nil {
		xlog.Logger.Error(err.Error())
		return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
	}

	if result.Count == 0 {
		xlog.Logger.Warn("key ", key)
		return ErrKeyNotFound
	}
	for _, kvs := range result.Kvs {
		v := reflect.New(ind.Type().Elem())
		if err := json.Unmarshal(kvs.Value, v.Interface()); err != nil {
			xlog.Logger.Error(err.Error())
			continue
		}
		slice = reflect.Append(slice, v.Elem())
	}
	ind.Set(slice)

	return nil
}

// RangeUnmarshal
func (c *Etcd) RangeUnmarshal(keyBeg, keyEnd string, valueArr interface{}, opts ...clientv3.OpOption) error {
	val := reflect.ValueOf(valueArr)
	ind := reflect.Indirect(val)
	slice := ind

	opts = append(opts, clientv3.WithRange(keyEnd))

	result, err := c.RawClient.Get(context.Background(), keyBeg, opts...)
	if err != nil {
		xlog.Logger.Error(err.Error())
		return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
	}

	if result.Count == 0 {
		xlog.Logger.Warn("key1 -- key2 :  ", keyBeg, keyEnd)
		return ErrKeyNotFound
	}
	for _, kvs := range result.Kvs {
		v := reflect.New(ind.Type().Elem())
		if err := json.Unmarshal(kvs.Value, v.Interface()); err != nil {
			xlog.Logger.Error(err.Error())
			continue
		}
		slice = reflect.Append(slice, v.Elem())
	}
	ind.Set(slice)

	return nil
}

// Delete
func (c *Etcd) Delete(key string) error {
	result, err := c.RawClient.Get(context.Background(), key)
	if err != nil {
		xlog.Logger.Error(err.Error())
		return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
	}

	if result.Count == 0 {
		xlog.Logger.Warn("key ", key)
		return ErrKeyNotFound
	}

	if _, err = c.RawClient.Delete(context.Background(), key); err != nil {
		xlog.Logger.Error(err.Error())
		return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
	}
	return nil
}

// IsExist
func (c *Etcd) IsExist(key string) bool {
	if key == "" {
		return false
	}
	resp, err := c.RawClient.Get(context.Background(), key, clientv3.WithKeysOnly())
	if err != nil {
		xlog.Logger.Error(err)
		return false
	}
	if resp.Count > 0 {
		return true
	}
	return false
}

// Close
func (c *Etcd) Close() error {
	if c.RawClient == nil {
		return nil
	}

	err := c.RawClient.Close()
	if err != nil {
		return fmt.Errorf("%v:%v", ErrEtcd, err.Error())
	}
	return nil
}

func (c *Etcd) MemberList() []string {
	if c.RawClient == nil {
		return nil
	}

	memberList, err := c.RawClient.MemberList(context.Background())
	if err != nil {
		xlog.Logger.Error("GetMemberList|error : ", err.Error())
		return nil
	}

	hosts := []string{}
	for _, member := range memberList.Members {
		for _, clientURL := range member.ClientURLs {
			httpHost := strings.TrimPrefix(clientURL, "http://")
			if httpHost != "" {
				hosts = append(hosts, httpHost)
				continue
			}

			httpsHost := strings.TrimPrefix(clientURL, "https://")
			if httpsHost != "" {
				hosts = append(hosts, httpsHost)
				continue
			}
		}
	}

	return hosts
}

func init() {
	userName := ""
	password := ""
	hosts := []string{"http://192.168.8.202:12379"}

	GlobalEtcd = NewEtcd(userName, password, hosts)
	if GlobalEtcd == nil {
		xlog.Logger.Panic(fmt.Sprintf("init etcd error %v %v %v ", userName, password, hosts))
	}

}
