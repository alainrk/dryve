# Dryve

A simple file storage service written in Go.

## Run

### Dependencies

- Docker `>= 20.10.12`
- Docker Compose `>= 3.8`
- Make

### Running on Docker

Run the whole stack (server and postgres database)

```sh
make start
```

If you want to change the configuration:
- `config.json` for local configuration (needed to run automigration)
- `config-docker.json` for docker configuration (picked up by docker compose)

## Development

- Go `>= 1.19`
- Docker `>= 20.10.12`
- Docker Compose `>= 3.8`
- Make

### Running local server for development

```sh
# Run the postgres database
make start-db

# Automigrate (basically creates database and tables)
make automigration

# Run live-reloading server
make dev
```

### Run tests

```sh
make test
```

## Implementation

This server provides APIs to handle file upload, download, deletion and metadata retrieval.
Some API endpoints are protected using a basic rate limiter to prevent abuse (on the single server instance).
The server architecture follows an exagonal architecture structure to be modular and flexible.

```sh
.
├── cmd
│   ├── automigrate   # Entrypoint for automigration script
│   └── server        # Entrypoint for API server
├── internal
│   ├── app           # API endpoints entrypoints
│   ├── config        # Configuration management
│   ├── datastruct    # Models
│   ├── dto           # Model structures for request/response
│   ├── repository    # Database layer management
│   └── service       # Business logic controllers
```

API Endpoints:

- `GET /files/{id}`: Retrieves the file metadata for the file with the given ID.
- `GET /files/range/{from}/{to}`: Retrieves the file metadata for all files within the specified date range.
- `POST /files`: Uploads a file to the server.
- `GET /files/{id}/download`: Downloads the file with the given ID.
- `DELETE /files/{id}`: Deletes the file with the given ID.
- `DELETE /files/range/{from}/{to}`: Deletes all files within the specified date range.

```sh
# Upload a file
curl -X POST -F "file=@{ABSOLUTE_PATH}" http://localhost:8666/files

# Get file metadata
curl http://localhost:8666/files/44fdac3e-5384-4eb3-94f4-e7a0fd0cee15

# Download a file
curl http://localhost:8666/files/2b0f8f45-7ffc-479d-8189-794bf02e0fa7/download

# Get files metadata in a date range
curl http://localhost:8666/files/range/2021-09-10/2024-04-30

# Delete a file
curl -X DELETE http://localhost:8666/files/range/2021-09-10/2024-04-30

# Delete files in a date range
curl -X DELETE http://localhost:8666/files/44fdac3e-5384-4eb3-94f4-e7a0fd0cee15
```
