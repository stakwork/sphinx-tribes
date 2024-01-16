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

### Run Golang backend

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


### Mocking interfaces for unit testing golang backend
We are using [mockery](https://vektra.github.io/mockery/latest/) to autogenerate mocks of an interface in our unit tests.

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

#### When adding a new function to the interface which is already mocked follow the below steps:
1. Update the corresponding interface with the function signature, for example if you are adding new function to the ```database``` structure make sure the interface file ```db/interface.go``` is updated with the function signature.
2. run the command ```mockery``` to update the mocks.


#### To create mocks for a new interface make follow the steps below:

1. Add the new entry in the ```.mockery.yml``` file like this
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
2. run the command ```mockery``` to update the mocks.


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

### Backend API Data Validations

We are currently using `gopkg.in/go-playground/validator.v9` for validation, to validate a struct add the `validate` property to it

e.g 

```
type Organization struct {
	Name        string     `gorm:"unique;not null" json:"name" validate:"required"`
	Website     string     `json:"website" validate:"omitempty,uri"`
	Github      string     `json:"github" validate:"omitempty,uri"`
	Description string     `json:"description" validate:"omitempty,lte=200"`
}
```

Then handle the validation errors in the request handler 

```
err = db.Validate.Struct(org)
if err != nil {
	w.WriteHeader(http.StatusBadRequest)
	msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
	json.NewEncoder(w).Encode(msg)
	return
}

```


