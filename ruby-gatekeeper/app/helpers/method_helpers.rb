module MethodHelpers
  def all_letters_or_digits(str)
    str[/[a-zA-Z0-9]+/] == str
  end

  # returns true if the string contents map to a positive, non float, integer
  def positive_integer?(sample_rate)
    unless /\A\d+\z/ =~ sample_rate
      # sample rate string maps to a negative number or non-integer false
      raise BadSampleRate.new("error": 'bad sample rate provided')
    end

    # Is all good ...continue
    sample_rate
  end

  def validate_write_key(users_write_key)
    unless all_letters_or_digits(users_write_key)
      raise AuthMishapenFailure.new("error": 'writekey malformed - expect only letters and numbers')
    end

    matching_team = $known_teams.select { |team| team.write_key == users_write_key }
    if matching_team.empty?
      raise AuthFailure.new("error": 'writekey does not match valid credentials')
    end

    matching_team[0]
  end

  def resolve_dataset(given_dataset)
    matching_dataset = $known_datasets.select { |dataset| dataset.name == given_dataset }
    if matching_dataset.empty?
      raise DatasetLookupFailure.new("error": 'failed to resolve dataset')
    end

    matching_dataset[0]
  end

  def grab_partition(given_dataset)
    available_partitions = given_dataset.partition_list

    raise 'no partitions found' if available_partitions.length <= 0

    available_partitions.sample
  end

  def get_schema(_given_dataset)
    hit_cache = true

    $last_cache_time = Time.now if $last_cache_time == 0

    if (Time.now - $last_cache_time) > CACHE_TIMEOUT
      # we fall through the cache every 10 seconds
      hit_cache = false
      # pretend to hit a slow database that takes 30-50ms
      sleep(rand(30...51))
      $last_cache_time = Time.now
    end
    Rack::Honeycomb.add_field(env, 'hit_schema_cache', hit_cache)

    # let's just fail sometimes to pretend
    return unless rand(60).zero?

    raise SchemaLookupFailure.new("error": 'failed to resolve schema')
  end

  def write_event(event)
    File.open("/tmp/api#{event.chosen_partition}.log", 'a') { |f| f.puts event.to_json }
  end
end
