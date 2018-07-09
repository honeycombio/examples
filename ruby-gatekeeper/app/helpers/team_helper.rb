module TeamHelper
    class Team
      attr_accessor :id, :name, :write_key

      def initialize(args)
        @id = args.fetch(:id, 0)
        @name = args.fetch(:name, "")
        @write_key = args.fetch(:write_key, "")
      end
    end

end
