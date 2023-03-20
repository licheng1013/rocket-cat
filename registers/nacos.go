package registers

import (
	"errors"
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"time"
)

// Nacos 请使用构造方法获取实例  NewNacos
type Nacos struct {
	namingClient       naming_client.INamingClient
	registerClientInfo RegisterInfo
	registerParam      vo.RegisterInstanceParam
	logoutParam        vo.DeregisterInstanceParam
	updateParam        vo.UpdateInstanceParam
}

func (n *Nacos) RegisterInfo() RegisterInfo {
	return n.registerClientInfo
}

func (n *Nacos) Close() {
	success, err := n.namingClient.DeregisterInstance(n.logoutParam)
	if err != nil {
		common.Logger().Println("注销错误:" + err.Error())
	}
	if success {
		common.Logger().Println("注销成功！")
	}
}

func (n *Nacos) Register(info RegisterInfo) {
	// 创建serverConfig的另一种方式 -> 此处链接nacos的配置
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(info.Ip, uint64(info.Port), constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}
	var err error
	// 创建服务发现客户端的另一种方式 (推荐)
	n.namingClient, err = clients.NewNamingClient(vo.NacosClientParam{ServerConfigs: serverConfigs})
	if err != nil {
		panic(err)
	}
	if success, err := n.namingClient.RegisterInstance(n.registerParam); err != nil || !success {
		panic(err)
	}
	common.Logger().Println("注册中心:", info.Ip+":"+fmt.Sprint(info.Port))
	go n.heartbeat() // 心跳功能
}

// GetIp 获取单个ip
func (n *Nacos) GetIp() (RegisterInfo, error) {
	// SelectList 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: n.registerClientInfo.RemoteName,
		GroupName:   n.registerParam.GroupName,
	})
	if err != nil {
		return RegisterInfo{}, errors.New("获取实例为空")
	}
	return RegisterInfo{Ip: instances.Ip, Port: uint16(instances.Port), ServiceName: instances.ServiceName}, nil
}

// ListIp 获取ip
func (n *Nacos) ListIp(serverName string) ([]RegisterInfo, error) {
	instances := n.SelectList(serverName)
	infos := make([]RegisterInfo, 0)
	if len(instances) == 0 {
		return infos, errors.New("获取实例为空")
	}
	for _, item := range instances {
		infos = append(infos, RegisterInfo{Ip: item.Ip, Port: uint16(item.Port), ServiceName: item.ServiceName})
	}
	return infos, nil
}

func NewNacos() *Nacos {
	return &Nacos{}
}

// TODO 这里作为保留教程
//func (n *Nacos) allInstances() {
//	// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0 包括关闭得实例
//	instances, err := n.namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
//		ServiceName: n.registerParam.ServiceName,
//		GroupName:   n.registerParam.GroupName,
//	})
//	if err != nil {
//		print(err)
//	}
//	common.Logger().Println("所有实例: ", instances)
//}

// SelectList 只获取存活的实例
func (n *Nacos) SelectList(serverName string) []model.Instance {
	// SelectList 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, _ := n.namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serverName,
		GroupName:   n.registerParam.GroupName,
		HealthyOnly: true, //true 健康的实例
	})
	return instances
}

// Heartbeat 心跳功能
func (n *Nacos) heartbeat() {
	for true {
		instance, err := n.namingClient.UpdateInstance(n.updateParam)
		if err != nil || !instance {
			common.Logger().Println("更新实例失败,请检查Nacos!", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func (n *Nacos) RegisterClient(info RegisterInfo) {
	n.registerClientInfo = info
	// 这里是设置注册客户端的参数
	n.registerParam = vo.RegisterInstanceParam{
		Ip:          info.Ip,
		Port:        uint64(info.Port),
		ServiceName: info.ServiceName,
		GroupName:   "DEFAULT_GROUP",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
	}
	n.logoutParam = vo.DeregisterInstanceParam{
		Ip:          n.registerParam.Ip,
		Port:        n.registerParam.Port,
		ServiceName: n.registerParam.ServiceName,
	}
	n.updateParam = vo.UpdateInstanceParam{
		Ip:          n.registerParam.Ip,
		Port:        n.registerParam.Port,
		Weight:      n.registerParam.Weight,
		Enable:      n.registerParam.Enable,
		Healthy:     n.registerParam.Healthy,
		ServiceName: n.registerParam.ServiceName,
	}
}
