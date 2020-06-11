package genstruct

import (
	"fmt"
	"strings"
)

func (m *Module) GenModel() string {
	fmt.Println("Generate " + title(m.Name) + " Model.")

	dataStruct := readTemplate("models/model.head.tmpt")

	datamodels, importAddition := m.GenDataModel()

	dataStruct = strings.Replace(dataStruct, "{{DATA_MODELS}}", datamodels, -1)

	 
	if strings.Index(m.Mode, "M") > -1 {
		dataStruct += readTemplate("models/model.media.tmpt")
	} else {

	// List 
		dataStruct += readTemplate("models/model.list.tmpt")

		filterInit := ""
		filterParam := ""
		if len(m.FilterList) > 0 {
			for _, filterColName := range m.FilterList {
				if column, ok := m.Columns[filterColName]; ok {

					typeConverted := convertStructType(column.Type)
					if typeConverted == "time.Time" {
						typeConverted = "string"
					}
					defaultType := defultStructType(column.Type)

					filterParam += fmt.Sprintf("%v %v, ", filterColName, typeConverted)
					filterInit += fmt.Sprintf("if %v != %v { dbCnn = dbCnn.Where(\"%v = ?\", %v) }\n", filterColName, defaultType, filterColName, filterColName)

				}
			}
		}

		if hasOne, ok := m.Relation["hasOne"]; ok {
			for _, relation := range hasOne {
				rel := strings.Split(relation, ":")
				filterInit += fmt.Sprintf("	dbCnn = dbCnn.Preload(\"%v\")\n", title(rel[0]))
			}
		}

		if hasMany, ok := m.Relation["hasMany"]; ok {
			for _, relation := range hasMany {
				rel := strings.Split(relation, ":")
				filterInit += fmt.Sprintf("	dbCnn = dbCnn.Preload(\"%v\")\n", title(rel[0]))
			}
		}

		dataStruct = strings.Replace(dataStruct, "{{FILTER_INIT}}", filterInit, -1)
		dataStruct = strings.Replace(dataStruct, "{{FILTER_PARAM}}", filterParam, -1)
 		
	// Read
		dataStruct += readTemplate("models/model.read.tmpt")
		preload := "" 
		filterInit = ""
		if hasOne, ok := m.Relation["hasOne"]; ok {
			for _, relation := range hasOne {
				rel := strings.Split(relation, ":")
				preload += fmt.Sprintf("	dbCnn = dbCnn.Preload(\"%v\")\n", title(rel[0]))
			}
		}

		if hasMany, ok := m.Relation["hasMany"]; ok {
			for _, relation := range hasMany {
				rel := strings.Split(relation, ":")
				preload += fmt.Sprintf("	dbCnn = dbCnn.Preload(\"%v\")\n", title(rel[0]))
			}
		} 

		if len(m.UniqueKey) > 0 { 
			for _, uniqueKey := range m.UniqueKey {
				if column, ok := m.Columns[uniqueKey]; ok {

					valueType := mysqlTypeToGoType(column.Type, column.NullAble == "YES", true)
					if strings.Contains(valueType, "null.String") {
						filterInit += fmt.Sprintf(" dbCnn = dbCnn.Where(\"%v = ?\", obj.%v.String)\n", column.Name, title(column.Name))
					} else if strings.Contains(valueType, "null.Int") {
					 	filterInit += fmt.Sprintf(" dbCnn = dbCnn.Where(\"%v = ?\", obj.%v.Int)\n", column.Name, title(column.Name))
					} else {
						filterInit += fmt.Sprintf(" dbCnn = dbCnn.Where(\"%v = ?\", obj.%v)\n", column.Name, title(column.Name))
					}
				}
			}
		}
		dataStruct = strings.Replace(dataStruct, "{{PRELOAD}}", preload, -1) 
		dataStruct = strings.Replace(dataStruct, "{{FILTER_INIT}}", filterInit, -1) 
		
	// Create
		dataStruct += readTemplate("models/model.create.tmpt")
	
	// Update
		dataStruct += readTemplate("models/model.update.tmpt")
		filterInit = "" 
		// if len(m.UniqueKey) > 0 { 
		// 	for _, uniqueKey := range m.UniqueKey {
		// 		if column, ok := m.Columns[uniqueKey]; ok {
		// 			filterInit += fmt.Sprintf(" dbCnn = dbCnn.Where(\"%v = ?\", obj.%v)\n", column.Name, title(column.Name))
		// 		}
		// 	}
		// }
		for _, column := range m.Columns { 
			if column.Key == "PRI" {
				filterInit += fmt.Sprintf(" dbCnn = dbCnn.Where(\"%v = ?\", obj.%v)\n", column.Name, title(column.Name))
			}
		}

		dataStruct = strings.Replace(dataStruct, "{{FILTER_INIT}}", filterInit, -1) 
	
	// Delete
		dataStruct += readTemplate("models/model.delete.tmpt")
		filterInit = ""
		if len(m.UniqueKey) > 0 {
			for _, uniqueKey := range m.UniqueKey {
				if column, ok := m.Columns[uniqueKey]; ok { 
					filterInit += fmt.Sprintf(" dbCnn = dbCnn.Where(\"%v = ?\", obj.%v)\n", column.Name, title(column.Name)) 
				}
			}  
		} 
		dataStruct = strings.Replace(dataStruct, "{{FILTER_INIT}}", filterInit, -1) 
	}

	

	dataStruct = strings.Replace(dataStruct, "{{MODULE_NAME}}", title(m.Name), -1)
	dataStruct = strings.Replace(dataStruct, "{{IMPORT_ADDITION}}", importAddition, -1)

	return dataStruct 
}
