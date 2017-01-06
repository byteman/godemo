// netclient
package czserver

import (
	"fmt"
	models "hellobee/models"
	"net"
	"sync"
	"time"

	"github.com/astaxie/beego"
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
	DeviceId  uint16
	Version   string
	GpsReport uint8
	DevReport uint8
	Plate     string
	//OnDateTime string
	timeStamp time.Time
	UnixTime  int64
	IpAddr    string
}
type DevInfoList []DevInfo

var timeoutS int = 120
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
	//cli.Device.OnDateTime = time.Now().String()
	cli.Device.timeStamp = time.Now()
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
func handleOnline(ipaddr string, dev *DevicePara) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clientList[ipaddr]; ok { //存在}
		fmt.Println("handle online")
		clientList[ipaddr].Device.DeviceId = dev.DeviceId
		clientList[ipaddr].Device.DevReport = dev.DevReport
		clientList[ipaddr].Device.GpsReport = dev.GpsReport

		enc := mahonia.NewDecoder("GBK")
		src := string(dev.Plate[:])

		clientList[ipaddr].Device.Plate = enc.ConvertString(src)
		clientList[ipaddr].Device.Version = fmt.Sprintf("v%d.%d.%d", (dev.Version>>16)&0xff, (dev.Version>>8)&0xff, dev.Version&0xff)
	}
}
func resetTimeout(con net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clientList[con.RemoteAddr().String()]; ok { //存在}
		fmt.Println("reset timeout")
		clientList[con.RemoteAddr().String()].Device.timeStamp = time.Now()
	}
}
func handleTimeout() {
	mutex.Lock()
	defer mutex.Unlock()
	//dead := make([]string, 10)
	for _, value := range clientList {
		//fmt.Printf("%s->%-10s", key, value)
		diff := time.Now().Sub(value.Device.timeStamp)
		//fmt.Println(diff)
		s := time.Duration(timeoutS) * time.Second
		if diff > s {
			fmt.Println(value.Device.IpAddr, " timeout")
			value.Con.Close()
			//dead = append(dead, value.Device.IpAddr)
		}
	}
	//	for _, v := range dead {
	//		delete(clientList, v)
	//	}
}
func handleMsg(msg Message, con net.Conn) {
	fmt.Println("cmd=", msg.Head.Cmd)
	resetTimeout(con)
	switch msg.Head.Cmd {
	case CMD_DEV2HOST_ONE_WEIGHT:
		//var p PointWet

		p, ok := msg.Val.(*PointWet)
		if !ok {
			fmt.Println("convt PointWet failed", p)
			return
		}
		insertOneWeight(p)

	case CMD_DEV2HOST_ALL_WEIGHT:
		fallthrough
	case CMD_DEV2HOST_WATER_WEIGHT:

		p, ok := msg.Val.(*CommWeight)
		if !ok {
			fmt.Println("convt CommWeight failed", p)
			return
		}
		insertCommonWeight(p, int32(msg.Head.Cmd))

	case CMD_DEV_ONLINE:
		d, ok := msg.Val.(*DevicePara)
		if !ok {
			fmt.Println("convt DevicePara failed", d)
			return
		}
		handleOnline(con.RemoteAddr().String(), d)

	case CMD_DEV2HOST_HEART:
	default:
		fmt.Println("unkown cmd")
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
	for _, v := range clientList {
		//fmt.Println(k, v)
		v.Device.UnixTime = v.Device.timeStamp.Unix() * 1000

		infos = append(infos, v.Device)

	}
	//fmt.Println(infos)
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
func onlineTimeout(input chan bool) {
	t1 := time.NewTimer(time.Second * 5)
	//	t2 := time.NewTimer(time.Second * 10)
	var msg bool = false
	for {
		select {
		case msg = <-input:
			println(msg)
			if msg {
				fmt.Println("exit online timeout")
				break
			}

		case <-t1.C:
			//println("5s timer")
			handleTimeout()
			t1.Reset(time.Second * 5)

			//		case <-t2.C:
			//			println("10s timer")
			//			t2.Reset(time.Second * 10)
		}
	}
}

var quit chan bool
var Cfg = beego.AppConfig

func init() {

	timeoutS, _ = Cfg.Int("timeout")

	fmt.Println("timeout = ", timeoutS)
	go onlineTimeout(quit)
}
