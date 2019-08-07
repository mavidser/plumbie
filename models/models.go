package models

import (
	"fmt"
	"strings"

	"github.com/plumbie/plumbie/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

var (
	tables     []interface{}
	x          *xorm.Engine
	Connection string
)

func init() {
	tables = append(tables,
		new(User),
		new(Session),
	)
}

func Initialize() error {
	var paramSeparator = "?"
	if strings.Contains(config.Database.Name, paramSeparator) {
		paramSeparator = "&"
	}
	Connection = fmt.Sprintf("%s:%s@%s(%s)/%s%scharset=%s&parseTime=true&tls=%v",
		config.Database.User,
		config.Database.Password,
		config.Database.Protocol,
		config.Database.Host,
		config.Database.Name,
		paramSeparator,
		config.Database.Charset,
		config.Database.SSL)

	log.Debugf("models: Database connection string: %s", Connection)
	var err error
	if x, err = xorm.NewEngine(config.Database.Driver, Connection); err != nil {
		return err
	}

	x.ShowSQL(config.Debug)

	if err = x.Sync2(tables...); err != nil {
		return err
	}

	return nil
}
