module DatasetHelper
    class Dataset
      attr_accessor :id, :name, :partition_list

      def initialize(args)
        @id = args.fetch(:id, 0)
        @name = args.fetch(:name, "")
        @partition_list = args.fetch(:partition_list, [])
      end
    end

end
