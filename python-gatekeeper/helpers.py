import beeline
import json
import random
import re
import time

HEADER_WRITE_KEY = "X-Honeycomb-Team"
HEADER_TIMESTAMP = "X-Honeycomb-Event-Time"
HEADER_SAMPLE_RATE = "X-Honeycomb-Samplerate"

RFC3999_REGEX = re.compile('[0-9-]+T[0-9:]+[^ ]*')
VALID_WRITE_KEY_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

CACHE_TIMEOUT = 10
last_cache_time = None


class Team():
    def __init__(self, id, name, write_key):
        self.id = id
        self.name = name
        self.write_key = write_key


class Dataset():
    def __init__(self, id, name, partition_list):
        self.id = id
        self.name = name
        self.partition_list = partition_list


KNOWN_TEAMS = [
    Team(1, "RPO", "abcd123EFGH"),
    Team(2, "b&w", "ijkl456MNOP"),
    Team(3, "Third", "qrst789UVWX"),
]

KNOWN_DATASETS = [
    Dataset(1, "wade", [1, 2, 3]),
    Dataset(2, "james", []),
    Dataset(3, "helen", [1, 3, 4]),
    Dataset(4, "peter", [1, 2, 4]),
    Dataset(5, "valentine", []),
    Dataset(6, "andrew", [2, 3, 4]),
]


class ParseFailure(Exception):
    pass


class AuthFailure(Exception):
    pass


class AuthMishapenFailure(Exception):
    pass


class DatasetLookupFailure(Exception):
    pass


class SchemaLookupFailure(Exception):
    pass


def timer(wrapped_function):
    """
    Decorator which adds the execution time of wrapped functions as a 
    field on the Honeycomb event.
    """
    def wrapper(*args, **kwargs):
        time_start = time.time()

        result = wrapped_function(*args, **kwargs)

        time_end = time.time()
        event_field_name = "timer.%s_dur_ms" % wrapped_function.__name__
        beeline.add_field(event_field_name, (time_end - time_start) * 1000)
        return result

    return wrapper


@timer
def get_headers(request, event):
    """
    Pulls three headers out of the HTTP request and ensures they are the correct type.
    Does no additional validation.
    """
    # pull raw values from headers
    write_key = request.headers.get(HEADER_WRITE_KEY)
    beeline.add_field(HEADER_WRITE_KEY, write_key)

    timestamp = request.headers.get(HEADER_TIMESTAMP)
    beeline.add_field(HEADER_TIMESTAMP, timestamp)

    sample_rate = request.headers.get(HEADER_SAMPLE_RATE)
    beeline.add_field(HEADER_SAMPLE_RATE, sample_rate)

    # ensure correct types
    # writekeys are strings, so no conversion needed (we will validate them later)
    event['WriteKey'] = write_key

    # timestamps should be in RFC3339 format
    # if not in the right format or if missing, we should note that and continue
    if timestamp:
        match = RFC3999_REGEX.match(timestamp)
        if not match:
            beeline.add_field("error_time_parsing",
                              "timestamp not in RFC3999 format")
        event['Timestamp'] = timestamp
    else:
        beeline.add_field("error_time_parsing", "no timestamp for event")

    # sample rate should be a positive int, defaults to 1 if empty
    if sample_rate == "":
        sample_rate = "1"
    try:
        parsed_sample_rate = int(sample_rate)
        event['SampleRate'] = parsed_sample_rate
        beeline.add_field("sample_rate", parsed_sample_rate)
    except ValueError:
        raise ParseFailure


@timer
def validate_write_key(write_key):
    """
    Ensures writekeys are in a valid format: only contain characters within [a-zA-Z0-9].
    Authenticates writekeys: 
    """
    for char in write_key:
        if char not in VALID_WRITE_KEY_CHARS:
            raise AuthMishapenFailure

    for team in KNOWN_TEAMS:
        if write_key == team.write_key:
            return team
    raise AuthFailure


@timer
def resolve_dataset(dataset_name):
    """
    Authenticates datasets: here we would call out to the database to validate the dataset,
    but in the interests of simplicity, we'll check that it's one of a few hardcoded values.
    """
    for dataset in KNOWN_DATASETS:
        if dataset_name == dataset.name:
            return dataset
    raise DatasetLookupFailure


@timer
def get_partition(dataset):
    """
    Checks for partitions assigned to the dataset and randomly chooses on of the assigned
    partitions.
    """
    partitions = dataset.partition_list
    if len(partitions) < 1:
        raise DatasetLookupFailure

    partition = random.choice(partitions)

    return partition


@timer
def get_schema(dataset):
    """
    Looks up the dataset schema: implements a fake cache and database call to better simulate
    what an actual call might look like. 
    """
    global last_cache_time

    hit_cache = True
    if last_cache_time is None:
        last_cache_time = time.time()
    if time.time() - last_cache_time > CACHE_TIMEOUT:
        # we fall through the cache every 10 seconds
        hit_cache = False
        # pretend to hit a slow database that takes 30-50ms
        time.sleep(random.uniform(.03, .05))
        last_cache_time = time.time()
    beeline.add_field("hitSchemaCache", hit_cache)
    # let's just fail sometimes to pretend
    if random.randint(0, 61) == 0:
        raise SchemaLookupFailure


def write_event(event):
    """
    Writes the event to disk using the chosen partition.
    """
    file_name = "/tmp/api%s.log" % event['ChosenPartition']
    with open(file_name, 'w') as f:
        f.write(json.dumps(event))
