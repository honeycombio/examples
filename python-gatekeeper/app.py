import json
import datetime
import dateutil.parser
import logging
import os

from flask import Flask, request, render_template
from helpers import *

log = logging.getLogger(__name__)

PARSE_FAILURE_RESPONSE = (
    '{"error":"unable to parse request headers"}', 400, None)
AUTH_FAILURE_RESPONSE = (
    '{"error":"writekey didn\'t match valid credentials"}', 401, None)
AUTH_MISHAPEN_FAILURE_RESPONSE = (
    '{"error":"writekey malformed - expect only letters and numbers"}', 400, None)
JSON_FAILURE_RESPONSE = (
    '{"error":"failed to unmarshal JSON body"}', 400, None)
DATASET_LOOKUP_FAILURE_RESPONSE = (
    '{"error":"failed to resolve dataset object"}', 400, None)
SCHEMA_LOOKUP_FAILURE_RESPONSE = (
    '{"error":"failed to resolve schema"}', 500, None)

honeycomb_write_key = os.environ.get("HONEYCOMB_WRITEKEY")
if not honeycomb_write_key:
    log.error(
        "Got empty writekey from the environment. Please set HONEYCOMB_WRITEKEY")

app = Flask(__name__)

@app.route('/')
def home():
    return render_template('home.html')


@app.route('/x/alive')
def health():
    return json.dumps({'alive': 'yes'})


@app.route('/1/events/<dataset_name>', methods=['POST'])
def handle_event(dataset_name):
    event = {}

    # parse JSON body
    try:
        data = json.loads(request.data)
        event['Data'] = data
    except (TypeError, json.decoder.JSONDecodeError):
        return JSON_FAILURE_RESPONSE

    # get writekey, timestamp, and sample rate out of HTTP headers
    try:
        get_headers(request, event)
    except ParseFailure:
        return PARSE_FAILURE_RESPONSE

    # authenticate writekey or return 401
    try:
        team = validate_write_key(event['WriteKey'])
    except AuthFailure:
        return AUTH_FAILURE_RESPONSE
    except AuthMishapenFailure:
        return AUTH_MISHAPEN_FAILURE_RESPONSE

    # take the writekey and the dataset name and get back a dataset object
    try:
        dataset = resolve_dataset(dataset_name)
    except DatasetLookupFailure:
        return DATASET_LOOKUP_FAILURE_RESPONSE

    # get partition info
    try:
        partition = get_partition(dataset)
        event['ChosenPartition'] = partition
    except DatasetLookupFailure:
        return DATASET_LOOKUP_FAILURE_RESPONSE

    # check time - set to now if not present
    if 'Timestamp' not in event:
        event['Timestamp'] = datetime.datetime.now(
            datetime.timezone.utc).isoformat()
    else:
        # record the difference between the event's timestamp and now to help identify
        # lagging events
        event_timestamp = dateutil.parser.parse(event['Timestamp'])
        event_time_delta = datetime.datetime.now(
            datetime.timezone.utc) - event_timestamp

    # verify schema
    try:
        get_schema(dataset)
    except SchemaLookupFailure:
        return SCHEMA_LOOKUP_FAILURE_RESPONSE

    # hand off to external service - write to local disk
    write_event(event)
    return ''
