package db

import (
	"archivit_Backend/src/db/entity"
)

func GormMigrate() {

	ds := dataSourceInstance

	ds.AutoMigrate(&entity.User{})

}
