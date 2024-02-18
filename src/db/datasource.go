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

func (ds *DataSource) MakeDataSource() *gorm.DB {

	once.Do(func() {

		ds.Host = os.Getenv("DB_HOST")
		ds.Port = os.Getenv("DB_PORT")
		ds.Username = os.Getenv("DB_USERNAME")
		ds.Password = os.Getenv("DB_PASSWORD")
		ds.Database = os.Getenv("DB_DATABASE")

		var err error

		if os.Getenv("PROFILE") == "local" {
			dataSourceInstance, err = gorm.Open("mysql", ds.Username+":"+ds.Password+"@tcp("+ds.Host+":"+ds.Port+")/"+ds.Database+"?parseTime=true")
		} else {
			dataSourceInstance, err = gorm.Open("mysql", ds.Username+":"+ds.Password+"@unix(/cloudsql/"+ds.Host+")/"+ds.Database+"?parseTime=true")
		}

		if err != nil {
			panic(err.Error())
		}
	})

	return dataSourceInstance
}

func GetDataSource() *gorm.DB {
	return dataSourceInstance
}

func InitAndGetDataSource() *gorm.DB {
	ds := &DataSource{}
	ds.MakeDataSource()
	return GetDataSource()
}
