package database

import (
	"context"
	"fmt"
	"redi/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

var Pool *pgxpool.Pool

func Config() *pgxpool.Config {
	c, err := pgxpool.ParseConfig("")
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return nil
	}

	c.ConnConfig.Host = config.Config.DatabaseHost
	c.ConnConfig.Port = config.Config.DatabasePort
	c.ConnConfig.Database = config.Config.DatabaseName
	c.ConnConfig.User = config.Config.DatabaseUser
	c.ConnConfig.Password = config.Config.DatabasePassword
	c.ConnConfig.ConnectTimeout = config.Config.DatabaseConnectTimeout

	c.MaxConns = config.Config.DatabaseMaxConns
	c.MinConns = config.Config.DatabaseMinConns
	c.MaxConnLifetime = config.Config.DatabaseMaxConnLifetime
	c.MaxConnIdleTime = config.Config.DatabaseMaxConnIdleTime
	c.HealthCheckPeriod = config.Config.DatabaseHealthCheckPeriod

	c.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		fmt.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	c.AfterRelease = func(conn *pgx.Conn) bool {
		fmt.Println("After releasing the connection pool to the database!!")
		return true
	}

	c.BeforeClose = func(c *pgx.Conn) {
		fmt.Println("Closed the connection pool to the database!!")
	}

	return c
}

func Initialize() error {
	var err error

	Pool, err = pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		return err
	}

	return nil
}
