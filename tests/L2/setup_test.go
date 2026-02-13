//go:build L2

package e2e

import (
	"context"
	"fmt"
	"glower/initializers"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const baseURL = "http://localhost:8088"

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "glower-db"
	dbUser := "glower-user"
	dbPassword := "glower-password"

	pg, err := postgres.Run(ctx,
		"postgres:18.1-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)

	defer func() {
		if err := testcontainers.TerminateContainer(pg); err != nil {
			log.Printf("Failed to terminate container: %s", err)
		}
	}()

	if err != nil {
		log.Printf("Failed to start container: %s", err)
		return
	}

	host, err := pg.Host(ctx)
	if err != nil {
		log.Printf("Failed to get host data: %s", err)
		return
	}

	port, err := pg.MappedPort(ctx, "5432")
	if err != nil {
		log.Printf("Failed to get port data: %s", err)
		return
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, dbUser, dbPassword, dbName, port.Port(),
	)

	os.Setenv("DB_DSN", dsn)
	app := initializers.NewApp()

	go func() {
		if err := app.Router.Run(":8088"); err != nil {
			log.Fatal(err)
		}
	}()

	m.Run()
}

func TestXxx(t *testing.T) {}
