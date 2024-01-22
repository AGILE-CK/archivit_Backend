package db

import (
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jinzhu/gorm"
	"os"
	"sync"
)

type DataSource struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

var dataSourceInstance *gorm.DB
var once sync.Once

func (ds *DataSource) GetDataSource() *gorm.DB {

	once.Do(func() {

		if os.Getenv("PROFILE") == "develop" {
			ds.Host = os.Getenv("DB_HOST")
			ds.Port = os.Getenv("DB_PORT")
			ds.Username = os.Getenv("DB_USERNAME")
			ds.Password = os.Getenv("DB_PASSWORD")
			ds.Database = os.Getenv("DB_DATABASE")

		} else {
			ds.Host = "localhost"
			ds.Port = "5535"
			ds.Username = "agile"
			ds.Password = "agile"
			ds.Database = "agile_database"
		}

		var err error

		if os.Getenv("PROFILE") == "" {
			dataSourceInstance, err = gorm.Open("mysql", ds.Username+":"+ds.Password+"@tcp("+ds.Host+":"+ds.Port+")/"+ds.Database)
		} else {
			dataSourceInstance, err = gorm.Open("mysql", ds.Username+":"+ds.Password+"@unix(/cloudsql/"+ds.Host+")/"+ds.Database)
		}

		if err != nil {
			panic(err.Error())
		}
	})

	return dataSourceInstance
}
