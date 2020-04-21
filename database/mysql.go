package database

import (
	"fmt"

	"xorm.io/xorm"
)

// DB a pointer to xorm.Engine
var DB *xorm.Engine

// const user string = "root"
// const passwd string = "root"
// const database string = "test"
// const securePivFile string = "C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/base.txt"

// New init global DB
func New(user, passwd, database, securePivFile string) {
	var err error
	url := user + ":" + passwd + "@" + "/" + database
	DB, err = xorm.NewEngine("mysql", url)
	if err != nil {
		fmt.Println(err)
	}
	DB.ShowSQL(true)
}
