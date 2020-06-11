package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"goapigen/genjson"
	"goapigen/genstruct"
	"os"
)

//Program to reverse engineer your mysql database into gorm models
func main() {
	godotenv.Load(os.Getenv("PWD") + "/.env")
	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOSTNAME")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	mode := flag.String("mode", "code", "Generate Mode")
	templatePath := flag.String("template", os.Getenv("PWD")+"/template", "Template Path")
	outputPath := flag.String("project", os.Getenv("PWD")+"/src", "Project Path")
	relationMode := flag.Bool("relation", false, "Relation Mode")
	flag.Parse()

	fmt.Println("Generate Mode: " + *mode)
	fmt.Println("Template Path: " + *templatePath)
	fmt.Println("Project Path: " + *outputPath)

	if *mode == "json" || *mode == "all" {
		fmt.Println("Generate:", "goapigen.json")
		fmt.Println("hostdb:", host)
		fmt.Println("Connecting to mysql server " + host + ":" + port)
		genjson.GenJson(user, pass, host, port, database, *relationMode)
	}

	if *mode == "code" || *mode == "all" {
		genstruct.GenStruct(*templatePath, *outputPath)
	}

	return
}
