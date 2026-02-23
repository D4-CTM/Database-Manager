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
