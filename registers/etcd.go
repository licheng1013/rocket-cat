package registers

import (
	"context"
	"errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Etcd struct {
	ClientInfo ClientInfo        // 客户端信息
	ServerInfo ServerInfo        // 服务端信息
	conn       *ServiceDiscovery // 连接
	key        string            // key
	index      int               // 服务索引
	lock       sync.Mutex        // 锁
	ID         clientv3.LeaseID  // 租约ID
}

func (e *Etcd) GetIp() (ClientInfo, error) {
	list := e.conn.GetServices()
	if len(list) == 0 {
		return ClientInfo{}, errors.New("找不到服务")
	}
	//fmt.Println("list-", list)
	// 获取所有value
	var values []string
	for _, v := range list {
		values = append(values, v)
	}
	if e.index >= len(values) {
		e.index = 0
	}
	defer e.lock.Unlock()
	e.lock.Lock()
	ip := values[e.index]
	e.index++
	return parseIp(ip), nil
}

// 解析ip localhost:80 -> localhost 80
func parseIp(ip string) ClientInfo {
	split := strings.Split(ip, ":")
	parseInt, _ := strconv.ParseInt(split[1], 10, 64)
	return ClientInfo{Ip: split[0], Port: uint16(parseInt)}
}

func (e *Etcd) ListIp(serverName string) (clients []ClientInfo, errorInfo error) {
	list, err := e.conn.cli.Get(context.Background(), serverName, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(list.Kvs) == 0 {
		return nil, errors.New("找不到服务")
	}
	for _, kv := range list.Kvs {
		clients = append(clients, parseIp(string(kv.Value)))
	}
	return clients, nil
}

// Run 运行
func (e *Etcd) Run() {
	// 验证参数
	if e.ServerInfo.Ip == "" || e.ServerInfo.Port == 0 {
		panic("etcd地址错误")
	}
	if e.ClientInfo.Ip == "" || e.ClientInfo.Port == 0 || e.ClientInfo.ServiceName == "" || e.ClientInfo.RemoteName == "" {
		panic("etcd客户端地址错误")
	}
	e.key = e.ClientInfo.RemoteName + "-" + e.ClientInfo.Addr()
	// 注册服务
	e.conn = NewServiceDiscovery([]string{e.ServerInfo.Addr()})
	// 监听服务
	err := e.conn.WatchService(e.ClientInfo.RemoteName)
	if err != nil {
		panic("etcd监听失败" + err.Error())
	}
	_, err = e.conn.cli.Put(context.Background(), e.key, e.ClientInfo.Addr())
	if err != nil {
		panic("etcd注册失败" + err.Error())
	}
	// 心跳
	e.Heartbeat()
}

// Heartbeat 心跳
func (e *Etcd) Heartbeat() {
	resp, err := e.conn.cli.Grant(context.TODO(), 2)
	if err != nil {
		panic("etcd心跳失败" + err.Error())
	}
	_, err = e.conn.cli.Put(context.Background(), e.key, e.ClientInfo.Addr(), clientv3.WithLease(resp.ID))
	if err != nil {
		panic("etcd数据失败" + err.Error())
	}
	_, err = e.conn.cli.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		panic("etcd心跳失败" + err.Error())
	}
	e.ID = resp.ID
}

// Close 关闭
func (e *Etcd) Close() {
	_, _ = e.conn.cli.Revoke(context.TODO(), e.ID)
	_, _ = e.conn.cli.Delete(context.Background(), e.key)
	_ = e.conn.cli.Close()
}

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli        *clientv3.Client  //etcd client
	serverList map[string]string //服务列表
	lock       sync.Mutex
}

// NewServiceDiscovery  新建发现服务
func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic("etcd连接失败" + err.Error())
	}

	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

// WatchService 初始化服务列表和监视
func (s *ServiceDiscovery) WatchService(prefix string) error {
	//根据前缀获取现有的key
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	//监视前缀，修改变更的server
	go s.watcher(prefix)
	return nil
}

// watcher 监听前缀
func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	//log.Printf("watching prefix:%s now...", prefix)
	for resp := range rch {
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// SetServiceList 新增服务地址
func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = val
	//log.Println("put key :", key, "val:", val)
}

// DelServiceList 删除服务地址
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	//log.Println("del key:", key)
}

// GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	adds := make([]string, 0)
	for _, v := range s.serverList {
		adds = append(adds, v)
	}
	return adds
}

// Close 关闭服务
func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}
