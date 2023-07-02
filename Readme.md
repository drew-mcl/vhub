# vhub

## Overview

vhub is a simple HTTP server written in Go that tracks version information for multiple applications across different environments. It provides an API to set and retrieve version information.

The server stores application version data in memory, and it initializes this data from a JSON configuration file upon startup. The version data can be updated through API calls.

The server also handles graceful shutdown upon receiving a termination signal, allowing any ongoing operations to complete before the server stops.

## Usage

You can start the server with the following command:

```bash
go run main.go -config=<path-to-config> -port=<port>
```

Here, <path-to-config> should be replaced with the path to your JSON configuration file, and <port> should be replaced with the port you want the server to listen on. If these flags are not provided, the server uses config.json as the default configuration file and 8080 as the default port.

## API

The server provides the following API endpoints:

* GET /api/version?env=<environment>&app=<app>: Returns the version number for the specified app in the specified environment.

* POST /api/version: Sets the version for a specific application in a specific environment. This endpoint expects form data with fields env, app, and version.

In addition to the API endpoints, the server provides a simple landing page at the root (GET /) endpoint.

## Configuration File Format

The configuration file should be in JSON format with the following structure:

```json
{
  "environments": [
    {
      "name": "<environment-name>",
      "versions": {
        "<app-name>": <version-number>,
        "<app-name>": <version-number>,
        ...
      }
    },
    ...
  ]
}
```
## Testing

The main_test.go file contains a suite of tests for the API handlers. You can run these tests using the go test command:

```bash
go test -v ./...
```
## Future Enhancements

The current implementation keeps the version data in memory, meaning that all data will be lost if the server stops. A future enhancement could be to store the data in a persistent database.