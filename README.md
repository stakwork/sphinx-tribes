# sphinx-tribes

![Tribes](https://github.com/stakwork/sphinx-tribes/raw/master/img/sphinx-tribes.png)

Decentralized message broker for public groups in Sphinx. Anyone can run a **sphinx-tribes** server, to route group messages.

**sphinx-tribes** clients can be **sphinx-relay** nodes, apps, websites, or IoT devices.

### How

**sphinx-tribes** is an MQTT broker that any node can subscribe to. Message topics always have two parts: `{receiverPubKey}/{groupUUID}`. Only the owner of the group is allowed to publish to it: all messages from group members must be submitted to the owner as an Lightning keysend payment. The group `uuid` is a timestamp signed by the owner.

![Tribes](https://github.com/stakwork/sphinx-tribes/raw/master/img/tribes.jpg)

### Authentication

Authentication is handled by [sphinx-auth](https://github.com/stakwork/sphinx-auth)

### build

docker build --no-cache -t sphinx-tribes .

docker tag sphinx-tribes sphinxlightning/sphinx-tribes:x

docker push sphinxlightning/sphinx-tribes:x

### run against sphinx-stack

To run tribes frontend locally, use these ports:

- tribes: `yarn start:tribes:docker` (localhost:23000)
- people: `yarn start:people:docker` (localhost:23007)

### Run frontend locally against people.sphinx.chat

If you would like to run just the frontend do the following

line 10 in `frontend/app/src/config/ModeDispatcher.tsx` change `'localhost:3000': Mode.TRIBES` -> `'localhost:3000': Mode.COMMUNITY`

line 5 in `frontend/app/src/config/host.ts` return `"people.sphinx.chat"`

### Run Golang backemd

- Create a .env file and populate the env file with thgitese variables

```
    TEST_MODE=true
    PORT=  // Server PORT
    DATABASE_URL= // Database URL
    TEST_ASSET_URL= // Test asset URL
```

- Build the application by running `go build .`
- Run the application with `./sphinx-tribes`

### Run Golang backend with local database

- Create a Postgres database
- Copy the table creation queries from tribes.sql file, and create the tables in the database created
- Create a .env file and populate the .env files with these variables

```
    RDS_HOSTNAME =
    RDS_PORT =
    RDS_DB_NAME =
    RDS_USERNAME =
    RDS_PASSWORD =
```

- Build and run the Golang server

### Run Test

```
    // you may need to install cover with this command first
    // go get golang.org/x/tools/cmd/cover
    go test ./... -tags mock -race -v -coverprofile=coverage.out && ./cover-check.sh coverage.out <min coverage amount> 
    // To get code coverage in html format do the following after running the code above
    go tool cover -html="coverage.out"
```

### Enable Relay

If you would like to enable Relay for invoice creation and keysend payment add the relay env key and values to the .env file

```
    RELAY_URL=
    RELAY_AUTH_KEY=
```

### Enable Meme Image Upload

If you would like to enable Meme image upload for organization's add the meme env key and value to the .env file,
Meme image upload works with Relay enabled, so a running Relay is required for Meme to work.

```
    MEME_URL=
```

### Add REDIS for cache

- Create a Redis instance
- Create a .env file and populate the .env files with these variables

If you have a Redis url add the REDIS_URL variable to .env

```REDIS_URL = ```

else add these variables to the env to enable Redis

```
    REDIS_HOST =
    REDIS_DB =
    REDIS_USER =
    REDIS_PASS =
```

### Add SuperAdmins to access admin dashboard

Add comma separated public keys to the SUPER_ADMINS env var in the .env file,
any user public key added to this comaa separated strings will have access to the admin dashboard
e.g '{pubkey}, {pubkey}, {pubkey}'

```
ADMINS
```

### For Contributions

Read the contribution doc [here](./docs/Contribution.md)

### Stakwork Youtube vidoes download for tribes feed

If you would like to enable Stakwork jobs for Youtube videos download add the Stakwork env key and values to the .env file

```
    STAKWORK_KEY=
```
