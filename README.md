
# Sphinx-Tribes Backend

![Tribes](https://github.com/stakwork/sphinx-tribes/raw/master/img/sphinx-tribes.png)

Sphinx-Tribes is a decentralized message broker for public groups within the Sphinx ecosystem. This README covers the setup and configuration of the Sphinx-Tribes backend.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Setup](#setup)
  - [Cloning the Repository](#cloning-the-repository)
  - [Building the Docker Image](#building-the-docker-image)
  - [Environment Configuration](#environment-configuration)
  - [Database Setup](#database-setup)
  - [Running the Backend](#running-the-backend)
- [Optional Features](#optional-features)
  - [Redis for Caching](#redis-for-caching)
  - [Relay Integration](#relay-integration)
  - [Meme Image Upload](#meme-image-upload)
  - [SuperAdmin Dashboard Access](#superadmin-dashboard-access)
  - [Stakwork YouTube Integration](#stakwork-youtube-integration)
- [Testing and Mocking](#testing-and-mocking)
  - [Unit Testing](#unit-testing)
  - [Mocking Interfaces](#mocking-interfaces)
- [API Data Validations](#api-data-validations)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

- Docker
- Go language environment
- PostgreSQL database
- Redis instance (optional)
- Relay server access (optional)

## Setup

### Cloning the Repository

Clone the Sphinx-Tribes repository:

```sh
git clone https://github.com/Vayras/sphinx-tribes.git
```

### Building the Docker Image

Navigate to the cloned directory and build the Docker image:

```sh
docker build --no-cache -t sphinx-tribes .
docker tag sphinx-tribes sphinxlightning/sphinx-tribes:x
docker push sphinxlightning/sphinx-tribes:x
```

### Environment Configuration

Create a `.env` file in the project root with the required environment variables.

### Database Setup

Set up a PostgreSQL database and execute the provided SQL scripts to create necessary tables.

### Running the Backend

Build and run the Golang backend:

```sh
go build .
./sphinx-tribes
```

## Optional Features

### Redis for Caching

Configure Redis by adding the `REDIS_URL` or other relevant variables to your `.env` file.

### Relay Integration

For invoice creation and keysend payment, add `RELAY_URL` and `RELAY_AUTH_KEY`.

### Meme Image Upload

Requires a running Relay. Enable it with `MEME_URL`.

### SuperAdmin Dashboard Access

Add public keys to `SUPER_ADMINS` in your `.env` file.

### Stakwork YouTube Integration

Add `STAKWORK_KEY` for YouTube video downloads.

## Testing and Mocking

### Unit Testing

Run unit tests with coverage:

```sh
go test ./... -tags mock -race -v -coverprofile=coverage.out
go tool cover -html="coverage.out"
```

### Mocking Interfaces

Use [mockery](https://vektra.github.io/mockery/latest/) for interface mocking.

## API Data Validations

Implement validations using `gopkg.in/go-playground/validator.v9`.

## Contributing

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the [LICENSE NAME] - see the [LICENSE.md](LICENSE.md) file for details.
