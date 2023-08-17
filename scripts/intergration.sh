#!/bin/bash

echo "INTEGRATION TESTS"
echo "Please do not run this script in production"
echo ""

handle_failure() {
    echo "================================================================="
    echo "Integration test failed on $1 $2 status_code:$3"
    echo "================================================================="
    exit 1
}

run_test() {
  method=$1
  url=$2
  data=$3
  content_type=$4

  status_code=$(curl -o /dev/null -s -w "%{http_code}\n" -X $method -H "Content-Type: $content_type" -d "$data" $url)
  if [ $status_code -ge 200 ] && [ $status_code -lt 300 ]; then
    echo "Test passed for $url"
  else
    handle_failure $method $url $status_code
  fi
}

# Testing health check
echo "-----Running health check-----"
run_test GET http://localhost:8080/healthcheck

# Testing regions
echo "-----Running tests for regions-----"
run_test GET http://localhost:8080/api/v1/regions
run_test POST http://localhost:8080/api/v1/regions '{"name":"newregion"}' 
run_test GET http://localhost:8080/api/v1/regions/newregion
run_test DELETE http://localhost:8080/api/v1/regions/newregion

# Testing environments within a region
echo "-----Running tests for environments-----"
run_test POST http://localhost:8080/api/v1/regions '{"name":"newregion"}' 

run_test GET http://localhost:8080/api/v1/regions/newregion/environments
run_test POST http://localhost:8080/api/v1/regions/newregion/environments '{"name":"newenvironment"}'
run_test GET http://localhost:8080/api/v1/regions/newregion/environments/newenvironment
run_test DELETE http://localhost:8080/api/v1/regions/newregion

# Testing apps within an environment
echo "-----Running tests for apps-----"
run_test POST http://localhost:8080/api/v1/regions '{"name":"newregion"}' 
run_test POST http://localhost:8080/api/v1/regions/newregion/environments '{"name":"newenvironment"}'

run_test GET http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps
run_test POST http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps '{"name":"app1"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps '{"name":"app2", "version":"1.0.0"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps '{"name":"app3", "version":"1.0.0", "route":"green"}' application/json
run_test GET http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps/app1
run_test PUT http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps/app2 '{"version":"2.0.0"}' application/json
run_test GET http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps/app2
run_test PUT http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps/app3 '{"version":"2.0.0", "route":"blue"}' application/json
run_test DELETE http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps/app2
run_test DELETE http://localhost:8080/api/v1/regions/newregion

echo "Integration tests completed successfully."
exit 0
