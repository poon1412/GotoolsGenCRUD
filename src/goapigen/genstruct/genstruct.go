package genstruct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec" 
	"strings"
)

var TemplatePath string
var OutputPath string

func GenStruct(templatePath string, outputPath string) {

	TemplatePath = templatePath
	OutputPath = outputPath
	// jsonBlob, err := ioutil.ReadFile(os.Getenv("PWD") + "/goapigen.json")
	// check(err)

	mainStruct := MainStruct{}
	// err = json.Unmarshal(jsonBlob, &mainStruct)
	// check(err)

	if _, err := os.Stat(os.Getenv("PWD") + "/blueprint/_structure.json"); os.IsNotExist(err) {
		 fmt.Println("Error: BluePrint Not found")
		 return 
	} else {
		mainStruct = readBlueprint(os.Getenv("PWD") + "/blueprint")
	}



	fmt.Println("\n\n\nVersion: " + mainStruct.Version)
	fmt.Println(fmt.Sprintf("Revision: %d", mainStruct.Revision))

	routeStruct := ""
	routeImport := ""
 

	for moduleName, module := range mainStruct.Modules {

		routeStr, routeController, routeInject := module.GenRoute()
		routeStruct += routeStr
		if routeController {
			if strings.Index(routeImport, "server/api/controllers") <= -1 { routeImport += "\n\"server/api/controllers\"" }
		}

		if routeInject {
			if strings.Index(routeImport, "server/api/controllers/injection") <= -1 { routeImport += "\n\"server/api/controllers/injection\"" }
		}
		
		if fileStr, fileStrInject ,  hasController, hasInject := module.GenController(); hasController {
			writeFile(fileStr, outputPath+"/server/api/controllers/"+moduleName+".go")	
			if hasInject {
				writeFileIfNotExist(fileStrInject, outputPath+"/server/api/controllers/injection/"+moduleName+".go")	
			}
		} else {
			deleteFile(outputPath+"/server/api/controllers/"+moduleName+".go")
			if hasInject {
				writeFileIfNotExist(fileStrInject, outputPath+"/server/api/controllers/injection/"+moduleName+".go")	
			}
		}

		modelStr := module.GenModel() 
		modelPlugin, routePlugin, routePluginImport := module.GenPlugin()
		modelStr += modelPlugin
		writeFile(modelStr, outputPath+"/server/api/models/"+moduleName+".go")	
	 	
		routeStruct += routePlugin
		routeImport += routePluginImport
	
		

		// if fileStr, filePath, ok := module.GenPlugin(); ok {
		// 	writeFile(fileStr, outputPath+"/server/api/plugin/"+filePath+".go")	
		// } else {
		// 	deleteFile(outputPath+"/server/api/plugin/"+filePath+".go")
		// }

		// if fileStr, ok := module.GenUnitTest(); ok {
		// 	writeFile(fileStr, outputPath+"/server/test/"+moduleName+"_test.go")	
		// } else {
		// 	deleteFile(outputPath+"/server/test/"+moduleName+"_test.go")
		// }
		
		// writeFile(module.GenMockingModels(), outputPath+"/mock/mockmodels/"+moduleName+".go")
	}

	routeTemplate := readTemplate("routes/route.tmpt")
	routeTemplate = strings.Replace(routeTemplate, "{{ROUTE}}", routeStruct, -1)
	routeTemplate = strings.Replace(routeTemplate, "{{IMPORT_ADDITION}}", routeImport, -1)
	writeFile(routeTemplate, outputPath+"/server/routes/routes.go")
 

	cmd := exec.Command("gofmt", "-l", "-s","-w", outputPath + "/server/api")
	out, err := cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(string(out))

}

func readPluginStruct(filePath string) Plugin {

	struc := Plugin{}
	jsonBlob, err := ioutil.ReadFile(TemplatePath + "/" +filePath )
	check(err)
	
	err = json.Unmarshal(jsonBlob, &struc)
	check(err)

	return struc
}

func readBlueprint(blueprintPath string) MainStruct {
	
	struc := MainStruct{}
	jsonBlob, err := ioutil.ReadFile(blueprintPath + "/_structure.json")
	check(err)
	
	err = json.Unmarshal(jsonBlob, &struc)
	check(err)

	files, errd := ioutil.ReadDir(blueprintPath)
    check(errd)

    struc.Modules = make(map[string]Module)

    for _, f := range files { 
    	if f.Name() != "_structure.json" {
	    	module := Module{}
	    	fmt.Println("ReadBlueprint: "+ f.Name())
	    	jsonModule, err := ioutil.ReadFile(blueprintPath + "/"+f.Name())
	    	check(err) 
	        err = json.Unmarshal(jsonModule, &module)
	        check(err)
	        struc.Modules[module.Name] = module
	    }
    }

    return struc
}

func writeFile(data string, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	fmt.Fprintf(file, data)
}


func makeDir(path string){
	err := os.MkdirAll(path,0755)
	check(err) 
}


func writeFileWithMakeDir(data string, path string, fileName string) {
	makeDir(path)
	writeFile(data, path + "/" + fileName) 
}

func writeFileIfNotExist(data string, filePath string) {
	var _, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		writeFile(data, filePath)	
	}
}

func deleteFile( filePath string) {
	// delete file
	_ = os.Remove(filePath)
	
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
