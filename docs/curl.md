### Create a new region:
```bash
curl -X POST http://localhost:8080/api/v1/regions/myregion
```

### List regions:
```bash
curl -X GET http://localhost:8080/api/v1/regions
```

### Create an environment within a region:
```bash
curl -X POST http://localhost:8080/api/v1/regions/myregion/environments/myenvironment
```

### List environments within a region:
```bash
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments
```

### Create an app with default version (0.0.1):
```bash
curl -X POST http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp
```

### Create an app with specific version:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"version":"1.2.3"}' http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp
```

### List apps within a region and environment:
```bash
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps
```

### Get a specific app:
```bash
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp
```

### Get app version:
```bash
curl -X GET http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp/version
```

### Update app version:
```bash
curl -X PUT -H "Content-Type: application/json" -d '{"version":"2.0.0"}' http://localhost:8080/api/v1/regions/myregion/environments/myenvironment/apps/myapp/version
```