package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/huandu/go-sqlbuilder"
	"lipeining.com/poem-gin/config"
	"lipeining.com/poem-gin/database"
	"lipeining.com/poem-gin/model"
	"lipeining.com/poem-gin/service"
)

// Ping -pong
func Ping(c *gin.Context) {
	c.String(200, "pong")
}

// Home home page show config
func Home(c *gin.Context) {
	fmt.Printf("config %v", config.Config)
	// c.String(http.StatusOK, "hello World!")
	c.JSON(http.StatusOK, config.Config)
}

// TblTmpID use struct and xorm create table
func TblTmpID(c *gin.Context) {
	exisit, err := database.DB.IsTableExist(model.TblTmpID{})
	if err != nil {
		fmt.Println(err)
	}
	if exisit {
		err = database.DB.DropTables(model.TblTmpID{})
		if err != nil {
			fmt.Println(err)
		}
	}
	err = database.DB.CreateTables(model.TblTmpID{})
	if err != nil {
		fmt.Println(err)
	}
	loadDataCommand := "load data infile 'C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/base.txt' replace into table test.tbl_tmp_id;"
	results, err := database.DB.Query(loadDataCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(results)
	c.JSON(http.StatusOK, results)
}

// TmpID use sql builder to create table
func TmpID(c *gin.Context) {
	tblTmpIDName := c.DefaultQuery("tmp_name", "tbl_tmp_id")
	ctb := sqlbuilder.NewCreateTableBuilder()
	fullTblTmpIDName := config.Config.Xorm.Database + "." + tblTmpIDName
	ctb.CreateTable(fullTblTmpIDName).IfNotExists()
	ctb.Define("id", "BIGINT(20)", "NOT NULL", "PRIMARY KEY", `COMMENT "id"`)
	ctb.Option("DEFAULT CHARACTER SET", "utf8mb4")
	fmt.Println(ctb)
	// c.JSON(http.StatusOK, ctb)
	// c.String(http.StatusOK, ctb.String())
	insertCommand := ctb.String()
	results, err := database.DB.Query(insertCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("create table results", results)
	loadDataCommand := "load data infile " + "'" + config.Config.Xorm.SecurePivFile + "'" + " replace into table " + fullTblTmpIDName + ";"
	results, err = database.DB.Query(loadDataCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("load data results", results)
	c.JSON(http.StatusOK, results)
}

// NewTbl use tmp_name tbl_name and sql builder to create and insert data
func NewTbl(c *gin.Context) {
	// http://localhost:8000/new_tbl?tmp_name=abc&tbl_name=user&cols=province,city,name&times=create_time,update_time,delete_time
	// http://localhost:8000/new_tbl?tmp_name=abc&tbl_name=food&cols=status,voted,name&times=create_time,update_time,delete_time
	tblTmpIDName := c.DefaultQuery("tmp_name", "tbl_tmp_id")
	fullTblTmpIDName := config.Config.Xorm.Database + "." + tblTmpIDName
	tblName := c.Query("tbl_name")
	fullTblName := config.Config.Xorm.Database + "." + tblName
	cols := c.Query("cols")
	columns := strings.Split(cols, ",")
	times := c.Query("times")
	timeCols := strings.Split(times, ",")
	insertList := make([]string, 0)
	insertList = append(insertList, "INSERT INTO "+fullTblName+" SELECT ")
	insertCols := make([]string, 0)
	insertCols = append(insertCols, "id")
	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(fullTblName).IfNotExists()
	ctb.Define("id", "BIGINT(20)", "NOT NULL", "PRIMARY KEY", "AUTO_INCREMENT", `COMMENT "id"`)
	for _, column := range columns {
		insertCols = append(insertCols, " uuid() ")
		ctb.Define(column, "VARCHAR(255)", "NOT NULL", `COMMENT "`+column+`"`)
	}
	for _, column := range timeCols {
		// insertCols = append(insertCols, " NOW() ")
		insertCols = append(insertCols, " date_add(NOW(), interval FLOOR(1 + (RAND() * 10)) month) ")
		ctb.Define(column, "datetime", "NOT NULL", `COMMENT "`+column+`"`)
	}
	insertList = append(insertList, strings.Join(insertCols, ","))
	insertList = append(insertList, " FROM "+fullTblTmpIDName+";")
	ctb.Option("DEFAULT CHARACTER SET", "utf8mb4")
	fmt.Println(ctb)
	insertCommand := ctb.String()
	results, err := database.DB.Query(insertCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("create table results", results)
	// 	INSERT INTO tbl_user
	// 	SELECT
	// 	  id,
	// 	  uuid(),
	// 	  CONCAT('userNickName', id),
	// 	  FLOOR(Rand() * 1000),
	// 	  FLOOR(Rand() * 100),
	// 	  NOW()
	// 	FROM
	//   tbl_tmp_id;
	insertDataCommand := strings.Join(insertList, "")
	results, err = database.DB.Query(insertDataCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("insert data results", results)
	c.JSON(http.StatusOK, results)
}

// BuildTmpTable use service.func create table
func BuildTmpTable(c *gin.Context) {
	// http://localhost:8000/build_tmp_table?tmp_name=efg
	tblTmpIDName := c.DefaultQuery("tmp_name", "tbl_tmp_id")
	cols := make([]model.Column, 0)
	createTableCommand, _, err := service.BuildCreateTableCommand(tblTmpIDName, tblTmpIDName, cols)
	results, err := database.DB.Query(createTableCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("create table results", results)
	loadIntoTableCommand := service.BuildLoadIntoTableCommand(tblTmpIDName)
	results, err = database.DB.Query(loadIntoTableCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("load into data results", results)
	c.JSON(http.StatusOK, results)
}

// BuildTable use service.func create table
func BuildTable(c *gin.Context) {
	// http://localhost:8000/build_table?tmp_name=efg&tbl_name=teacher&cols=province,city,name&times=create_time,update_time,delete_time
	// http://localhost:8000/build_table?tmp_name=efg&tbl_name=student&cols=province,city,name&times=create_time,update_time,delete_time
	tblTmpIDName := c.DefaultQuery("tmp_name", "tbl_tmp_id")
	tblName := c.Query("tbl_name")
	cols := strings.Split(c.Query("cols"), ",")
	timeCols := strings.Split(c.Query("times"), ",")
	columns := make([]model.Column, 0)
	for _, col := range cols {
		columns = append(columns, model.Column{
			Name: col,
			T:    "string",
		})
	}
	for _, col := range timeCols {
		columns = append(columns, model.Column{
			Name: col,
			T:    "datetime",
		})
	}

	createTableCommand, insertDataCommand, err := service.BuildCreateTableCommand(tblTmpIDName, tblName, columns)
	results, err := database.DB.Query(createTableCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("create table results", results)
	results, err = database.DB.Query(insertDataCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("insert into data results", results)
	c.JSON(http.StatusOK, results)
	// 	INSERT INTO tbl_user
	// 	SELECT
	// 	  id,
	// 	  uuid(),
	// 	  CONCAT('userNickName', id),
	// 	  FLOOR(Rand() * 1000),
	// 	  FLOOR(Rand() * 100),
	// 	  NOW()
	// 	FROM
	//   tbl_tmp_id;
}

// Infile create base.txt for 1,000,000
func Infile(c *gin.Context) {
	// http://localhost:8000/infile?filename=counter.txt&num=1000000
	filename := c.DefaultQuery("filename", "counter.txt")
	numInput := c.DefaultQuery("num", "1000000")
	num, ok := strconv.Atoi(numInput)
	if ok != nil {
		num = 1000000
	}
	fmt.Println(filename, num)
	service.WriteInfile(filename, num)
	c.String(http.StatusOK, "ok")
}
