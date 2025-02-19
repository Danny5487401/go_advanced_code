package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
)

func TestInsertRowAndQuert(t *testing.T) {
	insertRowDemo()
	queryRowDemo()
}

const (
	dbUsername string = "root"
	dbPassword string = "password"
	dbName     string = "test"
)

var createSql = `
CREATE TABLE Person (
	Id integer auto_increment NOT NULL,
	Name VARCHAR(30) NULL,
	City VARCHAR(50) NULL,
	AddTime DATETIME NOT NULL,
	UpdateTime DATETIME NOT NULL,
	CONSTRAINT Person_PK PRIMARY KEY (Id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_general_ci;

`

func TestMain(m *testing.M) {
	closeFunc, newDb, err := SetupMySQLContainer()
	db = newDb
	if err != nil {
		log.Fatalf("could not setup MySQL container: %v", err)
	}
	defer closeFunc()
	m.Run()
}

func SetupMySQLContainer() (func(), *sqlx.DB, error) {
	log.Printf("setup MySQL Container")
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": dbPassword,
			"MYSQL_DATABASE":      dbName,
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL"),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Panicf("error starting mysql container: %s", err)
	}

	closeContainer := func() {
		log.Printf("terminating container")
		err := mysqlC.Terminate(ctx)
		if err != nil {
			log.Panicf("error terminating mysql container: %s", err)
		}
	}

	host, _ := mysqlC.Host(ctx)
	p, _ := mysqlC.MappedPort(ctx, "3306/tcp")
	port := p.Int()

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUsername, dbPassword, host, port, dbName)

	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		log.Printf("error connect to db: %+v\n", err)
		return closeContainer, db, err
	}

	if err = db.Ping(); err != nil {
		log.Panicf("error pinging db: %+v\n", err)
		return closeContainer, db, err
	}
	if _, err = db.ExecContext(ctx, createSql); err != nil {
		log.Panicf("error create db: %+v\n", err)
		return closeContainer, db, err
	}

	return closeContainer, db, nil
}
