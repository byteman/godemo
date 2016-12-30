// protocal
package czserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"utils"

	models "hellobee/models"

	"github.com/astaxie/beego/orm"
	"github.com/mahonia"
	_ "github.com/mattn/go-sqlite3"
)

type MsgHead struct {
	DevId uint16
	Dir   uint8
	Cmd   uint8
	Oper  uint8
	Len   uint16
}
type GpsDef struct {
	Longitude float64 // 经度
	Latitude  float64 // 纬度
	Ns        uint8   // 南北值为,'n'或's'
	Ew        uint8   // 东西,'e'或'w'
}

type DateDef struct {
	Year  uint8 // 当前年减去2000,如2016年，year实际保存16。
	Month uint8
	Day   uint8
	Hour  uint8
	Min   uint8
	Sec   uint8
}

const (
	LICENSE_LEN = 10
	DUTY_LEN    = 16
)

// B组内容-总重量
type CommWeight struct {
	Wet    int32             // 单点重量、
	Plate  [LICENSE_LEN]byte // 车辆号牌信息（或本机信息）、
	Gps    GpsDef            // GPS信息、
	UpDate DateDef           // 发送的实时日期时间.

}

// A组内容-单点重量
type PointWet struct {
	Wet    int32             // 单点重量、
	Wdate  DateDef           // 单点重量的获取日期时间、
	Gps    GpsDef            // GPS信息、
	Plate  [LICENSE_LEN]byte // 车辆号牌信息（或本机信息）、
	Duty   [DUTY_LEN]byte    // 值班员（或司机信息）、
	UpDate DateDef           // 发送的实时日期时间.
}

func (h *MsgHead) Init(d []byte) {

	err := binary.Read(bytes.NewReader(d), binary.LittleEndian, h)
	if err != nil {

	}
	fmt.Printf("%v\n", *h) // {1 2 3 4 5 6}

}

type ProtoParser struct {
	Data     []byte
	Header   MsgHead
	waitHead bool
}

const (
	CMD_ONE_WEIGHT   = 1
	CMD_TOTAL_WEIGHT = 2
	CMD_WATER_WEIGHT = 3
)

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
func unSerial(data []byte, n uint16, msg interface{}) bool {

	r := bytes.NewReader(data[:n])
	err := binary.Read(r, binary.LittleEndian, msg)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Print("%v", msg)
	return true
}
func handleMsg(head MsgHead, d []byte, n uint16) {
	fmt.Println("handle msgtype ", head.Cmd)
	switch head.Cmd {
	case CMD_ONE_WEIGHT:

		pwt := &PointWet{}

		if unSerial(d, n, pwt) {
			insertOneWeight(pwt)
		}

		break
	case CMD_TOTAL_WEIGHT:
		pwt := &CommWeight{}
		if unSerial(d, n, pwt) {
			insertCommonWeight(pwt, int32(head.Cmd))
		}
		break
	case CMD_WATER_WEIGHT:
		pwt := &CommWeight{}
		if unSerial(d, n, pwt) {
			insertCommonWeight(pwt, int32(head.Cmd))
		}
		break

	}
}

//分析数据协议.
func (p *ProtoParser) Prase(data []byte, n int) {
	fmt.Printf("%d %d\n", len(p.Data), n)
	for i := 0; i < n; i++ {
		p.Data = append(p.Data, data[i])
	}
	fmt.Println("come here")
	for {
		fmt.Println("len", len(p.Data))
		if len(p.Data) <= 0 {
			break
		}
		if p.waitHead == true {
			fmt.Println("find head")
			if len(p.Data) >= 7 {
				fmt.Println("find head ok")
				p.Header.Init(p.Data[:7])
				//p.Data = p.Data[7:]
				p.waitHead = false

			} else {
				fmt.Println("< 7")
				break
			}
		} else {
			fmt.Println("find data len", p.Header.Len)
			var size int = int(p.Header.Len + 9)
			if len(p.Data) < size {
				break
			}
			var crc16 uint16 = utils.Reentrent_CRC16(p.Data, uint32(p.Header.Len+7))
			var crc16_data uint16 = utils.BytesToUint16(p.Data[p.Header.Len+7:])
			fmt.Printf("crc1=%d,crc2=%d\n", crc16, crc16_data)
			p.waitHead = true
			if crc16 != crc16_data {
				p.Data = p.Data[p.Header.Len+9:]
				fmt.Println("crc error")
			}
			p.Data = p.Data[7 : 7+p.Header.Len]
			handleMsg(p.Header, p.Data, p.Header.Len)
			p.Data = p.Data[p.Header.Len:]
		}
	}
}

func init() {
	// 需要在init中注册定义的model
	fmt.Println("init sqlite3")
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "database/orm_test.db")
	orm.RegisterModel(new(models.OneWeight))
	orm.RunSyncdb("default", false, true)

}
