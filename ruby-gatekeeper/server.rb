# frozen_string_literal: true

require "rails"
require "action_controller/railtie"
require "active_model/railtie"
require "honeycomb-beeline"

CACHE_TIMEOUT = 10
$last_cache_time = 0

Honeycomb.configure do |config|
  config.debug = true
end

class Gatekeeper < Rails::Application
  routes.append do
    get "/x/alive", to: "main#alive"
    post "/1/events/:dataset_name", to: "main#events"
  end

  config.api_only = true
end

class Team
  include ActiveModel::Model

  attr_accessor :id, :name, :write_key

  def self.all
    @teams ||= [
      Team.new(id: 1, name: "RPO", write_key: "abcd123EFGH"),
      Team.new(id: 2, name: "B&W", write_key: "ijkl456MNOP"),
      Team.new(id: 3, name: "Third", write_key: "qrst789UVWX"),
    ]
  end
end

class Dataset
  include ActiveModel::Model

  attr_accessor :id, :name, :partition_list

  def partition
    Honeycomb.start_span(name: "claim_paritition") do
      partition_list.sample.tap do |partition|
        Honeycomb.add_field "partition", partition
      end
    end
  end

  def schema
    Honeycomb.start_span(name: "schema_lookup") do
      hit_cache = true

      $last_cache_time = Time.now if $last_cache_time == 0

      if (Time.now - $last_cache_time) > CACHE_TIMEOUT
        # we fall through the cache every 10 seconds
        hit_cache = false

        # pretend to hit a slow database that takes 30-50ms
        sleep(rand(30...51).fdiv(1000))
        $last_cache_time = Time.now
      end

      Honeycomb.add_field "hit_schema_cache", hit_cache

      # let's just fail sometimes to pretend
      return unless rand(60).zero?

      raise StandardError.new("failed to resolve schema")
    end
  end

  def self.all
    @datasets ||= [
      Dataset.new(id: 1, name: "wade", partition_list: [1, 2, 3]),
      Dataset.new(id: 2, name: "james", partition_list: [1, 2, 4]),
      Dataset.new(id: 3, name: "helen", partition_list: [1, 3, 4]),
      Dataset.new(id: 4, name: "peter", partition_list: [1, 2, 4]),
      Dataset.new(id: 5, name: "valentine", partition_list: [1, 2, 4]),
      Dataset.new(id: 6, name: "andrea", partition_list: [2, 3, 4]),
    ]
  end
end

class Event
  include ActiveModel::Model

  attr_accessor :sample_rate, :timestamp, :write_key, :data, :dataset

  validates :sample_rate, numericality: { only_integer: true, greater_than_or_equal_to: 0 }
  validates :write_key, format: { with: /\A[a-zA-Z0-9]+\z/ }
  validates :dataset, presence: true
  validates :timestamp, presence: true
  validate :valid_json_data
  validate :valid_team

  def valid_json_data
    Honeycomb.start_span(name: "json_parse") do
      begin
        JSON.parse data
      rescue StandardError
        errors.add(:data, "must be valid JSON")
      end
    end
  end

  def valid_team
    Honeycomb.start_span(name: "team_lookup") do
      team = Team.all.find do |team|
        team.write_key == write_key
      end

      if team
        Honeycomb.add_field_to_trace "team_id", team.id
      else
        errors.add(:write_key, "doesn't match valid credentials")
      end
    end
  end

  def write
    Honeycomb.start_span(name: "write_event") do |span|
      begin
        schema = dataset.schema
        Tempfile.create(dataset.partition.to_s) do |f|
          f << data.to_json
        end
        true
      rescue StandardError => e
        puts e.backtrace
        errors.add(:data, "failed to write")
        false
      end
    end
  end
end

class MainController < ActionController::API
  def events
    write_key = request.headers["X-Honeycomb-Team"]
    timestamp = request.headers["X-Honeycomb-Event-Time"]
    sample_rate = request.headers["X-Honeycomb-Samplerate"]

    Honeycomb.add_field "sample_rate", sample_rate
    Honeycomb.add_field "write_key", write_key
    Honeycomb.add_field "incoming_timestamp", timestamp

    dataset_name = params[:dataset_name]

    sample_rate ||= "1"
    timestamp ||= Time.now.utc

    dataset = Honeycomb.start_span(name: "dataset_lookup") do
      Dataset.all.find do |dataset|
        dataset.name == dataset_name
      end
    end

    event = Event.new(
      write_key: write_key,
      timestamp: timestamp,
      sample_rate: sample_rate,
      data: request.body.read,
      dataset: dataset,
    )

    if event.valid? && event.write
      render status: :accepted
    else
      render json: event.errors.full_messages, status: :bad_request
    end
  end

  def alive
    Honeycomb.add_field "alive", true

    render status: :ok
  end
end
