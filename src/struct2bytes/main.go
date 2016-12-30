package main

/*
#include <stdio.h>
#include <string.h>
typedef unsigned char  uint8;
typedef unsigned short uint16;


#pragma pack(push) // 将当前pack设置压栈保存
#pragma pack(1)// 必须在结构体定义之前使用
typedef struct{
	uint16 dev_id;
	uint8 dir    ;
	uint8 cmd ;
	uint8 oper ;
	uint16 len  ;
}MsgHead;
void setMsgHead(void* ptr,void* data)
{
	memcpy(ptr,data,sizeof(MsgHead));
	//MsgHead* header = (struct MsgHead*)ptr;

}
#pragma pack(pop)
#include <stdint.h>
#pragma pack(push, 1)
typedef struct {
	uint16_t size;
	uint8 dir    ;
	uint8 cmd ;
	uint8 oper ;
	uint16_t data3;
} mydata2;
typedef struct {
	uint16_t size;
	uint16_t msgtype;
	uint32_t sequnce;
	uint8_t data1;
	uint32_t data2;
	uint16_t data3;
} mydata;
#pragma pack(pop)

mydata2 foo = {
	1, 2, 3, 4, 5,
};

int size() {
	return sizeof(mydata2);
}

*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

// A组内容-单点重量
type PointWet struct {
	Wet    uint32            // 单点重量、
	Wdate  DateDef           // 单点重量的获取日期时间、
	Gps    GpsDef            // GPS信息、
	Plate  [LICENSE_LEN]byte // 车辆号牌信息（或本机信息）、
	Duty   [DUTY_LEN]byte    // 值班员（或司机信息）、
	UpDate DateDef           // 发送的实时日期时间.
}

func main() {

	a := [60]byte{0}
	bs := a[:]
	//fmt.Printf("len %d data %v\n", len(d), d)
	fmt.Printf("len %d data %v\n", len(bs), bs)
	var data PointWet
	fmt.Printf("%v\n", data)
	reader := bytes.NewReader(bs)
	err := binary.Read(reader, binary.LittleEndian, &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", data) // {1 2 3 4 5 6}
}
