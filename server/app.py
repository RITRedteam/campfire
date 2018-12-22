from flask import Flask
app = Flask(__name__)

@app.route('/api/rule_send')
def get_rules():
	return "we got some rules"


