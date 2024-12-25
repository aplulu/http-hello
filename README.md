# HTTP Hello Server

This project is an HTTP server designed for testing purposes on Kubernetes, KNative, and Cloud Run.

## Project Structure

- `main.go`: The main application code.
- `index.tmpl.html`: HTML template for the server response.
- `Dockerfile`: Dockerfile to build the Docker image.
- `go.mod` and `go.sum`: Go module files.

## Prerequisites

- Go 1.23 or later
- Docker

## Building the Project

To build the project, run the following commands:

```sh
docker build -t aplulu/http-hello .
```

## Running the Server

To run the server locally, use the following command:

```sh
docker run -p 8080:8080 aplulu/http-hello
```

The server will be available at `http://localhost:8080`.

## Environment Variables

The server can be configured using the following environment variables:

## Environment Variables

The server can be configured using the following environment variables:

| Variable         | Description                                | Default Value                |
|------------------|--------------------------------------------|------------------------------|
| `PORT`           | Port for the server to listen on           | `8080`                       |
| `JWKS_URL`       | URL for the JSON Web Key Set               |                              |
| `JWK_HEADER_NAME`| Header name for the JWK                    | `Authorization`              |
| `TITLE`          | Title for the HTML response                | `Congratulations`            |
| `SUCCESS_MESSAGE`| Success message for the HTML response      | `Successfully deployed.`     |
| `CUSTOM_DATA`    | Custom user data to display in the HTML response | `Issuer:iss,Subject:sub` |
| `TLS_INSECURE`   | Allow insecure TLS connections             | `false`                      |
| `TLS_ADDITIONAL_CERT_BUNDLE` | Additional certificate bundle for TLS | |

## Deployment

### Kubernetes

To deploy the server on Kubernetes, create a deployment and service YAML file:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: http-hello
spec:
  replicas: 1
  selector:
    matchLabels:
      app: http-hello
  template:
    metadata:
      labels:
        app: http-hello
    spec:
      containers:
      - name: http-hello
        image: aplulu/http-hello
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: http-hello
spec:
  selector:
    app: http-hello
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

Apply the YAML file:

```sh
kubectl apply -f deployment.yaml
```

### KNative

To deploy the server on KNative, create a service YAML file:

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: http-hello
spec:
  template:
    spec:
      containers:
      - image: aplulu/http-hello
        ports:
        - containerPort: 8080
```

Apply the YAML file:

```sh
kubectl apply -f service.yaml
```

### Cloud Run

To deploy the server on Cloud Run, use the following command:

```sh
gcloud run deploy http-hello --image aplulu/http-hello --platform managed
```

## License

[MIT License](LICENSE)
