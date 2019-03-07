# Campfire
Client binary will gather iptable information and send it to a web server via POST request.  The server will be [Reach](https://github.com/degenerat3/reach), which will process the data and ship it to an ELK stack.

### Config:
In the "campfire.go" file, lines 18 and 19 contain the variables that must be changed  
`var serv = "127.0.0.1:5000"`  
`var loop_time = 60`  
The "serv" variable contains the IP/port of the Reach server, and the "loop_time" variable is an integer dictating how often the binary will post data back to the server (in seconds).    
Once these variables are updated, the code can be compiled into a binary using the following command, which can then be dropped on to the hosts:  
`go build campfire.go`  
Note: Target OS must be set to Linux (GOOS env variable)


### Usage:
##### Client
The client binary has two options when running, to loop or to only execute a single time.  Executing with no arguments will cause the binary to post data every X seconds, where X is defined as explained above. Executing with any argument will cause the binary to only post data back once, then terminate.  

### Old Version:
There is a standalone tracker/server that can be used, which runs a simple flask app for callbacks instead of using Reach.  Look in the "old" branch to figure it out.
