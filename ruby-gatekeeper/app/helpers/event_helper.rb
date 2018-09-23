module EventHelper
  class Event
    attr_accessor :write_key, :timestamp, :sample_rate, :data, :chosen_partition

    def initialize(args)
      @write_key        = args.fetch(:write_key, '')
      @timestamp        = args.fetch(:timestamp, '')
      @sample_rate      = args.fetch(:sample_rate, 0)
      @data             = args.fetch(:data, '')
      @chosen_partition = args.fetch(:chosen_partition, 0)
    end
  end
end
