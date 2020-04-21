package service

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"lipeining.com/poem-gin/config"
	"lipeining.com/poem-gin/model"
)

// BuildCreateTableCommand just sql
func BuildCreateTableCommand(tblTmpIDName, tblName string, cols []model.Column) (string, string, error) {
	database := config.Config.Xorm.Database
	fullTblTmpIDName := database + "." + tblTmpIDName
	fullTblName := database + "." + tblName
	insertList := make([]string, 0)
	insertList = append(insertList, "INSERT INTO "+fullTblName+" SELECT ")
	insertCols := make([]string, 0)
	insertCols = append(insertCols, "id")
	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(fullTblName).IfNotExists()
	ctb.Define("id", "BIGINT(20)", "NOT NULL", "PRIMARY KEY", "AUTO_INCREMENT", `COMMENT "id"`)
	// todo 需要检查正确性
	// todo 丰富的 column 属性 默认值，大小，
	for _, column := range cols {
		name := column.Name
		if column.T == "string" {
			insertCols = append(insertCols, " uuid() ")
			ctb.Define(name, "VARCHAR(255)", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "datetime" {
			insertCols = append(insertCols, " date_add(NOW(), interval FLOOR(1 + (RAND() * 10)) month) ")
			ctb.Define(name, "datetime", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "int64" {
			insertCols = append(insertCols, " FLOOR(1 + (RAND() * 100)) ")
			ctb.Define(name, "BIGINT(20)", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "int32" {
			insertCols = append(insertCols, " FLOOR(1 + (RAND() * 100)) ")
			ctb.Define(name, "INT(11)", "NOT NULL", `COMMENT "`+name+`"`)
		}
	}
	insertList = append(insertList, strings.Join(insertCols, ","))
	insertList = append(insertList, " FROM "+fullTblTmpIDName+";")
	ctb.Option("DEFAULT CHARACTER SET", "utf8mb4")
	fmt.Println(ctb)
	createTableCommand := ctb.String()
	insertDataCommand := strings.Join(insertList, "")
	return createTableCommand, insertDataCommand, nil
}

// BuildLoadIntoTableCommand just sql
func BuildLoadIntoTableCommand(fullTblTmpIDName string) string {
	return "load data infile " + "'" + config.Config.Xorm.SecurePivFile + "'" + " replace into table " + fullTblTmpIDName + ";"
}

// WriteInfile write 1,000,000
func WriteInfile(fileName string, num int) {
	// dir := filepath.Dir(securePivFile)
	// filePath := filepath.Join(dir, fileName)
	filePath := fileName
	outputFile, outputError := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Println("An error occurred with file opening or creation\n", outputError)
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	// 每次写 10000 个数字
	every := 10000
	total := num / every
	if num%every != 0 {
		total++
	}
	current := 0
	for i := 0; i < total; i++ {
		for j := 0; j < every; j++ {
			current++
			outputString := strconv.Itoa(current)
			if current != 1 {
				outputString = "\n" + outputString
			}
			outputWriter.WriteString(outputString)
		}
		outputWriter.Flush()
	}
}
