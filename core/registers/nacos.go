package registers

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"time"
)

// Nacos 请使用构造方法获取实例  NewNacos
type Nacos struct {
	registerParam vo.RegisterInstanceParam
	logoutParam   vo.DeregisterInstanceParam
	serverConfigs []constant.ServerConfig
	namingClient  naming_client.INamingClient
}

func NewNacos() *Nacos {
	return &Nacos{}
}

func (n *Nacos) SetServerConfig(ip string, port uint64) {
	// 创建serverConfig的另一种方式 -> 此处链接nacos的配置
	n.serverConfigs = []constant.ServerConfig{
		*constant.NewServerConfig(ip, port, constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}
}

func (n *Nacos) Register(info RegisterInfo, name string) {
	if len(n.serverConfigs) == 0 {
		panic("未设置Nacos配置: SetServerConfig")
	}
	n.registerParam = vo.RegisterInstanceParam{
		Ip:          info.Ip,
		Port:        uint64(info.Port),
		ServiceName: name,
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
}

func (n *Nacos) Init() {
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(vo.NacosClientParam{ServerConfigs: n.serverConfigs})
	if err != nil {
		print(err)
	}
	n.namingClient = namingClient
	success, err := n.namingClient.RegisterInstance(n.registerParam)
	if err != nil || !success {
		print("注册失败:", err)
	}
}

func (n *Nacos) Logout() {
	success, err := n.namingClient.DeregisterInstance(n.logoutParam)
	if err != nil {
		print(err)
	}
	if success {
		log.Println("注销成功！")
	}
}

func (n *Nacos) AllInstances() {
	// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
	instances, err := n.namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: n.registerParam.ServiceName,
		GroupName:   n.registerParam.GroupName,
	})
	if err != nil {
		print(err)
	}
	log.Println("所有实例: ", instances)
}

// SelectInstances 获取指定条件的实例
func (n *Nacos) SelectInstances(serverName string) []model.Instance {
	// SelectInstances 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, err := n.namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serverName,
		GroupName:   n.registerParam.GroupName,
		HealthyOnly: true, //true 健康的实例
	})
	if err != nil {
		print(err)
	}
	return instances
}

// SelectOneHealthyInstance 获取指定条件的实例
func (n *Nacos) SelectOneHealthyInstance(serverName string) *model.Instance {
	// SelectInstances 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
	instances, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serverName,
		GroupName:   n.registerParam.GroupName,
	})
	if err != nil {
		print(err)
	}
	return instances
}

// Heartbeat 心跳功能
func (n *Nacos) Heartbeat() {
	go func() {
		for true {
			instance, err := n.namingClient.UpdateInstance(
				vo.UpdateInstanceParam{
					Ip:          n.registerParam.Ip,
					Port:        n.registerParam.Port,
					Weight:      n.registerParam.Weight,
					Enable:      n.registerParam.Enable,
					Healthy:     n.registerParam.Healthy,
					ServiceName: n.registerParam.ServiceName,
				})
			if err != nil || !instance {
				log.Panicln("更新实例失败,请检查Nacos!", err)
			}
			//log.Println("实例更新成功！")
			time.Sleep(1 * time.Second)
		}
	}()
}
