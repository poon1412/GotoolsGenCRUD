package genfake

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"goapigen/genstruct"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GenJson(user string, pass string, host string, port string, database string, relationMode bool) {
	db, err := gorm.Open("mysql", user+":"+pass+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln("Failed to connect database")
	}
	defer db.Close()

	dataStruc := genstruct.MainStruct{}

	if jsonBlob, err := ioutil.ReadFile(os.Getenv("PWD") + "/goapigen.json"); err != nil {
		dataStruc = genstruct.MainStruct{
			Version:  "1.0",
			Revision: 0,
			Modules:  make(map[string]genstruct.Module),
		}

	} else {
		err = json.Unmarshal(jsonBlob, &dataStruc)
		check(err)
	}

	//Get all the tables from Database
	rows, err := db.Raw("SHOW TABLES").Rows()
	defer rows.Close()
	for rows.Next() {
		var table string
		rows.Scan(&table)

		fmt.Println("Module: " + table)
		// Store colum as map of maps
		columnDataTypes := make(map[string]genstruct.Column)
		// Select columnd data from INFORMATION_SCHEMA
		columnDataTypeQuery := "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_KEY, COLUMN_DEFAULT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ?"

		rowsCol, err := db.DB().Query(columnDataTypeQuery, database, table)

		if err != nil {
			fmt.Println("Error selecting from db: " + err.Error())
		}

		if rowsCol != nil {
			defer rowsCol.Close()
		} else {
			fmt.Println("No results returned for table")
		}

		uniqueKey := []string{}
		publishField := []string{}
		editableField := []string{}

		for rowsCol.Next() {

			var column string
			var dataType string
			var nullable string
			var columnkey string
			var columndefult string
			rowsCol.Scan(&column, &dataType, &nullable, &columnkey, &columndefult)

			// fmt.Println(fmt.Sprintf("-- Column: %s [%s] %s", column, dataType, nullable))
			columnDataTypes[column] = genstruct.Column{
				Name:     column,
				Type:     dataType,
				NullAble: nullable,
				Key:      columnkey,
				Default:  columndefult,
			}

			publishField = append(publishField, column)
			editableField = append(editableField, column)

			if columnkey == "PRI" {
				uniqueKey = append(uniqueKey, column)
			}
		}

		relation := make(map[string][]string)
		if relationMode {
			relationHasManyQuery := "SELECT TABLE_NAME,COLUMN_NAME,REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = ? AND REFERENCED_TABLE_NAME = ?"
			rowsRelMany, _ := db.DB().Query(relationHasManyQuery, database, table)
			if rowsRelMany != nil {
				defer rowsRelMany.Close()
			} else {
				fmt.Println("No results returned for table")
			}

			for rowsRelMany.Next() {
				var ref_table string
				var ref_col string
				var main_col string

				rowsRelMany.Scan(&ref_table, &ref_col, &main_col)
				relation["hasMany"] = append(relation["hasMany"], ref_table+":"+main_col+":"+ref_col)

			}

			relationHasOneQuery := "SELECT REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME, COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = ? AND TABLE_NAME = ?"
			rowsRelOne, _ := db.DB().Query(relationHasOneQuery, database, table)
			if rowsRelOne != nil {
				defer rowsRelOne.Close()
			} else {
				fmt.Println("No results returned for table")
			}

			for rowsRelOne.Next() {
				var ref_table string
				var ref_col string
				var main_col string

				rowsRelOne.Scan(&ref_table, &ref_col, &main_col)
				relation["hasOne"] = append(relation["hasOne"], ref_table+":"+main_col+":"+ref_col)

			}
		}
	

		if mod, found := dataStruc.Modules[table]; found {
			fmt.Println(fmt.Sprintf("-- Found: %s, Mode %s", table, mod.Mode))
			mod.UniqueKey = uniqueKey
			mod.Columns = columnDataTypes
			dataStruc.Modules[table] = mod
		} else {

			mode := "CRUDL"

			if strings.Index(table, "_media_file") > -1 {
				mode = "M"
			}

			dataStruc.Modules[table] = genstruct.Module{
				Name:               table,
				Mode:               mode,
				InjectionMode:      "",
				PublishField:       publishField,
				EditableField:      editableField,
				FilterList:         []string{},
				FilterDetail:       []string{},
				UniqueKey:          uniqueKey,
				Validator:          make(map[string]string),
				CacheGroupBackend:  []string{},
				CacheGroupFrontend: []string{table},
				CacheDuration:      0,
				Relation:           relation,
				Columns:            columnDataTypes,
			}
		}

		// module[table].Columns = columnDataTypes
	}

	// dataStruc.Modules = module

	jsonString, err := json.MarshalIndent(dataStruc, "", "    ")

	file, err := os.Create(os.Getenv("PWD") + "/goapigen.json")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	fmt.Fprintf(file, string(jsonString))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}