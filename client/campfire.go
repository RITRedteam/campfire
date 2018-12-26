// This go file will collect the host's firewall rules ship them back
// to a defined webserver, along with the hostname
// disclaimer: this barely works
// @author: degenerat3
package main

import "os/exec"
import "fmt"
import "log"
import "net/http"
import "io/ioutil"
import "bytes"
import "encoding/json"
import "strings"
import "time"
import "os"

var serv = "127.0.0.1:5000/api/rule_send"	//IP of flask serv
var loop_time = 60		//sleep time in seconds


// return output of "iptables -L" as one large string
func get_tables() string{
	cmd := exec.Command("iptables", "-L")
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
		return "Err"
    }
	return string(out)
}


// return hostname as string
func get_hn() string{
	cmd := exec.Command("hostname")
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
        return "Err"
    }
	o1 := string(out)
	o2 := strings.TrimSuffix(o1, "\n")
    return o2

}


// post strings to flask server
func send_data(rules string, host string){
	url1 := "http://" + serv;	// turn ip into valid url
    jsonData := map[string]string{"hostname": host, "rules": rules}
	jsonValue, _ := json.Marshal(jsonData)
	resp, err := http.Post(url1, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("Req failed: %s\n", err)
		return
	} else{
		// block below for debug
		// data, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(data))
		return
	}
}


// fetch data then send it
func run(){
	rules := get_tables()
	host := get_hn()
	send_data(rules, host)
}

func main(){
	loop_arg := os.Args[1]
	if loop_arg == "-s"{	// if "-s" is an arg, run once, otherwise loop
		run()
	} else{
		for {				// send data to webserver ever X seconds, until termination
			run()
			t := time.Duration(loop_time*1000)
			time.Sleep(t * time.Millisecond)
		}
	}
}


