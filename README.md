# Employee Management API

This is a simple API built with the [Fiber](https://github.com/gofiber/fiber) web framework for managing employee records in a MongoDB database.

## Prerequisites

Before you start, ensure you have the following:

- [Go](https://golang.org/dl/) installed on your machine.
- [MongoDB](https://www.mongodb.com/try/download/community) installed and running.

## Getting Started

1. Clone the repository:

    ```bash
    git clone https://github.com/KarkiAnmol/HRMS-GoFiber-Mongo.git
    ```

2. Navigate to the project directory:

    ```bash
    cd HRMS-GoFiber-Mongo
    ```

3. Install dependencies:

    ```bash
    go mod download
    ```

4. Start the MongoDB instance on your local machine:

    ```bash
    # Replace "path/to/mongod" with the actual path to your MongoDB executable
    path/to/mongod
    ```

5. Build and Run the application:

    ```bash
    go build main.go
    go run main.go
    ```

The API should now be running at `http://localhost:3000`.

## Endpoints

### 1. Retrieve all employees

```http
GET /employee
```

### 2. Create a new employee

```http
POST /employee
```

### 3. Update employee details

```http
PUT /employee/:id
```

### 4. Delete an employee

```http
DELETE /employee/:id
```

## Usage

- Use your preferred API client (e.g., [Postman](https://www.postman.com/)) or tools like `curl` to interact with the API.
- Ensure that you have a MongoDB instance running before making requests.

## License

This project is licensed under the MIT License.
```

Feel free to customize this template according to your specific project details and requirements. Make sure to include relevant information about the API, endpoints, and usage instructions.
