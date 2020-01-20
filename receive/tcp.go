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
	`net`
	`sync`
	`sync/atomic`

	gotcp `github.com/generalzgd/securegotcp`

	`github.com/generalzgd/msg-subscriber/iface`
)

type TcpProtocol struct {
}

func (p *TcpProtocol) ReadPacket(conn net.Conn) (gotcp.Packet, error) {
	panic("implement me")

}
//
type TcpConn struct {
	conn *gotcp.Conn
}

func (p *TcpConn) Send(packet gotcp.Packet) error {
	panic("implement me")
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
	Svr *gotcp.Server
	//
	lock    sync.RWMutex
	linkMap map[uint32]iface.IPackSender
	seed    uint32
}


func (p *TcpReceiver) OnConnect(conn *gotcp.Conn) bool {
	//panic("implement me")
	p.lock.Lock()
	defer p.lock.Unlock()

	id := atomic.AddUint32(&p.seed, 1)
	conn.PutExtraData(id)
	c := &TcpConn{conn:conn}
	p.linkMap[id] = c
	//
	p.Callback.OnConnect(c)
	return true
}

func (p *TcpReceiver) OnMessage(conn *gotcp.Conn, pack gotcp.Packet) bool {
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
	p.lock.Lock()
	defer p.lock.Unlock()

	id := conn.GetExtraData().(uint32)
	if c, ok := p.linkMap[id]; ok {
		p.Callback.OnClose(c)
		delete(p.linkMap, id)
	}
}
