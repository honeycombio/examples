const HONEYCOMB_DATASET = "tracing-example";

const beeline = require("honeycomb-beeline")({
  writeKey: process.env.HONEYCOMB_API_KEY,
  dataset: HONEYCOMB_DATASET,
  serviceName: "wall",
});

const https = require('https');
const express = require('express');
const bodyParser = require('body-parser');
const rp = require('request-promise-native');
const app = express();

const contents = ['first post'];
const hashtagRegexp = /#([a-z0-9]+)/gi;
const hashtagSearch = `<a href="https://twitter.com/hashtag/$1">#$1</a>`;

app.use(bodyParser.urlencoded({ extended: false, type: "*/*" }));

// = MIDDLEWARE ====================================================
// Wraps HTTP handlers to output evidence of HTTP calls + trace IDs
// to STDOUT for debugging.
// =================================================================
app.use((req, res, next) => {
  let traceContext = beeline.marshalTraceContext(beeline.getTraceContext());
  let { traceId } = beeline.unmarshalTraceContext(traceContext);
  console.log("Handling request with:", JSON.stringify({
    method: req.method,
    path: req.path,
    traceId: traceId,
  }));

  // Add some way to identify these requests as coming from you!
  // customContext.add ensures that this field will be populated on
  // *all* spans in the trace, not just the currently active one.
  beeline.customContext.add("username", "YOUR_USERNAME_HERE");

  next();
});
app.get('/favicon.ico', (req, res) => { res.status(404); });

// = HANDLER =======================================================
// Returns the current contents of our "wall".
// =================================================================
app.get('/', (req, res) => {
  res.send(`${ contents.join("<br />\n") }
    <br /><br /><a href="/message">+ New Post</a>`);
});

// = HANDLER =======================================================
// Returns a simple HTML form for posting a new message to our wall.
// =================================================================
app.get('/message', (req, res) => {
  res.send(`<form method="POST" action="/">
		<input type="text" autofocus name="message" /><input type="submit" />
	</form>`);
});

// = HANDLER =======================================================
// Processes a string from the client and saves the message contents.
// =================================================================
app.post('/', async(req, res) => {
  if (typeof(req.body.message) !== 'string') {
    beeline.customContext.add('error', 'non-string body');
    res.status(500).send("not a string body");
    return;
  }
  let body = req.body.message.trim();

  let analysisPromise = analyze(body);

  body = await twitterize(body);

  // Let's persist our wall contents! POST each message to a
  // third-party service (in this case, a Lambda function).
  await persist(body);

  let sentiment = await analysisPromise;
  beeline.addContext({ sentiment });
  if (sentiment >= 0.2) {
    body = `<b>${ body }</b>`;
  } else if (sentiment <= -0.2) {
    body = `<i>${ body }</i>`;
  }

  contents.push(body);

  res.redirect("/");
});

// = HELPER ========================================================
// Identifies hashtags and Twitter handle-like strings. Replaces
// hashtags with links to a Twitter search for the found hashtag and
// replaces handle-like strings with links to the Twitter profile
// *if* a valid profile is found.
// =================================================================
const twitterize = async(content) => {
  let newContent = content.replace(hashtagRegexp, hashtagSearch);

  let matches = newContent.match(/@([a-z0-9_]+)/g);
  let promiseArr = (matches || []).map((handle) => {
    const profile = `https://twitter.com/${ handle.substr(1) }`;
    return beeline.startAsyncSpan({
      name: "check_twitter",
      "app.twitter.handle": handle,
    }, span => {
      return rp({ uri: profile, resolveWithFullResponse: true }).then((resp) => {
        newContent = newContent.replace(handle, `<a href="${ profile }">${ handle }</a>`);
        beeline.addContext({ "app.twitter.response_status": resp.statusCode });
        beeline.finishSpan(span);
      }).catch((err) => {
       beeline.addContext({ "app.twitter.response_status": err.statusCode });
       beeline.finishSpan(span);
      });
    });
  });

  await Promise.all(promiseArr);
  return newContent;
};

// = HELPER ========================================================
// Calls out to a second service of ours (which may or may not be
// live) to perform some further analysis on the post content.
// =================================================================
const analyze = async(content) => {
  let options = {
    uri: "http://localhost:8088",
    method: "POST",
    headers: { "Content-Type": "text/plain" },
    body: content,
  };
  return rp(options).then((v) => v).catch((err) => 0.0);
};

// = HELPER ========================================================
// Calls out to a third service over which we *don't* have control,
// in order to persist the contents of a single message.
// =================================================================
const persist = async(content) => {
  let options = {
    uri: "https://p3m11fv104.execute-api.us-east-1.amazonaws.com/dev/",
    method: "POST",
    headers: { "Content-Type": "text/plain" },
    body: content,
  };
  return rp(options).then(() => {}).catch(() => {});
};

app.listen(8080, () => console.log(`'wall' service listening on port 8080!`));
