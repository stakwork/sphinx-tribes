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
- [Notification System](#notification-system)

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

- Have Docker installed on your machine
- Run Docker
- Spin up a Postgres DB container before the test with this command
  `docker compose -f docker/testdb-docker-compose.yml -p test_db  up -d`
- Change the rdHost `rdsHost := "172.17.0.1"` variable value in `db/test_config.go` to your `127.0.0.1`

```sh
    // you may need to install cover with this command first
    go get golang.org/x/tools/cmd/cover
    // run test
    RELAY_AUTH_KEY=TEST go test ./... -tags mock -race -v -coverprofile=coverage.out && ./cover-check.sh coverage.out <min coverage amount>
    // To get code coverage in html format do the following after running the code above
    go tool cover -html="coverage.out"
```

- Drop the Postgres DB container after the test with this command
  `docker compose -f docker/testdb-docker-compose.yml -p test_db down`
- Change the rdHost variable value in `db/test_config.go` to the default value for github workflow `rdsHost := "172.17.0.1"`

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

## Notification System

### Auto Notification System for Sphinx Users

We aim to implement a robust notification system to keep users updated about key events on the bounties platform. Notifications will cover events such as bounty assignments, payments, and more.

### System Overview

The notification system will use the v2 bot to send messages to users identified by their public keys (`pubkey`). This document outlines the steps required to build and implement this system effectively.

### Steps to Implement the Notification System

#### 1. Database for Notification Tracking
The first step is to create a database to store notification details, ensuring we can track, retry, and process notifications effectively. Suggested schema:

```sql
CREATE TABLE notifications (
  id SERIAL PRIMARY KEY,
  event VARCHAR(50) NOT NULL,         -- Event type (e.g., bounty_assigned)
  pubkey VARCHAR(100) NOT NULL,       -- User's pubkey
  content TEXT NOT NULL,              -- Notification content
  retries INT DEFAULT 0,              -- Retry count
  status VARCHAR(20) DEFAULT 'PENDING', -- Status: PENDING, COMPLETE, FAILED, WAITING_KEY_EXCHANGE
  uuid UUID NOT NULL,                 -- Unique identifier for the notification
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

This table will store all notification-related details, including their status (`PENDING`, `COMPLETE`, `FAILED`, or `WAITING_KEY_EXCHANGE`), retries, and event data.

#### 2. Verify User on the v2 Bot
Before sending a notification, ensure the v2 bot recognizes the user (`pubkey`):

##### a. Check User in v2 Bot
Call the `/contact/{pubkey}` endpoint:
```typescript
// Response
{
  version: number;
  my_idx: number;
  pubkey: string;
  lsp: string;
  scid: number;
  contact_key?: string;
}
```

- **Condition**: 
  - If `contact_key` is empty, the bot does not recognize the user.

##### b. Add User to Bot’s Known Contacts
If the user is not recognized (`contact_key` is empty), call the `/add_contact` endpoint to add the user to the bot’s contact list.

##### c. Verify Addition
Re-call the `/contact/{pubkey}` endpoint to ensure the `contact_key` is now populated.

- **Handling Missing `contact_key`**:  
  If `contact_key` is still empty at this point, store the notification details in the database with status `WAITING_KEY_EXCHANGE`. A cron job will handle these cases by periodically rechecking the `contact_key`.

#### 3. Send Notification
Once the user is verified, send a notification using the `/send` endpoint:
```typescript
// Request Body for Sending a Message
{
  dest: string;       // User's pubkey
  amt_msat?: number;  // Amount in milli-satoshis (optional, set to 0 for text messages)
  content?: string;   // Message content
  is_tribe?: boolean; // Set to false for direct messages
  reply_uuid?: string; // Optional for replies
  msg_type?: number;  // Optional, depends on use case
  wait?: boolean;     // Set to true for synchronous behavior
}

// Response
{
  status: SendOnionStatus; // "COMPLETE", "PENDING", or "FAILED"
  tag: string;
  preimage?: string;
  payment_hash?: string;
  message?: string; // Error message if status == "FAILED"
}
```

- **Configuration for Text Messages**:
  - `amt_msat`: `0`
  - `is_tribe`: `false`
  - `wait`: `true`

Log the response in the database with the status returned by the endpoint.

#### 4. Notification Processing Function
Create a function to handle the notification lifecycle:

```typescript
function sendNotification(pubkey: string, event: string, content: string, retries: number): void {
  // 1. Verify user on v2 bot:
  //    - Call /contact/{pubkey}.
  //    - If `contact_key` is empty, add the user via /add_contact and re-verify.

  // 2. Send the notification:
  //    - Call /send with the appropriate payload.

  // 3. Log the result:
  //    - Update the database with the current status of the notification.
}
```

#### 5. Cron Job for Key Exchange Handling
Sometimes, the `contact_key` may not be immediately available. Set up a cron job to handle these cases:

- **Behavior**:
  - Periodically check notifications in the database with the status `WAITING_KEY_EXCHANGE`.
  - Re-call the `/contact/{pubkey}` endpoint to recheck for the `contact_key`.
  - If `contact_key` is now available, proceed to send the notification and update the database.

#### 6. Future Enhancements
Although RabbitMQ for message queuing is not implemented in this phase, consider adding it for:
- Improved scalability and reliability.
- Retry logic for failed notifications.
- Batch processing of notification requests.
