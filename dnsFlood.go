package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"strings"
	"syscall"
)

type UdpHeader struct {
	srcPort uint16
	dstPort uint16
	len uint16
	checksum uint16
}

func (udp UdpHeader)makeUdpHeader(dp int)UdpHeader{
	udp.srcPort=uint16(rand.Intn(65535)+1)
	udp.dstPort=uint16(dp)

	return udp
}

func packUdpHeader(udpHeader UdpHeader,pstHeader PsdHeader)(UdpHeader,bytes.Buffer){
	var buffer bytes.Buffer
	err:=binary.Write(&buffer, binary.BigEndian, pstHeader)
	if err!=nil{
		fmt.Println("pstHeader err is ",err)
	}
	err=binary.Write(&buffer, binary.BigEndian, udpHeader)
	if err!=nil{
		fmt.Println("tcpHeader err is ",err)
	}
	udpHeader.checksum = checkSum(buffer.Bytes())
	buffer.Reset()
	err=binary.Write(&buffer, binary.BigEndian, udpHeader)
	if err!=nil{
		fmt.Println("reset udpHeader err is ",err)
	}
	return udpHeader,buffer
}

func makeIpPacket(destIp string,srcIp string) ipv4.Header{
	destIPbyte := net.IP(make([]byte, 4))
	binary.BigEndian.PutUint32(destIPbyte[0:4], inet_addr(destIp))
	srcIPbyte := net.IP(make([]byte, 4))
	binary.BigEndian.PutUint32(srcIPbyte[0:4], inet_addr(srcIp))//输入是反的需要check
	iph := ipv4.Header{
		Version: ipv4.Version,
		//IP头长一般是20
		Len:  ipv4.HeaderLen,
		TOS:  0,//0x00
		//buff为数据
		TotalLen: ipv4.HeaderLen ,//加buff
		TTL:  64,
		Flags: ipv4.DontFragment,
		FragOff: 0,
		Protocol: syscall.IPPROTO_UDP,
		Checksum: 0,
		Src:  srcIPbyte,
		Dst:  destIPbyte,
	}

	h, err := iph.Marshal()
	if err != nil {
		fmt.Println(err)
	}
	//计算IP头部校验值
	iph.Checksum = int(checkSum(h))
	return iph
}

type DNSHeader struct {
	ID            uint16
	Flag          uint16
	QuestionCount uint16
	AnswerRRs     uint16 //RRs is Resource Records
	AuthorityRRs  uint16
	AdditionalRRs uint16
}

func (dns DNSHeader)makeDnsPacket()DNSHeader{
	dns.ID=1
	dns.Flag=SetFlag(0,0,0,0,1,0,0)
	dns.QuestionCount=1
	dns.AnswerRRs=0
	dns.AuthorityRRs=0
	dns.AdditionalRRs=0
	return dns
}

func SetFlag(QR uint16, OperationCode uint16, AuthoritativeAnswer uint16, Truncation uint16, RecursionDesired uint16, RecursionAvailable uint16, ResponseCode uint16) uint16 {
	flag:= QR<<15 + OperationCode<<11 + AuthoritativeAnswer<<10 + Truncation<<9 + RecursionDesired<<8 + RecursionAvailable<<7 + ResponseCode
	return flag
}

type DNSQuery struct {
	QuestionType  uint16
	QuestionClass uint16
}

func(dq DNSQuery)makeDqPacket()DNSQuery{
	dq.QuestionClass=1
	dq.QuestionType=255//代表any查询，达到放大的效果
	return dq
}

func ParseDomainName(domain string) ([]byte,error) {
	//要将域名解析成相应的格式，例如：
	//"www.google.com"会被解析成"0x03www0x06google0x03com0x00"
	//就是长度+内容，长度+内容……最后以0x00结尾
	var (
		err error
		buffer   bytes.Buffer
		segments []string = strings.Split(domain, ".")
	)
	for _, seg := range segments {
		err=binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		err=binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	err=binary.Write(&buffer, binary.BigEndian, byte(0x00))

	return buffer.Bytes(),err
}


func workDns(srcIp string,destIp string,destPort int){
	var pstHeader PsdHeader
	var udpHeader UdpHeader
	var ipHeader ipv4.Header
	pstHeader=pstHeader.makePsdHeader(srcIp, destIp,false)
	udpHeader=udpHeader.makeUdpHeader(destPort)
	udpHeader,buffer:=packUdpHeader(udpHeader,pstHeader)
	//fmt.Println(udpHeader)
	udpByte:=buffer.Bytes()
	//fmt.Println(udpByte)
	ipHeader=makeIpPacket(destIp,srcIp)
	ipByte,err:=ipHeader.Marshal()
	if err!=nil{
		fmt.Println("ipHeader.Marshal() err:",err)
	}
	//todo 建立udp连接 ，确定建立的端口是否只为53
	var dnsQuestion DNSQuery
	var dnsHeader DNSHeader
	var dnsBuffer bytes.Buffer
	dnsHeader=dnsHeader.makeDnsPacket()
	dnsQuestion=dnsQuestion.makeDqPacket()
	domainName,err:=ParseDomainName("www.baidu.com")
	err=binary.Write(&dnsBuffer, binary.BigEndian, dnsHeader)
	err=binary.Write(&dnsBuffer, binary.BigEndian, domainName)
	err=binary.Write(&dnsBuffer, binary.BigEndian, dnsQuestion)
	dnsByte:=dnsBuffer.Bytes()



}
