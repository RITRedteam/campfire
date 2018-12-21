package main

import "os/exec"
import "fmt"
import "bytes"
import "log"

func get_tables(){
	cmd := exec.Command("iptables", "-F")
	var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
    fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

func main(){
	get_tables()
}
