# Database Manager Tool

A web-based database management tool with a Go backend and a Vue frontend, designed for Oracle databases.

## Features

-   **Connection Management:** Save, edit, and delete connection credentials for multiple Oracle databases.
-   **Database Explorer:** Browse database objects, including tables, views, procedures, functions, packages, sequences, triggers, and indexes.
-   **SQL Editor:** Execute arbitrary SQL queries and view the results in a table format.
-   **Metadata Viewer:** View detailed information about tables, columns, constraints, and function arguments.
-   **DDL Viewer:** View the DDL for a given database object.

## Architecture

The application is divided into two main components: a backend server and a frontend client.

### Backend

The backend is a Go application built with the Gin framework. It serves a RESTful API for the frontend to consume. Its responsibilities include:

-   Managing database connections.
-   Executing SQL queries against the user's database.
-   Querying database metadata.

### Frontend

The frontend is a single-page application (SPA) built with Vue.js. It provides the user interface for interacting with the application. Its responsibilities include:

-   Displaying a list of saved connections.
-   Providing a way to create, edit, and delete connections.
-   Displaying the database explorer.
-   Providing a SQL editor for executing queries.
-   Displaying query results and metadata.

## Tech Stack

-   **Backend:** Go, Gin, Godror
-   **Frontend:** Vue.js, TypeScript, Bootstrap
-   **Database:** Oracle
-   **Containerization:** Docker

## Project Structure

-   `backend/`: Contains the Go backend application.
-   `frontend/`: Contains the Vue.js frontend application.
-   `init-scripts/`: Contains SQL scripts for initial database setup.
-   `docker-compose.yml`: Defines the services, networks, and volumes for the application.

## Getting Started

### Prerequisites

-   Docker
-   Docker Compose

### Running the Application

1.  Clone the repository:
    ```bash
    git clone https://github.com/johndoe/Database-Manager.git
    ```
2.  Navigate to the project directory:
    ```bash
    cd Database-Manager
    ```
3.  Start the application using Docker Compose:
    ```bash
    docker-compose up -d
    ```
4.  The frontend will be available at `http://localhost:4173`.

## Services

-   **`goapp`**: The Go backend service.
    -   Port: `5461`
-   **`frontend`**: The Vue.js frontend service.
    -   Port: `4173`
-   **`oracle`**: The Oracle database service.
    -   Port: `1521`
    -   User: `appuser`
    -   Password: `appuser_pass`

## Database Connection

To connect to the Oracle database from an external tool, use the following credentials:

-   **Host**: `localhost`
-   **Port**: `1521`
-   **Service Name**: `FREEPDB1`
-   **User**: `appuser`
-   **Password**: `appuser_pass`

## How the Go Backend Connects to Oracle

This section explains how the Go backend connects to an Oracle database using the `godror` driver.

### Prerequisites

The connection relies on a few key components:

1.  **Oracle Instant Client:** The `godror` driver is a CGo-based driver, which means it requires the Oracle Instant Client libraries to be available in the environment.
2.  **`CGO_ENABLED=1`:** This environment variable must be set to enable CGo for the Go compiler to build the driver correctly.
3.  **`godror` driver:** The Go application uses the `github.com/godror/godror` package to interface with the Oracle database.

### Docker Configuration

The `backend/Dockerfile` handles the setup of the environment for running the application:

1.  **Installation of Dependencies:** The Dockerfile starts by installing `libaio1`, which is a prerequisite for the Oracle Instant Client.
2.  **Oracle Instant Client Setup:** It downloads and unzips the Oracle Instant Client into the `/opt/oracle/instantclient_23_6` directory within the Docker image.
3.  **Environment Variables:** It sets the following environment variables:
    *   `LD_LIBRARY_PATH="/opt/oracle/instantclient_23_6"`: This tells the system's dynamic linker where to find the Oracle Instant Client's shared libraries (`.so` files).
    *   `CGO_ENABLED=1`: This enables CGo, allowing the Go program to call the C code in the Oracle Instant Client.

### Go Implementation

The Go code for handling database connections is located in `backend/Service/Credentials.go`.

The connection is established using the standard `database/sql` package along with the `godror` driver. The driver is imported with a blank identifier (`_ "github.com/godror/godror"`) to register it with `database/sql` without explicitly using any of its exported functions.

The connection string is constructed dynamically using `fmt.Sprintf` in the following format:

```go
db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`, c.User, c.Password, c.Server, c.Port, c.Database))
```

The parameters are:

*   `user`: The database username.
*   `password`: The user's password.
*   `connectString`: An Easy Connect string composed of the server's hostname/IP, the port, and the database service name (`<server>:<port>/<database>`).

This `sql.Open` call creates a `*sql.DB` object that can be used to query the database.
