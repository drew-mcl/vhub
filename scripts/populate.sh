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

# Running population 


run_test POST http://localhost:8080/api/v1/regions '{"name":"amer"}' 
run_test POST http://localhost:8080/api/v1/regions/amer/environments '{"name":"qa"}'
run_test POST http://localhost:8080/api/v1/regions/amer/environments '{"name":"dev"}'

run_test POST http://localhost:8080/api/v1/regions '{"name":"newregion"}' 
run_test POST http://localhost:8080/api/v1/regions '{"name":"newregion2"}'
run_test POST http://localhost:8080/api/v1/regions/newregion/environments '{"name":"newenvironment"}'
run_test POST http://localhost:8080/api/v1/regions/newregion2/environments '{"name":"newenvironment"}'
run_test POST http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps '{"name":"app1"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps '{"name":"app2", "version":"1.0.0"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps '{"name":"app3", "version":"1.0.0", "route":"green"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion/environments/newenvironment/apps '{"name":"app4", "version":"1.0.0", "route":"blue"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion2/environments/newenvironment/apps '{"name":"app5", "version":"1.0.0", "route":"green"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion2/environments/newenvironment/apps '{"name":"app6", "version":"1.0.0", "route":"blue"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion2/environments/newenvironment/apps '{"name":"app7", "version":"1.0.0", "route":"green"}' application/json
run_test POST http://localhost:8080/api/v1/regions/newregion2/environments/newenvironment/apps '{"name":"app8", "version":"1.0.0", "route":"blue"}' application/json

echo "Population completed successfully."
exit 0
