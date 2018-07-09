module ErrorHelpers
    class Bad_Sample_Rate < StandardError
      attr_reader :message

      def initialize(message)
        @message = message
        super
      end
    end

    class Auth_Failure < StandardError
      attr_reader :message

      def initialize(message)
        @message = message
        super
      end
    end

    class Auth_Mishapen_Failure < StandardError
      attr_reader :message

      def initialize(message)
        @message = message
        super
      end
    end

    class Dataset_Lookup_Failure < StandardError
      attr_reader :message

      def initialize(message)
        @message = message
        super
      end
    end

    class Schema_Lookup_Failure < StandardError
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
