import os
import sys
import datetime
import libhoney
from flask_sqlalchemy import SQLAlchemy
from flask import Flask, jsonify, g, request
import logging

db_user = "root"
db_pass = ""
db_host = os.getenv("DB_HOST")
db_name = "example-python-api"
app = Flask(__name__)
app.logger.addHandler(logging.StreamHandler(sys.stderr))
app.logger.setLevel(logging.DEBUG)
app.config["SQLALCHEMY_DATABASE_URI"] = "mysql://{}:{}@tcp({}:3306)/{}".format(db_user, db_pass, db_host, db_name)
libhoney.init(writekey=os.getenv("HONEYCOMB_WRITEKEY"), dataset="examples.python-api")
libhoney_builder = libhoney.Builder()

def milliseconds_since(start):
    delta = (datetime.datetime.now() - start).total_seconds()
    app.logger.debug(delta)
    return delta*1000

@app.before_request
def before():
    # g is the thread-local / request-local variable, we will use it to store
    # information that will be used when we eventually send the event to
    # Honeycomb, including a timer for the whole request duration.
    g.req_start = datetime.datetime.now()
    g.ev = libhoney_builder.new_event()
    g.ev.add_field("request.path", request.path)
    g.ev.add_field("request.method", request.method)
    g.ev.add_field("request.user_agent", request.user_agent)
    g.ev.add_field("request.endpoint", request.endpoint)
    g.ev.add_field("request.url_pattern", str(request.url_rule))

@app.after_request
def after(response):
    g.ev.add_field("response.status_code", response.status_code)

    # Note that this isn't _strictly_ the total time to serve the request - it
    # accounts only for the time spent in the Flask handlers and not Werkzeug,
    # etc.
    g.ev.add_field("timers.flask_time_ms", milliseconds_since(g.req_start))

    app.logger.debug(g.ev)
    g.ev.send()
    return response

@app.route("/")
def index():
    return jsonify(up=True)

@app.route('/todos', defaults={"todo_id": None})
@app.route("/todos/<int:todo_id>")
def todo(todo_id):
    if todo_id is None:
        return jsonify([{"id": 1, "description": "wash hair"}, {"id": 2, "description": "walk dog"}, {"id": 3, "description": "buy flowers"}])
    else:
        return jsonify({})
