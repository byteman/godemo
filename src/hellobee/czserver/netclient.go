// netclient
package czserver

import (
	"fmt"
	models "hellobee/models"
	"net"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/mahonia"
	_ "github.com/mattn/go-sqlite3"
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
func handleOnline(devid uint16) {

}
func handleMsg(msg Message, con net.Conn) {
	switch msg.Head.Cmd {
	case CMD_DEV2HOST_ONE_WEIGHT:
		var p PointWet
		switch v := msg.Val.(type) {
		case PointWet:
			p = v

			insertOneWeight(&p)
		}

	case CMD_DEV2HOST_ALL_WEIGHT:
		fallthrough
	case CMD_DEV2HOST_WATER_WEIGHT:

		var c CommWeight
		switch v := msg.Val.(type) {
		case CommWeight:
			c = v

			insertCommonWeight(&c, int32(msg.Head.Cmd))
		}

	case CMD_DEV_ONLINE:
		var c uint16
		switch v := msg.Val.(type) {
		case uint16:
			c = v

			handleOnline(c)
		}
	case CMD_DEV2HOST_HEART:
	}
}
func (cli *NetClient) Handle(data []byte, n int) (err bool) {

	msgList := cli.parser.Prase(data, n)
	fmt.Println("Handle msg", msgList)
	for i, v := range msgList {
		fmt.Println("handle msg", i)
		handleMsg(v, cli.Con)
	}
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

func fmtDate(dt DateDef) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", 2000+int(dt.Year), dt.Month, dt.Day, dt.Hour, dt.Min, dt.Sec)
}
func fmtGps(gps GpsDef) string {
	return fmt.Sprintf("%.6f,%.6f,%c,%c", gps.Latitude, gps.Longitude, gps.Ew, gps.Ns)
}
func insertOneWeight(pwt *PointWet) {
	msg := new(models.OneWeight)
	enc := mahonia.NewDecoder("GBK")
	src := string(pwt.Plate[:])
	msg.WType = 1
	msg.Weight = pwt.Wet
	msg.LicensePlate = enc.ConvertString(src)

	src = string(pwt.Duty[:])
	msg.Duty = enc.ConvertString(src)

	msg.UpDate = fmtDate(pwt.UpDate)
	msg.WetDate = fmtDate(pwt.Wdate)
	msg.Gps = fmtGps(pwt.Gps)
	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库
	_, err := o.Insert(msg)
	if err != nil {
		fmt.Println(err)
	}
}
func insertCommonWeight(pwt *CommWeight, WType int32) {
	msg := new(models.OneWeight)
	enc := mahonia.NewDecoder("GBK")
	src := string(pwt.Plate[:])
	msg.Weight = pwt.Wet
	msg.WType = WType
	msg.LicensePlate = enc.ConvertString(src)
	msg.Gps = fmtGps(pwt.Gps)
	msg.UpDate = fmtDate(pwt.UpDate)

	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库
	_, err := o.Insert(msg)
	if err != nil {
		fmt.Println(err)
	}
}
