package database

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	_ "github.com/lib/pq"
)

func Init() {
	registerDatabase()
	verifyConnection()
}

func registerDatabase() {
	if err := orm.RegisterDriver("postgres", orm.DRPostgres); err != nil {
		panic(fmt.Errorf("failed to register postgres driver: %w", err))
	}

	user, _ := beego.AppConfig.String("POSTGRES_USER")
	password, _ := beego.AppConfig.String("POSTGRES_PASSWORD")
	dbName, _ := beego.AppConfig.String("POSTGRES_DB")
	host, _ := beego.AppConfig.String("POSTGRES_HOST")
	port, _ := beego.AppConfig.String("POSTGRES_PORT")
	sslMode, _ := beego.AppConfig.String("POSTGRES_SSLMODE")

	dataSource := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		user,
		password,
		dbName,
		host,
		port,
		sslMode,
	)

	if err := orm.RegisterDataBase("default", "postgres", dataSource); err != nil {
		panic(fmt.Errorf("failed to register database: %w", err))
	}
}

func verifyConnection() {
	o := orm.NewOrm()

	if err := o.Raw("SELECT 1").QueryRow(new(int)); err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}
}