package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type TCPHeader struct {
	SrcPort   uint16
	DstPort   uint16
	SeqNum    uint32
	AckNum    uint32
	Offset    uint8
	Flag      uint8
	Window    uint16
	Checksum  uint16
	UrgentPtr uint16
}

type PsdHeader struct {
	SrcAddr   uint32
	DstAddr   uint32
	Zero      uint8
	ProtoType uint8
	TcpLength uint16
}

func (tcp TCPHeader)makeTcpHeader(destPort int ,flag int) TCPHeader{
	rand.Seed(time.Now().UnixNano())
	tcp.SrcPort=uint16(rand.Intn(65535)+1)
	tcp.DstPort=uint16(destPort)
	tcp.SeqNum=uint32(rand.Intn(4000000000)+1)
	if flag==2{
		tcp.AckNum=0
	}else{
		tcp.AckNum=uint32(rand.Intn(4000000000)+1)
	}
	tcp.Flag=uint8(flag)
	tcp.Offset=uint8(uint16(unsafe.Sizeof(TCPHeader{}))/4) << 4
	tcp.Window=uint16(8192)
	tcp.UrgentPtr=uint16(0)
	return tcp
}

func (psd PsdHeader)makePsdHeader(srcIp string,destIp string) PsdHeader{
	psd.SrcAddr=inet_addr(srcIp)
	psd.DstAddr=inet_addr(destIp)
	psd.ProtoType= syscall.IPPROTO_TCP//6代表tcp
	psd.TcpLength=uint16(20)
	psd.Zero=0
	return psd
}

const (
	ipv4Version      = 4
	ipv4HeaderLen    = 20
	ipv4MaxHeaderLen = 60
)

// A ipv4 header
type IpHeader struct {
	Version  int    // 协议版本 4bit
	Len      int    // 头部长度 4bit
	TOS      int    // 服务类   8bit
	TotalLen int    // 包长		16bit
	ID       int    // id		8bit
	Flags    int    // flags	3bit
	FragOff  int    // 分段偏移量 13bit
	TTL      int    // 生命周期 4bit
	Protocol int    // 上层服务协议4bit
	Checksum int    // 头部校验和16bit
	Src      net.IP // 源IP  	32bit
	Dst      net.IP // 目的IP  	32bit
	Options  []byte // 选项, extension headers
}
//

// Marshal encode ipv4 header
func (h IpHeader) Marshal() ([]byte, error) {

	hdrlen := ipv4HeaderLen + len(h.Options)
	b := make([]byte, hdrlen)

	//版本和头部长度
	b[0] = byte(ipv4Version<<4 | (hdrlen >> 2 & 0x0f))
	b[1] = byte(h.TOS)

	binary.BigEndian.PutUint16(b[2:4], uint16(h.TotalLen))
	binary.BigEndian.PutUint16(b[4:6], uint16(h.ID))

	flagsAndFragOff := (h.FragOff & 0x1fff) | int(h.Flags<<13)
	binary.BigEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))

	b[8] = byte(h.TTL)
	b[9] = byte(h.Protocol)

	binary.BigEndian.PutUint16(b[10:12], uint16(h.Checksum))

	if ip := h.Src.To4(); ip != nil {
		copy(b[12:16], ip[:net.IPv4len])
	}

	if ip := h.Dst.To4(); ip != nil {
		copy(b[16:20], ip[:net.IPv4len])
	} else {
		return nil, errors.New("missing address")
	}

	if len(h.Options) > 0 {
		copy(b[ipv4HeaderLen:], h.Options)
	}

	return b, nil
}

func (ip IpHeader)makeIpHeader (srcIp net.IP,dstIp net.IP) IpHeader{
	ip.ID=1
	ip.TTL=255
	ip.Protocol=syscall.IPPROTO_TCP
	ip.Checksum=0 // 系统自动填充
	ip.Src=srcIp
	ip.Dst=dstIp
	return ip
}

type Handle uintptr

