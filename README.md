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
    go test ./...
```

### Enable Relay

If you would like to enable Relay for invoice creation and keysend payment add the relay env key and values to the .env file

```
    RELAY_URL=
    RELAY_AUTH_KEY=
```

### Enable Meme Image Upload

If you would like to enable Meme image upload for organization's add the meme env keyd and values to the .env file

```
    MEME_URL=
    RELAY_NODE_KEY=
```

### For Contributions

Read the contribution doc [here](./Contribution.md)

### Stakwork Youtube vidoes download for tribes feed

If you would like to enable Stakwork jobs for Youtube videos download add the Stakwork env key and values to the .env file

```
    STAKWORK_KEY=
```
