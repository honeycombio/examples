using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;

namespace dotnet_core_webapi_sample
{
  public class HoneycombContext
  {
    public Dictionary<string, object> Fields { get; }

    public HoneycombContext()
    {
      Fields = new Dictionary<string, object>();
    }
  }

  public class HoneycombMiddleware
  {
    private readonly RequestDelegate _next;
    private readonly HttpClient _client;
    public HoneycombMiddleware(RequestDelegate next)
    {
      _next = next;
      _client = new HttpClient();
      _client.DefaultRequestHeaders.Add("X-Honeycomb-Team","6b7919e9cb69152d49b8bd3aa1fffbc4");
    }

    public async Task InvokeAsync(HttpContext context, HoneycombContext honeycombContext)
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
      foreach (var field in honeycombContext.Fields) {
        fields.Add(field.Key, field.Value);
      }
      var dataset = "dotnet-core-webapi";
      // TODO: Think about sending this out of band from the web request
      await _client.PostAsJsonAsync($"https://api.honeycomb.io/1/events/{dataset}", fields);
    }
  }
}
