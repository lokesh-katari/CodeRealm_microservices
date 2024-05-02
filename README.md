# <h1 style="font-family: 'Arial'; color: #A5F3FC; text-align:center;  font-size: 60px" ><span style="color:#A5F3FC">Code<span style="color:rgb(236 254 255)">Realm</span>.</span></h1>

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

CodeRealm is a platform for coding evaluation and question evaluation that supports multiple programming languages. It enables users to compile and evaluate code snippets in languages such as Java, C++, Python, JavaScript, and Golang. The application was originally built as a monolithic MERN stack application, but has been migrated to microservices architecture with Golang as the backend and Next.js as the frontend.

## Features

- **Support for Multiple Languages**: CodeRealm supports 10+ programming languages for code compilation and evaluation.
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
   git clone https://github.com/your-username/CodeRealm.git
   ```

2. Navigate to the project directory:

   ```bash
   cd CodeRealm
   ```

3. Start the application using Docker Compose:

   ```bash
   docker-compose up
   ```

4. Access the application at `http://localhost:3000`.

### Production Deployment

1. Deploy the Kubernetes manifests in the `deployments/` directory:

   ```bash
   kubectl apply -f deployments/
   ```

2. Access the application using the provided Ingress configuration.

## Contributing

Contributions are welcome! Please feel free to fork the repository and submit pull requests to contribute new features, improvements, or bug fixes.

## License

This project is licensed under the [MIT License](LICENSE).
