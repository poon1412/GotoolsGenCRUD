package genstruct

import (
	"fmt"
	"strings"
)

func (m *Module) GenUnitTest() (string, bool) {

	hasUnittest := false

	testStruct := "/* UnitTest {{MODULE_NAME}} */\n"
	
	testStruct = readTemplate("unittest/unittest.tmpt") 

	dataStruct := ""

	if strings.Index(m.Mode, "C") > -1 {
		hasUnittest = true
		dataStruct += readTemplate("unittest/unittest.create.tmpt")
		dataForm :=""
		for _, column := range m.Columns { 
			if !contains(m.UneditableField, column.Name) {
				valueType := mysqlTypeToGoType(column.Type, false, true) 
				if strings.Contains(valueType, "time") {
					dataForm += fmt.Sprintf("WithFormField(\"%v\", DateTime.ConvFormat({{MODULE_NAME}}CreateStruct.%v)). \n", column.Name, title(column.Name))
				} else {
					dataForm += fmt.Sprintf("WithFormField(\"%v\", {{MODULE_NAME}}CreateStruct.%v). \n", column.Name, title(column.Name))
				}
				
			}
		}

		dataStruct = strings.Replace(dataStruct, "{{FORM_FIELD}}", dataForm, -1)
	}



	if strings.Index(m.Mode, "L") > -1 && strings.Index(m.InjectionMode, "L") <= -1  {
		hasUnittest = true
		dataStruct += readTemplate("unittest/unittest.list.tmpt") 
	}

	if strings.Index(m.Mode, "R") > -1 && strings.Index(m.InjectionMode, "R") <= -1 {
		hasUnittest = true
		dataStruct += readTemplate("unittest/unittest.read.tmpt")
		routeUrlKey := ""
		routeUrlVal := ""
		for _, uniqueKey := range m.UniqueKey {  
			if column , ok := m.Columns[uniqueKey]; ok && !contains(m.FilterDetail, uniqueKey) {
				
				routeUrlKey += "/%%v"

				valueType := mysqlTypeToGoType(column.Type, false, true) 
				if strings.Contains(valueType, "int") {
					routeUrlVal += ",2147483647"
				}

				if valueType == "string" {
					routeUrlVal += ",\"unittestkey\""
				}  
			}
		}
		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_KEY}}", routeUrlKey, -1)
		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_VALUE}}", routeUrlVal, -1) 
	}

	

	if strings.Index(m.Mode, "U") > -1 {
		hasUnittest = true
		dataStruct += readTemplate("unittest/unittest.update.tmpt")
		routeUrlKey := ""
		routeUrlVal := ""
		for _, uniqueKey := range m.UniqueKey {  
			if column , ok := m.Columns[uniqueKey]; ok && !contains(m.FilterDetail, uniqueKey) {
				
				routeUrlKey += "/%%v"

				valueType := mysqlTypeToGoType(column.Type, false, true) 
				if strings.Contains(valueType, "int") {
					routeUrlVal += ",2147483647"
				}

				if valueType == "string" {
					routeUrlVal += ",\"unittestkey\""
				}  
			}
		}

		dataForm :=""
		for _, column := range m.Columns { 
			if !contains(m.UneditableField, column.Name) && !contains(m.UniqueKey, column.Name)  {
				valueType := mysqlTypeToGoType(column.Type, false, true) 
				if strings.Contains(valueType, "time") {
					dataForm += fmt.Sprintf("WithFormField(\"%v\", DateTime.ConvFormat({{MODULE_NAME}}CreateStruct.%v)). \n", column.Name, title(column.Name))
				} else {
					dataForm += fmt.Sprintf("WithFormField(\"%v\", {{MODULE_NAME}}CreateStruct.%v). \n", column.Name, title(column.Name))
				}
			}
		}

		dataStruct = strings.Replace(dataStruct, "{{FORM_FIELD}}", dataForm, -1)
		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_KEY}}", routeUrlKey, -1)
		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_VALUE}}", routeUrlVal, -1) 
	}

	if strings.Index(m.Mode, "D") > -1 {
		hasUnittest = true
		dataStruct += readTemplate("unittest/unittest.delete.tmpt")
		routeUrlKey := ""
		routeUrlVal := ""
		for _, uniqueKey := range m.UniqueKey {  
			if column , ok := m.Columns[uniqueKey]; ok && !contains(m.FilterDetail, uniqueKey) {
				
				routeUrlKey += "/%%v"

				valueType := mysqlTypeToGoType(column.Type, false, true) 
				if strings.Contains(valueType, "int") {
					routeUrlVal += ",2147483647"
				}

				if valueType == "string" {
					routeUrlVal += ",\"unittestkey\""
				}  
			}
		}
		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_KEY}}", routeUrlKey, -1)
		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_VALUE}}", routeUrlVal, -1) 
	}

	// if strings.Index(m.Mode, "M") > -1 {
	// 	hasUnittest = true
	// 	dataStruct += readTemplate("unittest/unittest.media.tmpt")
	// 	mediaRoute := strings.Replace(m.Name, "_media_file", "", -1) + "/media"
	// 	dataStruct = strings.Replace(dataStruct, "{{MODULE_MEDIA_ROUTE}}", mediaRoute, -1)
	// }

	
	dataStruct = strings.Replace(dataStruct, "{{MODULE_ROUTE}}", m.Name, -1)

	dataStruct += "\n"

	testStruct = strings.Replace(testStruct, "{{UNITTEST}}", dataStruct, -1)
	testStruct = strings.Replace(testStruct, "{{MODULE_NAME}}", title(m.Name), -1)


	return testStruct, hasUnittest
}
