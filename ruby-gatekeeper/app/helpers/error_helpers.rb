module ErrorHelpers
  class BadSampleRate < StandardError
    attr_reader :message

    def initialize(message)
      @message = message
      super
    end
  end

  class AuthFailure < StandardError
    attr_reader :message

    def initialize(message)
      @message = message
      super
    end
  end

  class AuthMishapenFailure < StandardError
    attr_reader :message

    def initialize(message)
      @message = message
      super
    end
  end

  class DatasetLookupFailure < StandardError
    attr_reader :message

    def initialize(message)
      @message = message
      super
    end
  end

  class SchemaLookupFailure < StandardError
    attr_reader :message

    def initialize(message)
      @message = message
      super
    end
  end

  # extending class string
  class String
    def is_i?
      /\A[-+]?\d+\z/ === self
    end
  end
end
