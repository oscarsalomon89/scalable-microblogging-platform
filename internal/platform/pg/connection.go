package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oscarsalomon89/go-hexagonal/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connections struct {
	MasterConn *gorm.DB `name:"writeConn"`
}

func NewDBConnections(cfg config.Database) (Connections, error) {
	dialector, err := getDialector(cfg)
	if err != nil {
		return Connections{}, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return Connections{}, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return Connections{}, err
	}

	configureDBs(sqlDB, cfg)

	return Connections{
		MasterConn: db,
	}, nil
}

func getDialector(cfg config.Database) (gorm.Dialector, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.Sslmode)

	return postgres.Open(dsn), nil
}

func configureDBs(conn *sql.DB, appCfg config.Database) {
	conn.SetMaxOpenConns(appCfg.MaxConnections)
	conn.SetMaxIdleConns(appCfg.MaxIdleConnections)
	maxIdleTime := appCfg.MaxIdleTime
	maxLifeTime := appCfg.MaxLifeTime
	conn.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Second)
	conn.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Second)
}
