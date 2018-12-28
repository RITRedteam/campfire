# Campfire
Client binary will gather iptable information and send it to a web server via POST request.  The server is a dockerized flask server running the simple app.py file.  

### Config:
In the "campfire.go" file, lines 18 and 19 contain the variables that must be changed  
`var serv = "127.0.0.1:5000"`  
`var loop_time = 60`  
The "serv" variable contains the IP/port of the flask server, and the "loop_time" variable is an integer dictating how often the binary will post data back to the server.    
Once these variables are updated, the code can be compiled into a binary which can be dropped on to the hosts.


### Usage:
##### Client
The client binary has two options when running, to loop or to only execute a single time.  Executing with the "-l " argument will cause the binary to post data every X seconds, where X is defined as explained above. Executing with "-s" argument will cause the binary to only post data back once, then terminate.

##### Server
###### Install
The docker image can be built from the "server" directory by using the following command:  
`sudo docker build -t campfire:latest .`  
once the build is finished, it can be run with the following:  
`sudo docker run -d -p 5000:5000 campfire`  
The flask server is now accessible via port 5000  

###### Navigation
Once the server is running, the hosts will be able to send their post requests to "_IP_:5000/api/rule_send" which can then be viewed by going to the associated host URL.  For hostname "test1" the url would be ""_IP_:5000/api/hosts/test1" or users can navigate to "_IP_:5000/hosts/" and they will have a list of all tracked targets. 

###### Tracker
As an alternative to using a web browser to manually navigate the pages, users can run the "tracker.py" script, which offers a CLI to list tracked hosts and view firewall information of a specific host.  The only configuration needed for "tracker.py" is to install the requirements(the requests library), and to replace the "server" variable on line 10 with the IP address of your campfire flask server.  
`server = "127.0.0.1:5000"`  
 

