package db

import (
	"fmt"
	"log"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitTestDocker function initialize docker with postgres image used for integration tests
func InitTestDocker(exposedPort string) (*dockertest.Pool, *dockertest.Resource) {
	var passwordEnv = "postgres"
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"listen_addresses = '*'",
			fmt.Sprint(passwordEnv),
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
		},
		ExposedPorts: []string{exposedPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {
				{HostIP: "0.0.0.0", HostPort: exposedPort},
			},
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"} // Important option when container crash and you want to debug container
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// if err := resource.Expire(600); err != nil { // Tell docker to hard kill the container in 30 seconds
	// 	logrus.Error(err)
	// }

	// retry if container is not ready
	pool.MaxWait = 60 * time.Second
	if err = pool.Retry(func() error {
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return pool, resource
}

func OpenDatabaseConnection(pool *dockertest.Pool, resource *dockertest.Resource, exposedPort string) *gorm.DB {
	// Wait for 30 seconds for the databaase to spin up on docker
	time.Sleep(30 * time.Second)

	user := "test"
	password := "test"
	db := "postgres"
	port := "5432"
	dns := "host=%s port=%s user=%s sslmode=disable password=%s dbname=%s"

	retries := 5
	host := resource.GetBoundIP(fmt.Sprintf("%s/tcp", port))
	gdb, err := gorm.Open(postgres.Open(fmt.Sprintf(dns, host, exposedPort, user, password, db)), &gorm.Config{})

	// Sometimes it happens that after first time container is not ready.
	// It's always better to create retry if necessary and be sure that tests run without problems
	// retry every 6 seconds 5 times
	for err != nil {
		if retries > 1 {
			retries--
			time.Sleep(6 * time.Second)
			gdb, err = gorm.Open(postgres.Open(fmt.Sprintf(dns, host, exposedPort, user, password, db)), &gorm.Config{})
			continue
		}

		if err := pool.Purge(resource); err != nil {
			logrus.Error(err)
		}

		log.Panic("Fatal error in connection: ", err, resource.GetBoundIP("5432/tcp"))
	}

	return gdb
}

func MigrateTestDb(db *gorm.DB) {
	db.AutoMigrate(&Tribe{})
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Channel{})
	db.AutoMigrate(&LeaderBoard{})
	db.AutoMigrate(&ConnectionCodes{})
	db.AutoMigrate(&BountyRoles{})
	db.AutoMigrate(&UserInvoiceData{})
	db.AutoMigrate(&WorkspaceRepositories{})
	db.AutoMigrate(&WorkspaceFeatures{})
	db.AutoMigrate(&FeaturePhase{})
	db.AutoMigrate(&FeatureStory{})
	db.AutoMigrate(&NewBounty{})
	db.AutoMigrate(&BudgetHistory{})
	db.AutoMigrate(&NewPaymentHistory{})
	db.AutoMigrate(&NewInvoiceList{})
	db.AutoMigrate(&NewBountyBudget{})
	db.AutoMigrate(&Workspace{})
	db.AutoMigrate(&WorkspaceUsers{})
	db.AutoMigrate(&WorkspaceUserRoles{})
}

func StartTestDb() {
	exposedPort := fmt.Sprint(5532)
	pool, resource := InitTestDocker(exposedPort)
	if pool != nil && resource != nil {
		gdb := OpenDatabaseConnection(pool, resource, exposedPort)
		MigrateTestDb(gdb)
	}
}
