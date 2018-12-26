"""
This file will be run as the app for a flask server, it's responsible
for recieving/parsing POST'd firewall rules from hosts and writing 
them to files
disclaimer: this doesn't work yet
@author: degenerat3
"""
from flask import Flask, request
import os
import datetime

app = Flask(__name__)

@app.route('/api/rule_send', methods=['GET','POST'])
def get_rules():
	content = request.json
	print content
	hostname = content['hostname']
	hostname = hostname.lower()
	rules = content['rules']
	if not os.path.exists("/tmp/flask_files"):
		os.makedirs("/tmp/flask_files")	
	dir_str = "/tmp/flask_files/" + hostname + ".txt"
	with open(dir_str, "w+") as f:
		header_str = "Hostname: " + hostname
		t = datetime.datetime.now()
		t_str = "Updated at: " + str(t)
		f.write(header_str)
		f.write("\n")
		f.write(t_str)
		f.write("\n\n\n")
		f.write(rules)
	return str(content)


@app.route('/api/hosts/<host>')
def view_rules(host):
	try:
		dir_str = "/tmp/flask_files/" + host + ".txt"
		with open(dir_str, "r") as f:
			data = f.readlines()
	except:
		err = "No data on host: " + host
		return err

	fin = "<p>"
	for line in data:
		fin += line + "<br />"
	fin += "</p>"
	return fin


if __name__ == '__main__':
	app.run(debug=True, host='0.0.0.0')
