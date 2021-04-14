package ops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"

	"github.com/nanovms/ops/types"
)

const (
	basicProgram = `package main
import(
	"fmt"
)
func main(){
	fmt.Println("hello world")
}
`
)

// BuildBasicProgram generates binary from a hello world golang program in the current directory
func BuildBasicProgram() (binaryPath string) {
	program := []byte(basicProgram)
	binaryPath = fmt.Sprintf("basic%d.go", rand.Intn(100))
	sourcePath := binaryPath + ".go"

	err := ioutil.WriteFile(sourcePath, program, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer os.Remove(sourcePath)

	cmd := exec.Command("go", "build", sourcePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return
}

// WriteConfigFile writes configuration passed by argument to a file
func WriteConfigFile(config *types.Config) (filepath string) {
	filepath = "./config.json"
	json, _ := json.MarshalIndent(config, "", "  ")

	err := ioutil.WriteFile(filepath, json, 0666)
	if err != nil {
		panic(err)
	}

	return
}
