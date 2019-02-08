using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Http;

namespace dotnetcore
{
  public class HoneycombMiddleware
  {
    private readonly RequestDelegate _next;
    private readonly HttpClient _client;
    public HoneycombMiddleware(RequestDelegate next)
    {
      _next = next;
      _client = new HttpClient();
      _client.DefaultRequestHeaders.Add("X-Honeycomb-Team","<WRITEKEY>");
    }

    public async Task InvokeAsync(HttpContext context, UserManager<IdentityUser> userManager)
    {
      var now = DateTime.UtcNow;
      var stopwatch = new Stopwatch();
      var fields = new Dictionary<string, object> ();
      fields.Add("request.path", context.Request.Path.Value);
      fields.Add("request.method", context.Request.Method);
      fields.Add("request.content_length", context.Request.ContentLength);
      stopwatch.Start();
      await _next.Invoke(context);
      stopwatch.Stop();
      var user = await userManager.GetUserAsync(context.User);
      if (user != null) {
        fields.Add("user.id", user.Id);
      }
      fields.Add("duration_ms", stopwatch.ElapsedMilliseconds);
      fields.Add("response.http_status", context.Response.StatusCode);
      fields.Add("response.content_length", context.Response.ContentLength);
      var dataset = "<DATASET_NAME>";
      // TODO: Think about sending this out of band from the web request
      await _client.PostAsJsonAsync($"https://api.honeycomb.io/1/events/{dataset}", fields);
    }
  }
}
