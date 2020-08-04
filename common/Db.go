package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// Db
var Db *gorm.DB

func GetDbConfigString() string {
	dbConfig := Config.GetStringMapString("db")
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig["username"], dbConfig["password"], dbConfig["addr"], dbConfig["name"])
}

