package listenerprops

import (
	"fmt"
	"net"
)

type UDPListenerProps struct {
	UDPPort        string
	ShellPrmpt     string
	AgentConn      *net.UDPConn
	ReceivedResult []byte
}

type IUDPListenerFuncs interface {
	StartUDPController()
	SetControllerPrompt()
}

const BUFFSIZE = 1024

func (listenerprops *UDPListenerProps) ReadResultandPrint() {
	for {
		chunkbytes, _, _ := listenerprops.AgentConn.ReadFromUDP(listenerprops.ReceivedResult) //fmt.Println(string(recvdcmd[0:n]))
		if chunkbytes < BUFFSIZE {
			fmt.Println(string(listenerprops.ReceivedResult[0:chunkbytes]))
			break
		} else {
			fmt.Println(string(listenerprops.ReceivedResult[0:chunkbytes]))
		}
	}
}

func (listenerprops *UDPListenerProps) SetCutomControllerPrompt() string {
	return listenerprops.ShellPrmpt
}
func (listenerprops *UDPListenerProps) SetDefaultControllerPrompt() string {
	listenerprops.ShellPrmpt = "<<@dcrypT0R~UDP>>"
	return listenerprops.ShellPrmpt
}
func (listenerprops *UDPListenerProps) StartUDPController() (*net.UDPConn, error) {
	udplocaladdr, err := net.ResolveUDPAddr("udp4", listenerprops.UDPPort)
	if err != nil {
		err = fmt.Errorf("Could not resolve the address and port")
		return nil, err
	}
	udpconn, err := net.ListenUDP("udp4", udplocaladdr)
	if err != nil {
		err = fmt.Errorf("UDP Listen error")
		return nil, err
	}
	return udpconn, nil
}