func fakeIp()string{
	var fakeIp string
	rand.Seed(time.Now().UnixNano())
	for i:= 0;i<4 ;i++  {
		randomIp:=strconv.Itoa(rand.Intn(255)+1)//生成[0,256)的随机数
		fakeIp = fakeIp+randomIp
		if i!=3{
			fakeIp=fakeIp+"."
		}
	}
	return fakeIp
}

func checkSum(data []byte) uint16  {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	//以每16位为单位进行求和，直到所有的字节全部求完或者只剩下一个8位字节（如果剩余一个8位字节说明字节数为奇数个）
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	//如果字节数为奇数个，要加上最后剩下的那个8位字节
	if length > 0 {
		sum += uint32(data[index])
	}
	//加上高16位进位的部分
	sum += (sum >> 16)
	//别忘了返回的时候先求反
	return uint16(^sum)
}

func inet_addr(ipaddr string) uint32 {
	var (
		segments []string = strings.Split(ipaddr, ".")
		ip       [4]uint64
		ret      uint64
	)
	for i := 0; i < 4; i++ {
		ip[i], _ = strconv.ParseUint(segments[i], 10, 64)
	}
	ret = ip[3]<<24 + ip[2]<<16 + ip[1]<<8 + ip[0]
	return uint32(ret)
}

func packHeader(tcpHeader TCPHeader,pstHeader PsdHeader)(TCPHeader,bytes.Buffer){
	var buffer bytes.Buffer
	err:=binary.Write(&buffer, binary.BigEndian, pstHeader)
	if err!=nil{
		fmt.Println("pstHeader err is ",err)
	}
	err=binary.Write(&buffer, binary.BigEndian, tcpHeader)
	if err!=nil{
		fmt.Println("tcpHeader err is ",err)
	}
	tcpHeader.Checksum = checkSum(buffer.Bytes())
	buffer.Reset()
	err=binary.Write(&buffer, binary.BigEndian, tcpHeader)
	if err!=nil{
		fmt.Println("reset tcpHeader err is ",err)
	}
	return tcpHeader,buffer
}

func work(srcIp string,destIp string,destPort int,flag int,fd syscall.Handle) {
	var tcpHeader TCPHeader
	var pstHeader PsdHeader
	var ipHeader IpHeader
	var buffer bytes.Buffer
	tcpHeader = tcpHeader.makeTcpHeader(destPort, flag)
	pstHeader = pstHeader.makePsdHeader(srcIp, destIp)
	tcpHeader, buffer = packHeader(tcpHeader, pstHeader) //插入校验和
	tcpByte := buffer.Bytes()                            //转为数组
	destIPbyte := net.IP(make([]byte, 4))
	binary.BigEndian.PutUint32(destIPbyte[0:4], inet_addr(destIp))
	srcIPbyte := net.IP(make([]byte, 4))
	binary.BigEndian.PutUint32(srcIPbyte[0:4], inet_addr(srcIp))
	ipHeader = ipHeader.makeIpHeader(srcIPbyte, destIPbyte)
	ipByte, _ := ipHeader.Marshal()
	buffs := make([]byte, 0)
	buffs = append(buffs, ipByte...)
	buffs = append(buffs, tcpByte...)
	addr := syscall.SockaddrInet4{
		Port: destPort,
		//Addr: ip,
	}
	copy(addr.Addr[:4], destIp)
	fmt.Printf("Sendto %v %v ", destIp, destPort)
	err := syscall.Sendto(fd, buffs, 0, &addr)
	if err != nil {
		fmt.Println("Sendto error ", err)
	}
}

func synStart()  {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Println("socket err ",err)
		return
	}
	//设置IP层信息，使其能够修改IP层数据
	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, 2, 1)
	if err != nil {
		fmt.Println("ip err ",err)
		return
	}
	work(fakeIp(),"127.0.0.1",80,2,fd)//todo ack模式
}

func main(){
	synStart()
}
