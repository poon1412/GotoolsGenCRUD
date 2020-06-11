package genstruct

import (
	"fmt"
	"strings"
)

func (m *Module) GenController() (string,string, bool, bool) {
	fmt.Println("Generate " + title(m.Name) + " Controller.")

	dataStruct := readTemplate("controllers/controller.head.tmpt")
	injectStruct := readTemplate("controllers/inject.head.tmpt")
	importAddition := ""
	importAdditionInject := ""

	hasController := false 
	hasInject := false 

	countMain := 0
	countInject := 0

	if strings.Index(importAddition, "github.com/kataras/iris") <= -1 { importAddition += "\n\"github.com/kataras/iris\"" }
	if strings.Index(importAdditionInject, "github.com/kataras/iris") <= -1 { importAdditionInject += "\n\"github.com/kataras/iris\"" }

	if strings.Index(importAddition, "server/api/models") <= -1 { importAddition += "\n\"server/api/models\"" }
	if strings.Index(importAdditionInject, "server/api/models") <= -1 { importAdditionInject += "\n\"server/api/models\"" }

	if strings.Index(m.Mode, "L") > -1 {
		hasController = true
		structure := readTemplate("controllers/controller.list.tmpt") 

		filterParamStruct := ""
		filterDefault := ""
		filterParam := ""
		allowOrder := ""

		if len(m.UniqueKey) > 0 {
			filterDefault +=  fmt.Sprintf("Obj.Param.Sort = \"%v\"\n", m.UniqueKey[0])
		}
		
		if len(m.FilterList) > 0 {
			for _, filterColName := range m.FilterList {
				if column, ok := m.Columns[filterColName]; ok {
					typeConverted := convertStructType(column.Type)
					if typeConverted == "time.Time" {
						typeConverted = "string"
					}
					filterParamStruct += fmt.Sprintf("%v %v `json:\"%v\" form:\"%v\"`\n", title(filterColName), typeConverted, filterColName, filterColName)
					filterDefault += fmt.Sprintf("Obj.Param.%v = %v\n", title(filterColName), defultStructType(column.Type))
					filterParam += fmt.Sprintf("Obj.Param.%v, ", title(filterColName))
				}
			} 
		} 
		filterDefault += "Obj.Response.Param = &Obj.Param"

		if len(m.Columns) > 0 {
			for _, column := range m.Columns { 
				if !contains(m.UnpublishField, column.Name) {
					allowOrder += fmt.Sprintf("\"%v\", ", column.Name)
				}
			}
			allowOrder = strings.TrimRight(allowOrder, ", ")
		} 

		structure = strings.Replace(structure, "{{ALLOW_ORDER}}", allowOrder, -1)
		structure = strings.Replace(structure, "{{FILTER_PARAM_STRUCT}}", filterParamStruct, -1)
		structure = strings.Replace(structure, "{{FILTER_DEFAULT}}", filterDefault, -1)
		structure = strings.Replace(structure, "{{FILTER_PARAM}}", filterParam, -1)

		cacheGroup := ""
		if len(m.CacheGroupFrontend) > 0 {
			cacheGroup = strings.Join(m.CacheGroupFrontend, ",")
			cacheGroup = fmt.Sprintf("Obj.Response.SetFrontendCache([]string{\"%v\"}, %v)", cacheGroup, m.CacheDuration)
		}
		structure = strings.Replace(structure, "{{CACHE_FRONTEND}}", cacheGroup, -1)


		dataStruct += structure

		if strings.Index(m.InjectionMode, "L") > -1 { 
			hasInject = true
			injectStruct += structure

			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "/*  ", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "*/", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "controllers.", -1)

			injectStruct = strings.Replace(injectStruct, "{{OPEN_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{END_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{CTRL_DOT}}", "controllers.", -1)
			countInject++
		} else {
			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "", -1)
			countMain++ 
		}
	}

	if strings.Index(m.Mode, "R") > -1 {
		hasController = true
		
		structure := readTemplate("controllers/controller.read.tmpt")

		
		keyInit := ""
		keyParam := ""
		keyParamStruct := ""
		filterParamStruct := ""
		filterDefault := ""
		filterParam := ""


		if len(m.UniqueKey) > 0 {
			for _, uniqueKey := range m.UniqueKey {
				if column, ok := m.Columns[uniqueKey]; ok && !contains(m.FilterDetail, uniqueKey) {
					keyParamStruct += fmt.Sprintf("%v %v `json:\"%v\" form:\"%v\"`\n", title(uniqueKey), convertStructType(column.Type), uniqueKey, uniqueKey)
					keyInit += typeToParamGet(column.Type, column.Name)
					keyInit += fmt.Sprintf("Obj.Param.%v = %v\n", title(column.Name), column.Name)

					valueType := mysqlTypeToGoType(column.Type, column.NullAble == "YES", true)
					if strings.Contains(valueType, "null.String") {
						keyParam += fmt.Sprintf("Obj.Model.%v.String = Obj.Param.%v\n", title(column.Name), title(column.Name))
					} else if strings.Contains(valueType, "null.Int") {
					 	keyParam += fmt.Sprintf("Obj.Model.%v.Int = Obj.Param.%v\n", title(column.Name), title(column.Name))
					} else {
						keyParam += fmt.Sprintf("Obj.Model.%v = Obj.Param.%v\n", title(column.Name), title(column.Name))
					}
				}
			}
		}
		
		if len(m.FilterDetail) > 0 {

			for _, filterColName := range m.FilterDetail {
				if column, ok := m.Columns[filterColName]; ok {
					filterParamStruct += fmt.Sprintf("%v %v `json:\"%v\" form:\"%v\"`\n", title(filterColName), convertStructType(column.Type), filterColName, filterColName)
					filterDefault += fmt.Sprintf("Obj.Param.%v = %v\n", title(filterColName), defultStructType(column.Type))
					valueType := mysqlTypeToGoType(column.Type, column.NullAble == "YES", true)
					if strings.Contains(valueType, "null.String") {
						filterParam += fmt.Sprintf("Obj.Model.%v.String = Obj.Param.%v\n", title(column.Name), title(column.Name))
					} else if strings.Contains(valueType, "null.Int") {
					 	filterParam += fmt.Sprintf("Obj.Model.%v.Int = Obj.Param.%v\n", title(column.Name), title(column.Name))
					} else {
						filterParam += fmt.Sprintf("Obj.Model.%v = Obj.Param.%v\n", title(column.Name), title(column.Name))
					}
				}
			}

		}

		if keyInit != "" || filterDefault != "" {
			filterDefault += "Obj.Response.Param = &Obj.Param"
		}

		structure = strings.Replace(structure, "{{KEY_INIT}}", keyInit, -1)
		structure = strings.Replace(structure, "{{KEY_PARAM}}", keyParam, -1)
		structure = strings.Replace(structure, "{{KEY_PARAM_STRUCT}}", keyParamStruct, -1)
		structure = strings.Replace(structure, "{{FILTER_PARAM_STRUCT}}", filterParamStruct, -1)
		structure = strings.Replace(structure, "{{FILTER_DEFAULT}}", filterDefault, -1)
		structure = strings.Replace(structure, "{{FILTER_PARAM}}", filterParam, -1)


		cacheGroup := ""
		if len(m.CacheGroupFrontend) > 0 {
			cacheGroup = strings.Join(m.CacheGroupFrontend, ",")
			cacheGroup = "Obj.Response.SetFrontendCache([]string{\"" + cacheGroup + "\"}, 0)"
		}
		structure = strings.Replace(structure, "{{CACHE_FRONTEND}}", cacheGroup, -1)


		dataStruct += structure

		if strings.Index(m.InjectionMode, "R") > -1 { 
			hasInject = true
			injectStruct += structure

			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "/*  ", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "*/", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "controllers.", -1)
			
			injectStruct = strings.Replace(injectStruct, "{{OPEN_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{END_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{CTRL_DOT}}", "controllers.", -1)
			countInject++
		} else {
			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "", -1)
			countMain++ 
		}
	}

	if strings.Index(m.Mode, "C") > -1 {
		hasController = true
		
		structure := readTemplate("controllers/controller.create.tmpt")

		structure = strings.Replace(structure, "{{DATA_DEFAULT}}", "", -1)

		dataStruct += structure

		if strings.Index(m.InjectionMode, "C") > -1 { 
			hasInject = true
			injectStruct += structure

			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "/*  ", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "*/", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "controllers.", -1)
			
			injectStruct = strings.Replace(injectStruct, "{{OPEN_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{END_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{CTRL_DOT}}", "controllers.", -1)
			countInject++
		} else {
			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "", -1)
			countMain++ 
		}

	}

	if strings.Index(m.Mode, "U") > -1 {
		hasController = true
		
		structure := readTemplate("controllers/controller.update.tmpt")
		
		updateKey := ""
		updateInit := ""

		// if len(m.UniqueKey) > 0 {

		// 	for _, uniqueKey := range m.UniqueKey {
		// 		if column, ok := m.Columns[uniqueKey]; ok {
		// 			updateKey += typeToParamGet(column.Type, column.Name)
		// 			updateKey += fmt.Sprintf("Obj.Param.%v = %v\n", title(column.Name), column.Name)
		// 			updateInit += fmt.Sprintf("Obj.Model.%v = Obj.Param.%v\n", title(column.Name), title(column.Name))
					
		// 		}
		// 	}
		// }

		for _, column := range m.Columns { 
			if column.Key == "PRI" {
					updateKey += typeToParamGet(column.Type, column.Name)
					updateKey += fmt.Sprintf("Obj.Param.%v = %v\n", title(column.Name), column.Name)
					updateInit += fmt.Sprintf("Obj.Model.%v = Obj.Param.%v\n", title(column.Name), title(column.Name))
			}
		}
		
		structure = strings.Replace(structure, "{{UPDATE_KEY}}", updateKey, -1)
		structure = strings.Replace(structure, "{{UPDATE_KEY_INIT}}", updateInit, -1)

		dataStruct += structure

		if strings.Index(m.InjectionMode, "U") > -1 { 
			hasInject = true
			injectStruct += structure

			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "/*  ", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "*/", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "controllers.", -1)
			
			injectStruct = strings.Replace(injectStruct, "{{OPEN_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{END_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{CTRL_DOT}}", "controllers.", -1)
			countInject++
		} else {
			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "", -1)
			countMain++ 
		}
	}

	if strings.Index(m.Mode, "D") > -1 {
		hasController = true
		
		structure := readTemplate("controllers/controller.delete.tmpt")

		keyInit := ""
		keyParam := ""
		if len(m.UniqueKey) > 0 {

			for _, uniqueKey := range m.UniqueKey {
				if column, ok := m.Columns[uniqueKey]; ok {
					keyInit += typeToParamGet(column.Type, column.Name)
					keyParam += fmt.Sprintf("Obj.Model.%v = %v\n", title(column.Name), column.Name)
				}
			}

		}

		structure = strings.Replace(structure, "{{KEY_INIT}}", keyInit, -1)
		structure = strings.Replace(structure, "{{KEY_PARAM}}", keyParam, -1)

		dataStruct += structure

		if strings.Index(m.InjectionMode, "D") > -1 { 
			hasInject = true
			injectStruct += structure

			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "/*  ", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "*/", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "controllers.", -1)
			
			injectStruct = strings.Replace(injectStruct, "{{OPEN_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{END_INJECT}}", "", -1)
			injectStruct = strings.Replace(injectStruct, "{{CTRL_DOT}}", "controllers.", -1)
			countInject++
		} else {
			dataStruct = strings.Replace(dataStruct, "{{OPEN_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{END_INJECT}}", "", -1)
			dataStruct = strings.Replace(dataStruct, "{{CTRL_DOT}}", "", -1)
			countMain++ 
		}
	}

	if strings.Index(m.Mode, "M") > -1 {
		hasController = true
		countMain++ 
		importAddition += "\n\"strconv\""
		importAddition += "\n\"time\""
		if strings.Index(importAddition, "server/api/models") <= -1 { importAddition += "\n\"server/api/models\"" }
		dataStruct += readTemplate("controllers/controller.media.tmpt")
	}

	if countMain <= 0 {
		importAddition = ""
	}

	dataStruct = strings.Replace(dataStruct, "{{MODULE_NAME}}", title(m.Name), -1)
	dataStruct = strings.Replace(dataStruct, "{{IMPORT_ADDITION}}", importAddition, -1)


	injectStruct = strings.Replace(injectStruct, "{{MODULE_NAME}}", title(m.Name), -1)
	injectStruct = strings.Replace(injectStruct, "{{IMPORT_ADDITION}}", importAdditionInject, -1)
	
 

	return dataStruct, injectStruct, hasController, hasInject
}

func typeToParamGet(mysqlType string, input string) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		return fmt.Sprintf("%v , _ := ctx.Params().GetInt(\"%v\")\n", input, input)
	case "bigint":
		return fmt.Sprintf("%v , _ := ctx.Params().GetInt64(\"%v\")\n", input, input)
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
		return fmt.Sprintf("%v := ctx.Params().GetTrim(\"%v\")\n", input, input)
	case "decimal", "double":
		return fmt.Sprintf("%v , _ := ctx.Params().GetFloat64(\"%v\")\n", input, input)
	case "float":
		return fmt.Sprintf("%v , _ := ctx.Params().GetFloat32(\"%v\")\n", input, input)
	}
	return ""
}
