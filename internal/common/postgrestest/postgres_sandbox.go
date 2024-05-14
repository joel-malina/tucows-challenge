package postgrestest

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreatePostgresSandbox(t testing.TB) *sqlx.DB {
	client := resty.New().SetBaseURL("http://localhost:8081")
	response, err := client.NewRequest().Get("/pgsandbox")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode() != http.StatusOK {
		t.Fatal("did not get a success response code creating sandbox", response.StatusCode(), string(response.Body()))
	}

	datasourceName := string(response.Body())
	datasourceName = strings.Replace(datasourceName, "host=postgres", "host=localhost", 1)

	db, err := sqlx.Connect("postgres", datasourceName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("created sandbox", datasourceName)

	return db
}
