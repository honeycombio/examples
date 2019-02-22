using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;

namespace dotnet_core_webapi_sample
{
  public class HoneycombMiddleware
  {
    private readonly RequestDelegate _next;
    private readonly HttpClient _client;
    public HoneycombMiddleware(RequestDelegate next)
    {
      _next = next;
      _client = new HttpClient();
      _client.DefaultRequestHeaders.Add("X-Honeycomb-Team","<HONEYCOMB_API_KEY>");
    }

    public async Task InvokeAsync(HttpContext context)
    {
      var now = DateTime.UtcNow;
      var stopwatch = new Stopwatch();
      var fields = new Dictionary<string, object> ();
      fields.Add("request.path", context.Request.Path.Value);
      fields.Add("request.method", context.Request.Method);
      fields.Add("request.content_length", context.Request.ContentLength);
      fields.Add("request.host", context.Request.Host.ToString());
      stopwatch.Start();
      await _next.Invoke(context);
      stopwatch.Stop();
      fields.Add("duration_ms", stopwatch.ElapsedMilliseconds);
      fields.Add("response.status_code", context.Response.StatusCode);
      fields.Add("response.content_length", context.Response.ContentLength);
      var dataset = "<DATASET_NAME>";
      // TODO: Think about sending this out of band from the web request
      await _client.PostAsJsonAsync($"https://api.honeycomb.io/1/events/{dataset}", fields);
    }
  }
}
