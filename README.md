# Greenlight API

## Overview

Greenlight is a simple movie database API built with Go. It allows users to create, read, update, and delete movies in the database.

## Features

- **Health Check**: Verify the health of the API.
- **Create Movie**: Add a new movie to the database.
- **List Movies**: Retrieve all movies from the database.
- **Show Movie**: Get details of a specific movie by ID.
- **Update Movie**: Update details of an existing movie.
- **Delete Movie**: Remove a movie from the database.

## Endpoints

- **GET /v1/healthcheck**: Check the API health status.
- **POST /v1/movies**: Create a new movie.
- **GET /v1/movies**: List all movies.
- **GET /v1/movies/{id}**: Get a specific movie by ID.
- **PATCH /v1/movies/{id}**: Update a specific movie by ID.
- **DELETE /v1/movies/{id}**: Delete a specific movie by ID.

## Getting Started

### Prerequisites

- Go 1.18+
- PostgreSQL

### Setup

1. **Clone the repository:**

   ```sh
   git clone https://github.com/vivekup3424/greenlight.git
   cd greenlight
   ```

2. **Set up PostgreSQL:**

   Ensure PostgreSQL is installed and running. Create a database and user for the application.

   ```sql
   CREATE DATABASE greenlight;
   CREATE USER greenlight WITH PASSWORD 'pa55word';
   GRANT ALL PRIVILEGES ON DATABASE greenlight TO greenlight;
   ```

3. **Set up environment variables:**

   Create a `.env` file in the project root and add your PostgreSQL DSN:

   ```env
   DB_DSN=postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable
   ```

4. **Run the application:**

   ```sh
   go run ./cmd/api
   ```

   The server will start on `http://localhost:4000`.

## Project Structure
