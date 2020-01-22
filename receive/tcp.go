/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: tcp.go
 * @time: 2020/1/19 8:05 下午
 * @project: msgsubscribesvr
 */

package receive

import (
	`encoding/binary`
	`errors`
	`fmt`
	`io`
	`net`
	`sync`
	`sync/atomic`
	`time`

	`github.com/astaxie/beego/logs`
	gotcp `github.com/generalzgd/securegotcp`

	`github.com/generalzgd/msg-subscriber/codec`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/util`
)

type TcpProtocol struct {
	HeadSize  int
	LenPos    int
	LenSize   int
	HeadCodec codec.IHeadCodec
	BodyCodec codec.IBodyCodec
}

func (p *TcpProtocol) ReadPacket(conn net.Conn) (gotcp.Packet, error) {
	var (
		lengthBytes = make([]byte, p.HeadSize)
		length      int
	)

	// 设置读超时
	// conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	// defer conn.SetReadDeadline(time.Time{})

	// read header
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}

	length = util.GetIntFromBuf(lengthBytes, p.LenPos, p.LenSize, binary.LittleEndian)
	// 	length := uint32(0)
	if length > 1024*32 {
		return nil, errors.New(fmt.Sprintf("the size of post packet is larger than the limit %d.", 1024*32))
	}

	buff := make([]byte, p.HeadSize+length)
	copy(buff[0:p.HeadSize], lengthBytes)

	if _, err := io.ReadFull(conn, buff[p.HeadSize:]); err != nil {
		return nil, err
	}

	pack := codec.NewDataPack(buff, p.HeadCodec, p.BodyCodec)
	return pack, nil
}

//
type TcpConn struct {
	conn *gotcp.Conn
}

func (p *TcpConn) Send(packet gotcp.Packet) error {
	if p.conn != nil {
		p.conn.AsyncWritePacket(packet, time.Second*3)
	}
	return nil
}

func (p *TcpConn) GetAddress() string {
	return p.conn.GetRawConn().RemoteAddr().String()
}

func (p *TcpConn) GetConnId() uint32 {
	return p.conn.GetExtraData().(uint32)
}

type IReceiveCallback interface {
	OnConnect(conn iface.IPackSender)
	OnClose(conn iface.IPackSender)
	OnMessage(conn iface.IPackSender, pack gotcp.Packet)
}

//
type TcpReceiver struct {
	Callback IReceiveCallback
	//
	svr     *gotcp.Server
	lock    sync.RWMutex
	linkMap map[uint32]iface.IPackSender
	exclude map[string]struct{}
	seed    uint32
}

func NewTcpReceiver(callback IReceiveCallback) *TcpReceiver {
	return &TcpReceiver{linkMap: map[uint32]iface.IPackSender{}, Callback: callback, exclude: map[string]struct{}{}}
}

func (p *TcpReceiver) SetExcludeIp(ips ...string) {
	for _, it := range ips {
		p.exclude[it] = struct{}{}
	}
}

func (p *TcpReceiver) Start(cfg *gotcp.Config, addr string, pro gotcp.Protocol) error {
	svr := gotcp.NewServer(cfg, p, pro)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		logs.Error("Resolve tcp addr fail: ", tcpAddr)
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logs.Warning("Listen tcp fail: ", tcpAddr.String())
		return err
	}

	go svr.Start(listener, time.Second)

	logs.Info("Start listen:", listener.Addr())
	return nil
}

func (p *TcpReceiver) OnConnect(conn *gotcp.Conn) bool {
	host, _, err := net.SplitHostPort(conn.GetRawConn().RemoteAddr().String())
	if err != nil {
		return false
	}
	if _, ok := p.exclude[host]; ok {
		return true
	}
	//logs.Debug("OnConnect() %v", conn.GetRawConn().RemoteAddr())
	//panic("implement me")
	p.lock.Lock()
	defer p.lock.Unlock()

	id := atomic.AddUint32(&p.seed, 1)
	conn.PutExtraData(id)
	c := &TcpConn{conn: conn}
	p.linkMap[id] = c
	//
	p.Callback.OnConnect(c)
	return true
}

func (p *TcpReceiver) OnMessage(conn *gotcp.Conn, pack gotcp.Packet) bool {
	//logs.Debug("OnMessage() %v", conn.GetRawConn().RemoteAddr())
	//panic("implement me")
	id := conn.GetExtraData().(uint32)
	p.lock.RLock()
	defer p.lock.RUnlock()

	if c, ok := p.linkMap[id]; ok {
		p.Callback.OnMessage(c, pack)
	}
	return true
}

func (p *TcpReceiver) OnClose(conn *gotcp.Conn) {
	host, _, err := net.SplitHostPort(conn.GetRawConn().RemoteAddr().String())
	if err != nil {
		return
	}
	if _, ok := p.exclude[host]; ok {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	id := conn.GetExtraData().(uint32)
	if c, ok := p.linkMap[id]; ok {
		p.Callback.OnClose(c)
		delete(p.linkMap, id)
	}
}
