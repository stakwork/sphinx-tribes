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
- [Backend API Data Validations](#backend-api-data-validations)
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

```sh
    RDS_HOSTNAME =
    RDS_PORT =
    RDS_DB_NAME =
    RDS_USERNAME =
    RDS_PASSWORD =
```

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
    // you may need to install cover with this command first
    go get golang.org/x/tools/cmd/cover
    // run test
    RELAY_AUTH_KEY=TEST go test ./... -tags mock -race -v -coverprofile=coverage.out && ./cover-check.sh coverage.out <min coverage amount>
    // To get code coverage in html format do the following after running the code above
    go tool cover -html="coverage.out"
```

### Mocking Interfaces

Use [mockery](https://vektra.github.io/mockery/latest/) for interface mocking.

#### Installing mockery

There are multiple options to install mockery. Use any one of the following to download.

##### Download the mockery binary

Use the release page link [mockery releases](https://github.com/vektra/mockery/releases/tag/v2.39.1) to download the artifact for your respective device.

##### Using go install

If you have go already installed on your device you can use the go install command to download mockery.

```sh
go install github.com/vektra/mockery/v2@v2.39.1
```

##### Using homebrew

If you are on mac you can use homebrew to download mockery

```zsh
brew install mockery
brew upgrade mockery
```

#### When adding a new function to the interface which is already mocked follow the below steps

1. Update the corresponding interface with the function signature, for example if you are adding new function to the `database` structure make sure the interface file `db/interface.go` is updated with the function signature.
2. run the command `mockery` to update the mocks.

#### To create mocks for a new interface make follow the steps below

1. Add the new entry in the `.mockery.yml` file like this

```yml

with-expecter: true
dir: "mocks"
packages:
    github.com/stakwork/sphinx-tribes/db:
        interfaces:
            Database:
    github.com/stakwork/sphinx-tribes/*your-package-name*:
        interfaces:
            *your-interface-1*:
            *your-interface-2*:
```

2. run the command `mockery` to update the mocks.

### Backend API Data Validations

We are currently using `gopkg.in/go-playground/validator.v9` for validation, to validate a struct add the `validate` property to it

```golang
type Workspace struct {
  Name        string     `gorm:"unique;not null" json:"name" validate:"required"`
  Website     string     `json:"website" validate:"omitempty,uri"`
  Github      string     `json:"github" validate:"omitempty,uri"`
  Description string     `json:"description" validate:"omitempty,lte=200"`
}
```

Then handle the validation errors in the request handler

```golang
err = db.Validate.Struct(org)
if err != nil {
  w.WriteHeader(http.StatusBadRequest)
  msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
  json.NewEncoder(w).Encode(msg)
  return
}

```

## Contributing

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the [LICENSE NAME] - see the [LICENSE.md](LICENSE.md) file for details.
