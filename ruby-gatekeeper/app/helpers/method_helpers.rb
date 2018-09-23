module MethodHelpers
  def all_letters_or_digits(str)
    str[/[a-zA-Z0-9]+/] == str
  end

  # returns true if the string contents map to a positive, non float, integer
  def positive_integer?(sample_rate)
    if !/\A\d+\z/.match(sample_rate)
      # sample rate string maps to a negative number or non-integer
      # false
      raise Bad_Sample_Rate.new("error": 'bad sample rate provided')
    else
      # Is all good ..continue
      sample_rate
    end
  end

  def validate_write_key(users_write_key)
    unless all_letters_or_digits(users_write_key)
      raise Auth_Mishapen_Failure.new("error": 'writekey malformed - expect only letters and numbers')
    end

    matching_team = $known_teams.select { |team| team.write_key == users_write_key }
    if matching_team.empty?
      raise Auth_Failure.new("error": 'writekey does not match valid credentials')
    end

    matching_team[0]
  end

  def resolve_dataset(given_dataset)
    matching_dataset = $known_datasets.select { |dataset| dataset.name == given_dataset }
    if matching_dataset.empty?
      raise Dataset_Lookup_Failure.new("error": 'failed to resolve dataset')
    end

    matching_dataset[0]
  end

  def grab_partition(given_dataset)
    available_partitions = given_dataset.partition_list

    if available_partitions.length <= 0
      raise 'no partitions found'
    else
      return available_partitions.sample
    end
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
    if rand(60) == 0
      raise Schema_Lookup_Failure.new("error": 'failed to resolve schema')
    end
  end

  def write_event(event)
    File.open("/tmp/api#{event.chosen_partition}.log", 'a') { |f| f.puts event.to_json }
  end
end
