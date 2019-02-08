# Dotnet Core

## Overview

A default dotnet core application built with `dotnet new webapi --no-https`. Added a piece of middleware to capture basic information about the http request in [HoneycombMiddleware.cs](./HoneycombMiddleware.cs) and added it to the middleware stack in [Startup.cs](./Startup.cs#L33).

## Build and Run

- Install the [.Net Core SDK](https://dotnet.microsoft.com/download). This application was built and tested with version 2.2
- Add your write key to the [HoneycombMiddleware.cs](./HoneycombMiddleware.cs#L18) and your dataset name to the [HoneycombMiddleware.cs](./HoneycombMiddleware.cs#L35)
- run `dotnet run` from this directory
- You can query the webapi by curling http://localhost:5000/api/values
