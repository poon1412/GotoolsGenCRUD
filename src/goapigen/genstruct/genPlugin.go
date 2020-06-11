package genstruct

import (
	"fmt"
	"strings"
)

func (m *Module) GenPlugin() (string,string,string) {  

	modelPlugin := ""
	routePlugin  := ""
	routePluginImport := ""
	for pluginName , _ := range m.Plugin {

		// Controller
		pluginStruc := readPluginStruct("plugin/"+pluginName+"/_plugin.json")
		fmt.Println(pluginStruc)
		dataStruct := ""
		dataStruct += readTemplate("plugin/"+pluginName+"/template/controller.tmpt")

		keyInit := ""
		keyParam := ""
		keyParamStruct := ""

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

		dataStruct = strings.Replace(dataStruct, "{{KEY_INIT}}", keyInit, -1)
		dataStruct = strings.Replace(dataStruct, "{{KEY_PARAM}}", keyParam, -1)
		dataStruct = strings.Replace(dataStruct, "{{KEY_PARAM_STRUCT}}", keyParamStruct, -1)
		dataStruct = strings.Replace(dataStruct, "{{MODULE_NAME}}", title(m.Name), -1)
		writeFileWithMakeDir(dataStruct, OutputPath+"/server/api/controllers/plugin/"+pluginName, m.Name+".go")

		


		// Model
		dataStruct = ""
		dataStruct += readTemplate("plugin/"+pluginName+"/template/model.tmpt") 
		dataStruct = strings.Replace(dataStruct, "{{MODULE_NAME}}", title(m.Name), -1)

	 
		filterInit := ""	
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
		dataStruct = strings.Replace(dataStruct, "{{FILTER_INIT}}", filterInit, -1) 


		for _ , schema := range pluginStruc.Schema {
			for keyCol  , pluginCol  := range schema.Columns { 
				colName := strings.Replace(pluginCol.Name, "{m}" , m.Name , -1 ) 
				dataStruct = strings.Replace(dataStruct, "{{COLUMN_" +title(keyCol)+ "}}", title(colName), -1)
			}
		}

		modelPlugin += dataStruct+"\n\n"


		routePlugin += readTemplate("plugin/"+pluginName+"/template/route.tmpt")

		routeUniqKey := ""
		for _, uniqueKey := range m.UniqueKey {
			if _, ok := m.Columns[uniqueKey]; ok {
				routeUniqKey += fmt.Sprintf("/{%v}", uniqueKey)
			}
		} 
		routePlugin = strings.Replace(routePlugin, "{{ROUTE_UNIQ_KEY}}", routeUniqKey, -1) 
		routePlugin = strings.Replace(routePlugin, "{{MODULE_NAME}}", title(m.Name), -1) 
		routePlugin = strings.Replace(routePlugin, "{{MODULE_ROUTE}}", m.Name, -1) 

		routePluginImport += "\n\"server/api/controllers/plugin/" + pluginName + "\""

		// writeFileWithMakeDir(dataStruct, OutputPath+"/server/api/models/plugin/"+pluginName, m.Name+".go")

	}  

	return modelPlugin, routePlugin , routePluginImport

}