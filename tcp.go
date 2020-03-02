package main

//type TCPHeader struct {
//	SrcPort   uint16
//	DstPort   uint16
//	SeqNum    uint32
//	AckNum    uint32
//	Offset    uint8
//	Flag      uint8
//	Window    uint16
//	Checksum  uint16
//	UrgentPtr uint16
//}
//
//type PsdHeader struct {
//	SrcAddr   uint32
//	DstAddr   uint32
//	Zero      uint8
//	ProtoType uint8
//	TcpLength uint16
//}
//
//func (tcp TCPHeader)makeTcpHeader(destPort int ,flag int) TCPHeader{
//	rand.Seed(time.Now().UnixNano())
//	tcp.SrcPort=uint16(rand.Intn(65535)+1)
//	tcp.DstPort=uint16(destPort)
//	tcp.SeqNum=uint32(rand.Intn(4000000000)+1)
//	if flag==2{
//		tcp.AckNum=0
//	}else{
//		tcp.AckNum=uint32(rand.Intn(4000000000)+1)
//	}
//	tcp.Flag=uint8(flag)
//	tcp.Offset=uint8(uint16(unsafe.Sizeof(TCPHeader{}))/4) << 4
//	tcp.Window=uint16(8192)
//	tcp.UrgentPtr=uint16(0)
//	return tcp
//}
//
//func (psd PsdHeader)makePsdHeader(srcIp string,destIp string) PsdHeader{
//	psd.SrcAddr=inet_addr(srcIp)
//	psd.DstAddr=inet_addr(destIp)
//	psd.ProtoType= syscall.IPPROTO_TCP//6代表tcp
//	psd.TcpLength=uint16(20)
//	psd.Zero=0
//	return psd
//}
