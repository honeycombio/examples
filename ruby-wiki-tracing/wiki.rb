require "sinatra/base"
require "libhoney"

# Page represents the data (and some basic operations) on a wiki page.
#
# While the tracing instrumentation in this example is constrained to the
# handlers, we could just as easily propagate context down directly into this
# class if needed.
class Page
  attr_reader :filename, :title
  attr_accessor :body
  def initialize(title)
    @title = title
    @filename = "#{title}.txt"
  end

  def exist?
    File.exist?(@filename)
  end

  def save(body)
    File.write(@filename, body)
    true
  rescue => e
    false
  end
end

# Generate a new unique identifier for our spans and traces. This can be any
# unique string -- Zipkin uses hex-encoded base64 ints, as we do here; other
# folks may prefer to use their UUID library of choice.
def new_id
  rand(2**63).to_s(16)
end

# This middleware treats each HTTP request as a distinct "trace." Each trace
# begins with a top-level ("root") span indicating that the HTTP request has
# begun.
VALID_PATH = Regexp.new("^/(edit|save|view)/")
class RequestTracer
  def initialize(app)
    @app = app
  end
  def call(env)
    Thread.current[:request_id] = new_id
    m = env["REQUEST_PATH"].match(VALID_PATH)
    @app.with_span(m ? m[1] : env["REQUEST_PATH"]) do
      @app.call(env)
    end
  end
end

# This is our basic wiki webapp. It uses our RequestTracer middleware to track
# all HTTP requests with a root span, then defines a handful of handlers to
# handle the display / edit / saving of wiki pages on disk.
class App < Sinatra::Base
  use RequestTracer

  # Initialize our Honeycomb client once, and pull Honeycomb credentials from
  # an environment variable.
  configure do
    set :libhoney, Libhoney::Client.new(
      writekey: ENV["HONEYCOMB_WRITEKEY"],
      dataset: "ruby-wiki-tracing-example"
    )
  end

  # Redirect to a default wiki page.
  get "/" do
    redirect "/view/Index"
  end

  # Our "View" handler. Tries to load a page from disk and render it. Falls back
  # to the Edit handler if the page does not yet exist.
  get "/view/:title" do |title|
    @page = with_span("load_page", { title: title }) do
      load_page(title)
    end
    if @page.nil?
      return redirect "/edit/#{title}"
    end
    with_span("render_template", { template: "view" }) do
      erb(:view)
    end
  end

  # Our "Edit" handler. Tries to load a page from disk to seed the edit screen,
  # then renders a form to allow the user to define the content of the requested
  # wiki page.
  get "/edit/:title" do |title|
    @page = with_span("load_page", { title: title }) do
      load_page(title)
    end
    if @page.nil?
      @page = Page.new(title)
    end

    with_span("render_template", { template: "edit" }) do
      erb(:edit)
    end
  end

  # Our "Save" handler simply persists a page to disk.
  post "/save/:page" do |title|
    saved = with_span("File.write", { title: title, body_len: params["body"].size }) do
      p = Page.new(title)
      p.save(params["body"])
    end
    if saved
      return redirect "/view/#{title}"
    end
    "error"
  end

  # This wrapper takes a span name, some optional metadata, and a block; then
  # emits a "span" to Honeycomb as part of the trace begun in the RequestTracer
  # middleware.
  #
  # The special sauce in this method is the definition / resetting of thread
  # local variables in order to correctly propagate "parent" identifiers down
  # into the block.
  def with_span(name, metadata=nil)
    id = new_id
    start = Time.new
    # Field keys to trigger Honeycomb's tracing functionality on this dataset
    # defined at:
    # https://www.honeycomb.io/docs/working-with-data/tracing/send-trace-data/#manual-tracing
    data = {
      name: name,
      id: id,
      traceId: Thread.current[:request_id],
      serviceName: "wiki",
    }
    # Capture the calling scope's span ID, then restore it at the end of the
    # method.
    parent_id = Thread.current[:span_id]
    if parent_id
      data[:parentId] = parent_id
    end

    # Set the current span ID before invoking the provided block, then capture
    # the return value to return after emitting the Honeycomb event.
    Thread.current[:span_id] = id
    ret = yield

    data[:durationMs] = (Time.new - start)*1000
    if metadata
      data.merge!(metadata)
    end

    ev = settings.libhoney.event
    # NOTE: Don't forget to set the timestamp to `start` -- because spans are
    # emitted at the *end* of their execution, we want to be doubly sure that
    # our manually-emitted events are timestamped with the time that the work
    # (the span's actual execution) really begun.
    ev.timestamp = start
    ev.add(data)
    ev.send

    ret
  ensure
    Thread.current[:span_id] = parent_id
  end

 private

  # Helper method for returning a Page object for easy rendering
  def load_page(title)
    p = Page.new(title)
    return nil if !p.exist?
    with_span("File.read") do
      p.body = File.read(p.filename)
      p
    end
  end
end

# Let's go!
App.run!
