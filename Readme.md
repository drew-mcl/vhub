# Version Hub
Version hub is a RESTful web service designed to provide an interface for managing versions of applications across regions, environments, and applications. Users can list, create, and retrieve details about regions, environments, and apps within those environments.

## Prerequisites
Go  - The programming language used.

## Endpoints

GET /regions - Lists all regions.
POST /regions/{regionName} - Creates a new region.
GET /regions/{regionName}/environments - Lists environments within a region.
POST /regions/{regionName}/{environmentName} - Creates an environment within a region.
GET /regions/{regionName}/environments/{environmentName}/apps - Lists all apps within a specified environment.
POST /regions/{regionName}/environments/{environmentName}/{appName} - Creates an app within a specified environment.
GET /regions/{regionName}/environments/{environmentName}/apps/{appName} - Retrieves details about a specific app within a specified environment.
