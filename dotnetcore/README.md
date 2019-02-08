# Dotnet Core

## Overview

A default dotnet core application built with `dotnet new mvc --auth Individual`. Added a piece of middleware to capture basic information about the http request in [HoneycombMiddleware.cs](./HoneycombMiddleware.cs) and added it to the middleware stack in [Startup.cs](./Startup.cs#L51).

## Build and Run

- Install the [.Net Core SDK](https://dotnet.microsoft.com/download). This application was built and tested with version 2.2
- Add your write key to the [HoneycombMiddleware.cs](./HoneycombMiddleware.cs#L19) and your dataset name to the [HoneycombMiddleware.cs](./HoneycombMiddleware.cs#L40)
- run `dotnet run` from this directory
- You can access the sample site at http://localhost:5000
