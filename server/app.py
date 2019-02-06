"""
This file will be run as the app for a flask server, it's responsible
for recieving/parsing POST'd firewall rules from hosts and writing 
them to files
disclaimer: this barely works
@author: degenerat3
"""
from flask import Flask, request
import os
import datetime

app = Flask(__name__)


"""
This function will recieve the ip/hostnames/rules as a POST request 
and write the data to a file with the timestamp of when the data 
was recieved.
"""
@app.route('/api/rule_send', methods=['GET','POST'])
def get_rules():
	content = request.json
	print content
	hostname = content['hostname']
	hostname = hostname.lower()
        ip = content['ip']
	rules = content['rules']
	if not os.path.isfile("/tmp/host_list.txt"):
		with open("/tmp/host_list.txt", "w+") as f:
			f.write(ip + "\n")
	else:
		with open("/tmp/host_list.txt", "a") as f:
			f.write(ip + "\n")
	if not os.path.exists("/tmp/flask_files"):
		os.makedirs("/tmp/flask_files")		#make dir in /tmp	
	dir_str = "/tmp/flask_files/" + ip + ".txt"
	with open(dir_str, "w+") as f:
            header_str = "IP: " + ip + "\nHostname: " + hostname 
		t = datetime.datetime.now()
		t_str = "Updated at: " + str(t)
		f.write(header_str)		#write hostname/times/padding
		f.write("\n")
		f.write(t_str)
		f.write("\n\n\n")
		f.write(rules)			#write firewall rules
	return str(content)


"""
This function takes the hostname from the url and outputs 
that host's iptables as well as a timestamp of when the 
data was last updated
"""
@app.route('/api/hosts/<host>')
def view_rules(host):
	try:
		dir_str = "/tmp/flask_files/" + host + ".txt"
		with open(dir_str, "r") as f:
			data = f.readlines()	#read rules into string
	except:
		err = "No data on host: " + host
		return err					#err if no matching file

	fin = "<p>"
	for line in data:			#string -> html
		fin += line + "<br />"
	fin += "</p>"
	return fin					#return html block


"""
This function displays all tracked hosts with links
to their unique URLs
"""
@app.route('/hosts/')
def view_hosts():
	if not os.path.exists("/tmp/host_list.txt"):
		return "No tracked hosts..."
	with open("/tmp/host_list.txt", "r") as f:
		lst1 = f.readlines()
		lst = []
		for item in lst1:
			if item not in lst:
				lst.append(item)
		lst = sorted(lst)
		s = "<b>Tracked Targets: </b> <br />"
		for h in lst:
			host = h.strip("\n")
			tmp = "<a href=\"/api/hosts/$\">$</a> <br />"
			tmp = tmp.replace("$", host)
			s += tmp
	return str(s)


if __name__ == '__main__':
	app.run(debug=True, host='0.0.0.0')


