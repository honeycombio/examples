import os
import sys
import datetime
import libhoney
from flask_sqlalchemy import SQLAlchemy
from flask import Flask, jsonify, g, request
import logging

app = Flask(__name__)

app.logger.addHandler(logging.StreamHandler(sys.stderr))
app.logger.setLevel(logging.DEBUG)

db_user = "root"
db_pass = ""
db_host = os.getenv("DB_HOST")
db_name = "example-python-api"
app.config["SQLALCHEMY_DATABASE_URI"] = "mysql://{}:{}@{}:3306/{}".format(db_user, db_pass, db_host, db_name)
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
db = SQLAlchemy(app)

libhoney.init(writekey=os.getenv("HONEYCOMB_WRITEKEY"), dataset="examples.python-api")
libhoney_builder = libhoney.Builder()

class Todo(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    description = db.Column(db.String(200), nullable=False)
    completed = db.Column(db.Boolean(False), nullable=False)
    due = db.Column(db.DateTime())

    def serialize(self):
        return {
            "id": self.id,
            "description": self.description,
            "completed": self.completed,
            "due": self.due,
        }

    def __repr__(self):
        return "<Todo {}>".format(self.id)

db.create_all()

class InvalidUsage(Exception):
    status_code = 400

    def __init__(self, message, status_code=None, payload=None):
        Exception.__init__(self)
        self.message = message
        if status_code is not None:
            self.status_code = status_code
        self.payload = payload

    def to_dict(self):
        rv = dict(self.payload or ())
        rv['message'] = self.message
        return rv

def milliseconds_since(start):
    delta = (datetime.datetime.now() - start).total_seconds()
    return delta*1000

@app.errorhandler(InvalidUsage)
def handle_invalid_usage(error):
    g.ev.add_field("errors.message", error.message)
    response = jsonify(error.to_dict())
    response.status_code = error.status_code
    return response

@app.before_request
def before():
    # g is the thread-local / request-local variable, we will use it to store
    # information that will be used when we eventually send the event to
    # Honeycomb, including a timer for the whole request duration.
    g.req_start = datetime.datetime.now()
    g.ev = libhoney_builder.new_event()
    g.ev.add_field("request.path", request.path)
    g.ev.add_field("request.method", request.method)
    g.ev.add_field("request.user_agent.browser", request.user_agent.browser)
    g.ev.add_field("request.user_agent.platform", request.user_agent.platform)
    g.ev.add_field("request.user_agent.language", request.user_agent.language)
    g.ev.add_field("request.user_agent.string", request.user_agent.string)
    g.ev.add_field("request.user_agent.version", request.user_agent.version)
    g.ev.add_field("request.python_function", request.endpoint)
    g.ev.add_field("request.url_pattern", str(request.url_rule))

@app.after_request
def after(response):
    g.ev.add_field("response.status_code", response.status_code)

    # Note that this isn"t the total time to serve the request, i.e., how long
    # the end user is waiting. It accounts for the time spent in the Flask
    # handlers but not Werkzeug, etc. Ingesting edge data from ELB or
    # nginx etc. is usually much better for that kind of (total request time) info.
    g.ev.add_field("timers.flask_time_ms", milliseconds_since(g.req_start))

    app.logger.debug(g.ev)
    g.ev.send()
    return response

@app.route("/")
def index():
    return jsonify(up=True)

@app.route("/todos", defaults={"todo_id": None}, methods=["GET", "POST"])
@app.route("/todos/<int:todo_id>", methods=["GET", "DELETE", "PUT"])
def todo(todo_id):
    if todo_id is None:
        if request.method == "POST":
            json = request.get_json()
            todo = Todo(
                description=json.get("description", ""),
                completed=json.get("completed", False),
                due=datetime.datetime.fromtimestamp(json.get("due", None)),
            )

            insert_todo_start = datetime.datetime.now()

            db.session.add(todo)
            db.session.commit()

            g.ev.add_field("timers.db.insert_todo_ms", milliseconds_since(insert_todo_start))

            # Augment events with high cardinality information!
            g.ev.add_field("todo.id", todo.id)

            return jsonify(todo.serialize())

        if request.method == "GET":
            select_all_todos_start = datetime.datetime.now()
            todos = [todo.serialize() for todo in Todo.query.all()]
            g.ev.add_field("timers.db.select_all_todos", milliseconds_since(select_all_todos_start))
            return jsonify(todos)

        raise InvalidUsage("Method not allowed", status_code=405)
    else:
        # Augment events with high cardinality information!
        g.ev.add_field("todo.id", todo_id)

        select_todo_start = datetime.datetime.now()
        todo = Todo.query.get(todo_id)
        g.ev.add_field("timers.db.select_todo", milliseconds_since(select_todo_start))

        if todo is None:
            err = "Todo not found"
            g.ev.add_field("errors.message", err)
            return jsonify({"error": err}), 404

        if request.method == "GET":
            return jsonify(todo.serialize())

        elif request.method == "DELETE":
            delete_todo_start = datetime.datetime.now()
            db.session.delete(todo)
            db.session.commit()
            g.ev.add_field("timers.db.delete_todo", milliseconds_since(delete_todo_start))

            return jsonify({"success": True, "id": todo_id})

        elif request.method == "PUT":
            json = request.get_json()
            todo.description = json.get("description", todo.description)
            todo.completed = json.get("completed", todo.completed)
            app.logger.debug(json.get("completed", todo.completed))
            todo.due = json.get("due", todo.due),

            update_todo_start = datetime.datetime.now()
            db.session.add(todo)
            db.session.commit()
            g.ev.add_field("timers.db.update_todo", milliseconds_since(update_todo_start))

            return jsonify(todo.serialize())

        raise InvalidUsage("Method not allowed", status_code=405)
