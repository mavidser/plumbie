package models

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/plumbie/plumbie/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var (
	tables        []interface{}
	Driver        string
	ConnectionStr string
	db            *sqlx.DB
)

func Initialize() error {
	Driver = "postgres"
	ConnectionStr = fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		strings.ReplaceAll(config.Database.Password, "'", "\\'"),
		config.Database.Name)

	log.Debugf("models: Database connection string: %s", ConnectionStr)

	var err error
	db, err = sqlx.Connect(Driver, ConnectionStr)
	if err != nil {
		return err
	}

	db.MapperFunc(ToSnakeCase)

	return nil
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
