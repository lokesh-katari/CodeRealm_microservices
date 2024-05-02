<div align="center">
  <img src="https://github.com/lokesh-katari/CodeRealm_microservices/assets/111894942/ec05afd6-43b1-49c0-bb6c-ffb2e7e53d8c" alt="code" />
</div>

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

CodeRealm is a platform for coding evaluation and question evaluation that supports multiple programming languages. It enables users to compile and evaluate code snippets in languages such as Java, C++, Python, JavaScript, and Golang, ans also supports 10+ programming languages .It also has the frontend web editor builtin for HTML,CSS and JS . The application was originally built as a monolithic MERN stack application, but has been migrated to microservices architecture with Golang as the backend and Next.js as the frontend.
## Frontend Repo Link :
 - [Coderealm frontend](https://github.com/lokesh-katari/coderealm_frontend)

## Architecture
![code final asdf](https://github.com/lokesh-katari/CodeRealm_microservices/assets/111894942/793cb4f8-2f63-4219-a285-e4113e4af6a3)

## Previous repository links for the monolith architecture:
 - [Code Judge Online](https://github.com/lokesh-katari/Code-judge-Online)

## Features

- **Support for Multiple Languages**: CodeRealm supports 10+ programming languages for code compilation .
- **Microservices Architecture**: The application is built using microservices architecture, allowing for scalability and maintainability.
- **Queueing with Kafka**: Code submissions are queued using Kafka for efficient processing.
- **Database Integration**: PostgreSQL is used as the user database, while Redis is used for handling run requests.
- **Docker Compose and Kubernetes Deployment**: CodeRealm can be deployed using Docker Compose for local development or Kubernetes for production environments.

## Installation and Usage

### Prerequisites

- Docker and Docker Compose for local development.
- Kubernetes cluster for production deployment.


### Local Development

1. Clone the repository:

   ```bash
   git clone https://github.com/lokesh-katari/CodeRealm_microservices.git
   ```

2. Navigate to the project directory:

   ```bash
   cd CodeRealm_microservices
   ```

3. Start the application using Docker Compose:

   ```bash
   docker-compose up
   ```

4. Access the application at `http://localhost:3000`.

### Production Deployment
Before doing the below steps you need to configure the kafka strimzi cluster
  ```bash
    helm repo add strimzi https://strimzi.io/charts/

    kubectl create ns coderealm
    
    helm install strimzi-operator strimzi/strimzi-kafka-operator -n coderealm
  ```
make sure that you specify the namespace before deployment ,here :coderealm

1. Deploy the Kubernetes manifests in the `deployments/` directory:

   ```bash
   kubectl create ns coderealm 
   kubectl apply -f deployments/
   ```

2. Access the application using the provided Ingress configuration.
3. If you are using minikube cluster then use the command
  ```bash
  minikube ip
  ```
  for getting the ip from the minikube and you can access tit from the browser
4. also configure the port for the envoy to expose it to the local system using the command:
  ```bash
  kubectl port-forward service/code-frontend-service 8000:8000 -n coderealm

  ```

## Contributing

Contributions are welcome! Please feel free to fork the repository and submit pull requests to contribute new features, improvements, or bug fixes.

## License

This project is licensed under the [MIT License](LICENSE).

