"""
This script will hit the Campfire API and function as a CLI
to view tracked hosts' firewall configurations
@author: degenerat3
"""

import requests
import re

server = "127.0.0.1:5000"		#IP of flask server

"""
query API and return string of all tracked hosts
"""
def list_hosts(srv):
	endpoint = "http://" + srv + "/hosts/"
	r = requests.get(endpoint)
	txt = r.text
	txt = txt.split("<br />")
	hosts = []
	for item in txt:
		if "<b>" in item:
			continue
		if "hosts" in item:
			res = re.search('/hosts/(.*)">', item)
			h = res.group(1)
			hosts.append(h)
	hstr = ""
	for item in hosts:
		hstr += item + "\n"	
	return hstr

"""
query api, return firewall config of host
"""
def get_host(srv, host):
	endpoint = "http://" + srv + "/api/hosts/" + host
	r = requests.get(endpoint)
	txt = r.text
	if "No data on host" in txt:
		return "404"
	txt = txt.split("<br />")
	
	a = "\n"
	for line in txt:
		#print(line)
		line = line.replace("<p>", '')
		line = line.replace("</p>", '')
		a += line
	return a

"""
Display help/syntax message to user
"""
def help_msg():
	s = """	'help': print help message
	'hosts':list tracked hosts
	'show X': show IP tables for hostname X
	'exit': close the program"""	
	return s	


"""
Prase/process user input
"""
def inp_loop(srv):
	inp = input("Tracker> ")
	if inp == "help":
		st = help_msg()
	elif inp == "hosts":
		st = list_hosts(srv)
	elif "show" in inp:
		h = inp.lstrip("show ")
		r = get_host(srv, h)
		if r == "404":
			st = "No data for host: " + h
		else:
			st = r
	elif inp == "":
		return
	elif inp == "exit":
		exit()
	else:
		st = "Invalid input, type 'help' for assistance"
	
	print(st)	
		
"""
Run it all
"""
print("Welcome to the Campfire Tracker!")
print(help_msg())
while True:
	inp_loop(server)


