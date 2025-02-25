package containersinfo

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const data = "hello"

func ping(ipAddress string, pingCount int, pingTimeout time.Duration) error {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return err
	}
	defer conn.Close()

	pingAttemptTimeout := pingTimeout / time.Duration(pingCount)

	for i := 0; i < pingCount; i++ {
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  i + 1,
				Data: []byte(data),
			},
		}

		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			return err
		}

		addr := &net.IPAddr{IP: net.ParseIP(ipAddress)}
		_, err = conn.WriteTo(msgBytes, addr)
		if err != nil {
			return err
		}

		reply := make([]byte, 1500)
		conn.SetDeadline(time.Now().Add(pingAttemptTimeout))

		n, _, err := conn.ReadFrom(reply)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		responseMsg, err := icmp.ParseMessage(1, reply[:n])
		if err != nil {
			return err
		}

		switch responseMsg.Type {
		case ipv4.ICMPTypeEchoReply:
			return nil
		case ipv4.ICMPTypeDestinationUnreachable:
			return errDestinationUnreachable
		case ipv4.ICMPTypeRedirect:
			return errRedirect
		case ipv4.ICMPTypeExtendedEchoRequest:
			return errCantReply
		case ipv4.ICMPTypeTimeExceeded:
			continue
		default:
			slog.Info("not expected icmp message type", slog.String("response msg", fmt.Sprintf("%+#v", responseMsg)))
		}
	}

	return errTimeout
}
