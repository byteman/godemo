// netclient
package czserver

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type NetClient struct {
	Con    net.Conn
	parser ProtoParser
	Device DevInfo
}
type DevInfo struct {
	DevID      uint16
	IpAddr     string
	Plate      string
	OnDateTime string
}
type DevInfoList []DevInfo

var clientList map[string]*NetClient = make(map[string]*NetClient, 100)
var mutex sync.Mutex

func CreateClient(con net.Conn) (client *NetClient) {
	fmt.Println("new Client")

	cli := &NetClient{}
	cli.Con = con
	cli.parser = ProtoParser{}
	cli.parser.Data = make([]byte, 0, 512)
	cli.parser.waitHead = true
	cli.Device.IpAddr = con.RemoteAddr().String()
	cli.Device.OnDateTime = time.Now().String()

	mutex.Lock()
	defer mutex.Unlock()
	clientList[con.RemoteAddr().String()] = cli
	return cli
}
func RemoveClient(con net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(clientList, con.RemoteAddr().String())
}
func (cli *NetClient) Handle(data []byte, n int) (err bool) {

	cli.parser.Prase(data, n)

	return true
}
func GetClient() DevInfoList {
	mutex.Lock()
	defer mutex.Unlock()
	//infos := make([]DevInfo, 0, 30)
	infos := DevInfoList{}
	for k, v := range clientList {
		fmt.Println(k, v)

		infos = append(infos, v.Device)

	}
	fmt.Println(infos)
	return infos

}
