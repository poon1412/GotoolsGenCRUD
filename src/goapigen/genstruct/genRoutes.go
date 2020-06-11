package genstruct

import (
	"fmt"
	"strings"
)

func (m *Module) GenRoute() (string ,bool ,bool) {
	fmt.Println("Generate " + title(m.Name) + " Route.")

	dataStruct := "/* Route {{MODULE_NAME}} */\n"
	hasInject := false
	hasController := false
	if strings.Index(m.Mode, "L") > -1 {
		dataStruct += readTemplate("routes/route.list.tmpt")
		if strings.Index(m.InjectionMode, "L") > -1 {  
			hasInject = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "injection", -1)
		} else {
			hasController = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "controllers", -1)
		}
	}

	if strings.Index(m.Mode, "R") > -1 {

		dataStruct += readTemplate("routes/route.read.tmpt")
		routeUrlKey := ""
		for _, uniqueKey := range m.UniqueKey {
			if _, ok := m.Columns[uniqueKey]; ok && !contains(m.FilterDetail, uniqueKey) {
				routeUrlKey += fmt.Sprintf("/{%v}", uniqueKey)
			}
		}
		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_KEY}}", routeUrlKey, -1)

		if strings.Index(m.InjectionMode, "R") > -1 {  
			hasInject = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "injection", -1)
		} else {
			hasController = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "controllers", -1)
		}
	}

	if strings.Index(m.Mode, "C") > -1 {
		dataStruct += readTemplate("routes/route.create.tmpt")

		if strings.Index(m.InjectionMode, "C") > -1 {  
			hasInject = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "injection", -1)
		} else {
			hasController = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "controllers", -1)
		}
	}

	if strings.Index(m.Mode, "U") > -1 {
		dataStruct += readTemplate("routes/route.update.tmpt")
		routeUrlKey := ""
		// for _, uniqueKey := range m.UniqueKey {
		// 	if _, ok := m.Columns[uniqueKey]; ok {
		// 		routeUrlKey += fmt.Sprintf("/{%v}", uniqueKey)
		// 	}
		// }

		for _, column := range m.Columns { 
			if column.Key == "PRI" {
				routeUrlKey += fmt.Sprintf("/{%v}", column.Name)
			}
		}

		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_KEY}}", routeUrlKey, -1)

		if strings.Index(m.InjectionMode, "U") > -1 {  
			hasInject = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "injection", -1)
		} else {
			hasController = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "controllers", -1)
		}
	}

	if strings.Index(m.Mode, "D") > -1 {
		dataStruct += readTemplate("routes/route.delete.tmpt")
		routeUrlKey := ""
		for _, uniqueKey := range m.UniqueKey {
			if _, ok := m.Columns[uniqueKey]; ok {
				routeUrlKey += fmt.Sprintf("/{%v}", uniqueKey)
			}
		}
		dataStruct = strings.Replace(dataStruct, "{{ROUTE_URL_KEY}}", routeUrlKey, -1)

		if strings.Index(m.InjectionMode, "D") > -1 {  
			hasInject = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "injection", -1)
		} else {
			hasController = true
			dataStruct = strings.Replace(dataStruct, "{{MODE}}", "controllers", -1)
		}
	}

	if strings.Index(m.Mode, "M") > -1 {
		hasController = true
		dataStruct += readTemplate("routes/route.media.tmpt")
		mediaRoute := strings.Replace(m.Name, "_media_file", "", -1) + "/media"
		dataStruct = strings.Replace(dataStruct, "{{MODULE_MEDIA_ROUTE}}", mediaRoute, -1)
	}

	dataStruct = strings.Replace(dataStruct, "{{MODULE_NAME}}", title(m.Name), -1)
	dataStruct = strings.Replace(dataStruct, "{{MODULE_ROUTE}}", m.Name, -1)

	dataStruct += "\n"

	return dataStruct, hasController, hasInject
}
