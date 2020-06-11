package genstruct

import (
	"fmt"
	"strings"
)

func (m *Module) GenDataModel() (string, string) {
	fmt.Println("Generate " + title(m.Name) + " DataModel.")

	dataStruct := readTemplate("datamodels/datamodels.tmpt")

	dataFeild := ""
	dataForm := ""

	hasTime := false
	hasSql := false
	hasNull := false 
	
		for _, column := range m.Columns {
			
			if !contains(m.UnpublishField, column.Name) {
				valueType := mysqlTypeToGoType(column.Type, column.NullAble == "YES", true)

				if strings.Contains(valueType, "time") {
					hasTime = true
				}

				if strings.Contains(valueType, "sql") {
					hasSql = true
				}

				if strings.Contains(valueType, "null") {
					hasNull = true
				}

				tags := column.GetTag(column.Name, "publish", "")
				dataFeild += fmt.Sprintf("	%v %v %v\n", title(column.Name), valueType, tags)
			}
		}

	if hasOne, ok := m.Relation["hasOne"]; ok {
		for _, relation := range hasOne {
			rel := strings.Split(relation, ":")
			dataFeild += fmt.Sprintf("	%v %v `gorm:\"foreignkey:%v;association_foreignkey:%v\" json:\"%v\"`\n", title(rel[0]), title(rel[0]), title(rel[1]), title(rel[2]), rel[0])
		}
	}

	if hasMany, ok := m.Relation["hasMany"]; ok {
		for _, relation := range hasMany {
			rel := strings.Split(relation, ":")
			dataFeild += fmt.Sprintf("	%v []%v `gorm:\"foreignkey:%v;association_foreignkey:%v\" json:\"%v\"`\n", title(rel[0]), title(rel[0]), title(rel[1]), title(rel[2]), rel[0])
		}
	}


		for _, column := range m.Columns { 
			if !contains(m.UneditableField, column.Name) {
				valueType := mysqlTypeToGoType(column.Type, false, true)

				if strings.Contains(valueType, "time") {
					hasTime = true
				}
				validate := ""
				if val, hasValidate := m.Validator[column.Name]; hasValidate {
					validate = val
				}
				tags := column.GetTag(column.Name, "form", validate) 	

				dataForm += fmt.Sprintf("	%v %v %v\n", title(column.Name), valueType, tags)
			}
		} 


	libImport := ""

	if hasTime || hasSql || hasNull {
		if hasSql {
			libImport += "\"database/sql\"\n"
		}
		if hasTime {
			libImport += "\"time\"\n"
		}
		if hasNull {
			libImport += "\"gopkg.in/guregu/null.v3\"\n"
		}
	}

	dataStruct = strings.Replace(dataStruct, "{{MODULE_NAME}}", title(m.Name), -1)
	dataStruct = strings.Replace(dataStruct, "{{TABLE_NAME}}", m.Name, -1)
	dataStruct = strings.Replace(dataStruct, "{{DATA_FIELD}}", dataFeild, -1)
	dataStruct = strings.Replace(dataStruct, "{{DATA_FORM}}", dataForm, -1)
	// dataStruct = strings.Replace(dataStruct, "{{LIB_IMPORT}}", libImport, -1)

	return dataStruct, libImport
}

func (c *Column) GetTag(key string, typ string, validate string) (tags string) {

	// GORM Tags
	tags = "gorm:\"column:" + c.Name
	if c.Key == "PRI" {
		tags += ";primary_key:true"
	}

	if c.Default != "" {
		tags += ";default:'" + c.Default + "'"
	}

	if c.Type == "PRI" {
		tags += ";primary_key:true"
	}
	tags += "\""

	// JSON Tags
	tags += " json:\"" + key + "\""

	if typ == "form" || typ == "mock" {
		tags += " form:\"" + key + "\""
		
		if validate != "" {
			tags += " validate:\"" + validate + "\""	
		} 
	}

	if typ == "mock" {
		if keyFaker := genFakeTag(key, mysqlTypeToGoType(c.Type, false, true)); keyFaker != "" {
			tags += " faker:\"" + keyFaker + "\""	
		}
			
	}

	tags = strings.Trim(tags, " ")
	tags = "`" + tags + "`"

	return
}
