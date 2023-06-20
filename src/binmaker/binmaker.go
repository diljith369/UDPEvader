package binmaker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type UserOptions struct {
	AppType           string
	AgentOS           string
	AgentArch         string
	ListenerOS        string
	ListenerArch      string
	Rhost             string
	Rport             string
	TargetBinaryPath  string
	SaveAs            string
	SourceCodePath    string
	DecodedSourceCode string
	UpdatedSourceCode string
}

type IBinaryBuilder interface {
	SourceToBinary(finflag chan string, exepath, gofilepath string)
	GetUpdateSourceCode() string
	ClearAllSource()
	CreateSourceFile()
}

func checkerror(err error) {
	if err != nil {
		fmt.Println(err)
		//return
	}
}

func (oUserOptions *UserOptions) UpdateSourceCode(isListener bool) {
	oUserOptions.UpdatedSourceCode = strings.Replace(oUserOptions.DecodedSourceCode, "RHOST", oUserOptions.Rhost, 1)
	oUserOptions.UpdatedSourceCode = strings.Replace(oUserOptions.UpdatedSourceCode, "RPORT", oUserOptions.Rport, 1)
	if !isListener {
		oUserOptions.SourceCodePath = `out` + string(os.PathSeparator) + `udptest.go`
		if oUserOptions.AgentOS == "windows" {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs + ".exe"
		} else {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs
		}
	} else {
		//oUserOptions.UpdatedSourceCode = strings.Replace(oUserOptions.DecodedSourceCode, "LPORT", oUserOptions.Rport, 1)
		oUserOptions.SourceCodePath = `out` + string(os.PathSeparator) + `udptestlistener.go`
		if oUserOptions.ListenerOS == "windows" {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs + "listener.exe"
		} else {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs + "listener"
		}
	}
}
func (oUserOptions *UserOptions) CreateSourceFile() {
	newFile, err := os.Create(oUserOptions.SourceCodePath)
	checkerror(err)
	newFile.WriteString(oUserOptions.UpdatedSourceCode)
	newFile.Close()
}
func (oUserOptions *UserOptions) SourceToBinary(isListener bool, finflag chan string) {
	binpath := oUserOptions.TargetBinaryPath + " " + oUserOptions.SourceCodePath
	fmt.Println("building binary.....")
	if runtime.GOOS == "linux" {
		var execargs string
		if isListener {
			execargs = "GOOS=" + oUserOptions.ListenerOS + " GOARCH=" + oUserOptions.ListenerArch + ` go build -ldflags="-s -w" -buildmode=exe -H=windowsgui -o ` + binpath
		} else {
			execargs = "GOOS=" + oUserOptions.AgentOS + " GOARCH=" + oUserOptions.AgentArch + ` go build -ldflags="-s -w" -buildmode=exe -H=windowsgui -o ` + binpath
			if strings.ToLower(oUserOptions.AppType) == "console" {
				execargs = "GOOS=" + oUserOptions.AgentOS + " GOARCH=" + oUserOptions.AgentArch + ` go build -ldflags="-s -w" -buildmode=exe -o ` + binpath
			}
		}
		buildfromlinux(execargs, oUserOptions.TargetBinaryPath)

	} else {
		buildpath := `out` + string(os.PathSeparator) + `build.bat`
		buildbat, err := os.Create(buildpath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buildpath)
		if isListener {
			buildbat.WriteString("SET GOOS=" + oUserOptions.ListenerOS + "\n")
			buildbat.WriteString("SET GOARCH=" + oUserOptions.ListenerArch + "\n")
		} else {
			buildbat.WriteString("SET GOOS=" + oUserOptions.AgentOS + "\n")
			buildbat.WriteString("SET GOARCH=" + oUserOptions.AgentArch + "\n")
		}
		if strings.ToLower(oUserOptions.AppType) == "console" {
			buildbat.WriteString(`go build -o -ldflags="-s -w" -buildmode=exe -o ` + binpath)
		} else {
			buildbat.WriteString(`go build -ldflags="-s -w" -buildmode=exe -H=windowsgui -o ` + binpath)
		}
		buildbat.Close()
		buildfromwindows(buildpath, oUserOptions.TargetBinaryPath)
	}
	finflag <- "Build Success"
}
func buildfromwindows(buildpath string, targetbinpath string) {
	err := exec.Command(buildpath).Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Build Success !")
		fmt.Println(targetbinpath)
	}
}
func buildfromlinux(execargs string, targetbinpath string) {
	cmdpath, _ := exec.LookPath("bash")
	cmd := exec.Command(cmdpath, "-c", execargs)
	res, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Build Success !")
		fmt.Println(targetbinpath)
	}
	fmt.Println(string(res))
}

func (oUserOptions *UserOptions) ClearAllSource() {
	gofiles, err := filepath.Glob("out/*.go")
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range gofiles {
		if err := os.Remove(f); err != nil {
			fmt.Println(err)
		}
	}

	batfiles, err := filepath.Glob("out/*.bat")
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range batfiles {
		if err := os.Remove(f); err != nil {
			fmt.Println(err)
		}
	}
}
