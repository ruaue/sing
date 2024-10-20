package network

import (
	"net"

	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

type HandshakeFailure interface {
	HandshakeFailure(err error) error
}

type HandshakeSuccess interface {
	HandshakeSuccess() error
}

type ConnHandshakeSuccess interface {
	ConnHandshakeSuccess(conn net.Conn) error
}

type PacketConnHandshakeSuccess interface {
	PacketConnHandshakeSuccess(conn net.PacketConn) error
}

func ReportHandshakeFailure(reporter any, err error) error {
	if handshakeConn, isHandshakeConn := common.Cast[HandshakeFailure](reporter); isHandshakeConn {
		return E.Append(err, handshakeConn.HandshakeFailure(err), func(err error) error {
			return E.Cause(err, "write handshake failure")
		})
	}
	return err
}

func CloseOnHandshakeFailure(reporter any, onClose CloseHandler, err error) error {
	if handshakeConn, isHandshakeConn := common.Cast[HandshakeFailure](reporter); isHandshakeConn {
		err = E.Append(err, handshakeConn.HandshakeFailure(err), func(err error) error {
			return E.Cause(err, "write handshake failure")
		})
	} else {
		if tcpConn, isTCPConn := common.Cast[interface {
			SetLinger(sec int) error
		}](reporter); isTCPConn {
			tcpConn.SetLinger(0)
		}
		common.Close(reporter)
	}
	if onClose != nil {
		onClose(err)
	}
	return err
}

// Deprecated: use ReportConnHandshakeSuccess/ReportPacketConnHandshakeSuccess instead
func ReportHandshakeSuccess(reporter any) error {
	if handshakeConn, isHandshakeConn := common.Cast[HandshakeSuccess](reporter); isHandshakeConn {
		return handshakeConn.HandshakeSuccess()
	}
	return nil
}

func ReportConnHandshakeSuccess(reporter any, conn net.Conn) error {
	if handshakeConn, isHandshakeConn := common.Cast[ConnHandshakeSuccess](reporter); isHandshakeConn {
		return handshakeConn.ConnHandshakeSuccess(conn)
	}
	if handshakeConn, isHandshakeConn := common.Cast[HandshakeSuccess](reporter); isHandshakeConn {
		return handshakeConn.HandshakeSuccess()
	}
	return nil
}

func ReportPacketConnHandshakeSuccess(reporter any, conn net.PacketConn) error {
	if handshakeConn, isHandshakeConn := common.Cast[PacketConnHandshakeSuccess](reporter); isHandshakeConn {
		return handshakeConn.PacketConnHandshakeSuccess(conn)
	}
	if handshakeConn, isHandshakeConn := common.Cast[HandshakeSuccess](reporter); isHandshakeConn {
		return handshakeConn.HandshakeSuccess()
	}
	return nil
}
