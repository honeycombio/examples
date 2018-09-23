require 'sinatra'
require 'rubygems'
require 'honeycomb-beeline'
require 'httparty'
require 'time'
require 'benchmark'
require 'json'
require './app/helpers/dataset_helper'
require './app/helpers/event_helper'
require './app/helpers/team_helper'
require './app/helpers/error_helpers'
require './app/helpers/method_helpers'

include DatasetHelper
include EventHelper
include TeamHelper
include ErrorHelpers
include MethodHelpers

Honeycomb.init(debug: true)

# Initalize Datasets
datasets = [
  { id: 1, name: 'wade',      partition_list: [1, 2, 3] },
  { id: 2, name: 'james',     partition_list: [1, 2, 4] },
  { id: 3, name: 'helen',     partition_list: [1, 3, 4] },
  { id: 4, name: 'peter',     partition_list: [1, 2, 4] },
  { id: 5, name: 'valentine', partition_list: [1, 2, 4] },
  { id: 6, name: 'andrea',    partition_list: [2, 3, 4] }
]

$known_datasets = datasets.map { |dataset| Dataset.new(dataset) }
# Initalize Teams
teams = [
  { id: 1, name: 'RPO',   write_key: 'abcd123EFGH' },
  { id: 2, name: 'B&W',   write_key: 'ijkl456MNOP' },
  { id: 3, name: 'Third', write_key: 'qrst789UVWX' }
]

$known_teams = teams.map { |team| Team.new(team) }

HEADER_WRITE_KEY   = 'HTTP_X_HONEYCOMB_TEAM'.freeze
HEADER_TIMESTAMP   = 'HTTP_X_HONEYCOMB_EVENT_TIME'.freeze
HEADER_SAMPLE_RATE = 'HTTP_X_HONEYCOMB_SAMPLERATE'.freeze
CACHE_TIMEOUT      = 10

post '/1/events/:dataset_name' do
  content_type :json
  # get writekey, timestamp, and sample rate from HTTP Headers
  write_key = request.get_header(HEADER_WRITE_KEY)
  timestamp = request.get_header(HEADER_TIMESTAMP)
  sample_rate = request.get_header(HEADER_SAMPLE_RATE)

  Rack::Honeycomb.add_field(env, HEADER_SAMPLE_RATE, sample_rate)
  Rack::Honeycomb.add_field(env, HEADER_WRITE_KEY, write_key)
  Rack::Honeycomb.add_field(env, HEADER_TIMESTAMP, timestamp)

  users_dataset = params[:dataset_name]

  sample_rate = '1' if sample_rate.nil?

  # add error handling for bad sample rate
  begin
    positive_integer?(sample_rate)
    int_sample_rate = sample_rate.to_i
  rescue Bad_Sample_Rate => e
    status 400
    return e.message.to_json
  end

  # Initialize new event and set some fields
  event = Event.new({})
  event.timestamp   = timestamp
  event.sample_rate = int_sample_rate
  event.write_key   = write_key

  Rack::Honeycomb.add_field(env, 'sample_rate', int_sample_rate)

  # parse JSON body
  json_parse_timer_start = (Time.now.to_f * 1000).to_i

  begin
    parsed_json = JSON.parse(request.body.read)
  rescue StandardError
    status 400
    return { "error": 'unable to parse request headers' }.to_json
  end

  json_parse_timer_end = (Time.now.to_f * 1000).to_i
  json_parse_timer_ms  = (json_parse_timer_start - json_parse_timer_end)

  event.data = parsed_json
  Rack::Honeycomb.add_field(env, 'timer.parse_json_dur_ms', json_parse_timer_ms)

  # authenticate the writekey

  write_key_timer_start = (Time.now.to_f * 1000).to_i

  begin
    current_team = validate_write_key(write_key)
  rescue Auth_Mishapen_Failure => e
    status 401
    return e.message.to_json
  rescue Auth_Failure => e
    status 400
    return e.message.to_json
  end

  write_key_timer_end = (Time.now.to_f * 1000).to_i
  write_key_timer_ms  = (write_key_timer_end - write_key_timer_start)

  Rack::Honeycomb.add_field(env, 'team', current_team)
  Rack::Honeycomb.add_field(env, 'timer.validated_writekey_dur_ms', write_key_timer_ms)

  # take the writekey and the dataset name and get back a dataset object

  dataset_timer_start = (Time.now.to_f * 1000).to_i

  begin
    current_dataset = resolve_dataset(users_dataset)
  rescue Dataset_Lookup_Failure => e
    status 400
    return e.message.to_json
  end

  dataset_timer_end = (Time.now.to_f * 1000).to_i
  dataset_time_ms   = (dataset_timer_end - dataset_timer_start)

  Rack::Honeycomb.add_field(env, 'dataset', current_dataset)
  Rack::Honeycomb.add_field(env, 'timer.resolve_dataset_timer_dur_ms', dataset_time_ms)

  # get partition info - stub about

  grab_partition_timer_start = (Time.now.to_f * 1000).to_i

  begin
    chosen_partition = grab_partition(current_dataset)
  rescue Dataset_Lookup_Failure => e
    status 405
    return e.message.to_json
  end

  grab_partition_timer_end = (Time.now.to_f * 1000).to_i
  grab_partition_timer_ms  = (grab_partition_timer_start - grab_partition_timer_end)

  event.chosen_partition = chosen_partition
  Rack::Honeycomb.add_field(env, 'chosen_partition', chosen_partition)
  Rack::Honeycomb.add_field(env, 'timer.grab_partition_dur_ms', grab_partition_timer_ms)

  # check time - use or set to now if broken
  event_time_delta = 0

  if event.timestamp.zero?
    event.timestamp = (Time.now.to_f * 1000).to_i
  else
    event_time_delta = (Time.now.to_i - event.timestamp.to_i)
  end

  Rack::Honeycomb.add_field(env, 'event_time', event.timestamp)
  Rack::Honeycomb.add_field(env, 'timer.event_time_delta_sec', event_time_delta)

  # verify schema - stub out
  $last_cache_time = 0

  get_schema_timer_start = (Time.now.to_f * 1000).to_i

  begin
    get_schema(users_dataset)
  rescue Schema_Lookup_Failure => e
    status 500
    return e.message.to_json
  end

  get_schema_timer_end = (Time.now.to_f * 1000).to_i
  get_schema_timer_ms  = (get_schema_timer_end - get_schema_timer_start)

  Rack::Honeycomb.add_field(env, 'timer.get_schema', get_schema_timer_ms)

  # Hand off to external service (aka write to local disk)
  write_event(event)
end

get '/x/alive' do
  Rack::Honeycomb.add_field(env, 'alive', true)
  'Sending stuff via Beeline'
end

get '/' do
  'This App is up and running!'
end
