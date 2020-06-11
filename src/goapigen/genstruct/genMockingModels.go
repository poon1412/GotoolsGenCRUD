package genstruct

import (
	"fmt"
	"strings"
)

func (m *Module) GenMockingModels() string {
	fmt.Println("Generate asdfds " + title(m.Name) + " MockModel.")

	dataStruct := readTemplate("mock/mockmodel.tmpt")

	dataForm := ""

	hasTime := false
	hasSql := false
	hasNull := false
 
	
		for _, column := range m.Columns { 
			if !contains(m.UneditableField, column.Name) {
				valueType := mysqlTypeToGoType(column.Type, false, true)

				if strings.Contains(valueType, "time") {
					// hasTime = true
					valueType = "string"
				}
				validate := ""
				if val, hasValidate := m.Validator[column.Name]; hasValidate {
					validate = val
				}
				tags := column.GetTag(column.Name, "mock", validate) 	

				dataForm += fmt.Sprintf("	%v %v %v\n", title(column.Name), valueType, tags)
			}
		}
	

	libImport := ""
	libImport = "import(\n"
	libImport = libImport + "\"fmt\"\n"
	libImport = libImport + "\"github.com/bxcodec/faker\"\n"
	// libImport = libImport + "\"server/models\"\n" 
	if hasTime || hasSql || hasNull {
		
		if hasSql {
			libImport = libImport + "\"database/sql\"\n"
		}
		if hasTime {
			libImport = libImport + "\"time\"\n"
		}
		if hasNull {
			libImport = libImport + "\"gopkg.in/guregu/null.v3\"\n"
		} 
	}

	libImport = libImport + ")"
	dataStruct = strings.Replace(dataStruct, "{{MODULE_NAME}}", title(m.Name), -1)
	dataStruct = strings.Replace(dataStruct, "{{TABLE_NAME}}", m.Name, -1)
	dataStruct = strings.Replace(dataStruct, "{{DATA_FORM}}", dataForm, -1)
	dataStruct = strings.Replace(dataStruct, "{{LIB_IMPORT}}", libImport, -1)

	return dataStruct
}


func genFakeTag(key string, typeC string) string {
	if strings.Index(key, "first_name") > -1 && strings.Index(typeC, "string") > -1  {
		return "first_name_male"
	}

	if strings.Index(key, "last_name") > -1 && strings.Index(typeC, "string") > -1  {
		return "last_name"
	}

	if strings.Index(key, "name") > -1 && strings.Index(typeC, "string") > -1  {
		return "name"
	}

	if strings.Index(key, "phone") > -1 && strings.Index(typeC, "string") > -1  {
		return "phone_number"
	}

	if key == "tel" && strings.Index(typeC, "string") > -1  {
		return "phone_number"
	}

	if strings.Index(key, "title") > -1 && strings.Index(typeC, "string") > -1  {
		return "sentence"
	}

	if strings.Index(key, "desc") > -1 && strings.Index(typeC, "string") > -1  {
		return "sentences"
	}

	if strings.Index(key, "type") > -1 && strings.Index(typeC, "string") > -1 {
		return "word"
	}

	if strings.Index(key, "date") > -1  && strings.Index(typeC, "string") > -1 {
		return "timestamp"
	}
	if strings.Index(key, "create") > -1  && strings.Index(typeC, "time") > -1  {
		return "timestamp"
	}
	if strings.Index(key, "update") > -1  && strings.Index(typeC, "time") > -1 {
		return "timestamp"
	}

	if strings.Index(key, "link") > -1  && strings.Index(typeC, "string") > -1{
		return "url"
	}
	if strings.Index(key, "url") > -1  && strings.Index(typeC, "string") > -1 {
		return "url"
	}

	if strings.Index(key, "email") > -1  && strings.Index(typeC, "string") > -1 {
		return "email"
	}
	return ""
}




