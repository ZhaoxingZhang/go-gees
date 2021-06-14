package nio

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	go Service()
	time.Sleep(3*time.Second)
	fmt.Println(string(CUrlPOST()))
}

func CUrlPOST() []byte{
	cmd := exec.Command("/bin/sh", "-c", `curl -X POST 'http://127.0.0.1:8089' --data "hello world"`)
	resp, err := cmd.Output()
	if err != nil {
		fmt.Println("Output error ",err, resp)
		os.Exit(1)
	}
	return resp
}