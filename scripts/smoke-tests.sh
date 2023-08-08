#!/bin/bash

echo "SMOKE TESTS"
echo "please do not run this script in production"
echo "" 

echo "-----Running tests for regions-----"
curl -X GET http://localhost:8080/api/v1/regions
curl -X POST http://localhost:8080/api/v1/regions/myregion
curl -X GET http://localhost:8080/api/v1/regions/myregion

echo "" 
echo "-----Running tests for environments-----"
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments
curl -X POST http://localhost:8080/api/v1/regions/myregion/environments/myenvironment
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment

echo "" 
echo "-----Running tests for apps-----"
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps

echo "--> Creating app with default version"
curl -X POST http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp 
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp

echo "--> Creating app with version"
curl -X POST -H "Content-Type: application/json" -d '{"version":"1.2.3"}' http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp2
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp2

echo "" 
echo "-----Running tests for versions-----"
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp/version

echo "--> Updating version"
curl -X PUT -H "Content-Type: application/json" -d '{"version":"2.0.0"}' http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp/version
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp/version
