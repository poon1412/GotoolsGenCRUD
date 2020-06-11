package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"strconv" 
	"godb2struc/Shelnutt2/db2struct"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)



//Program to reverse engineer your mysql database into gorm models
func main() {

	packagename := "datamodels"
	user := flag.String("user", "username", "Username") 
	pass := flag.String("pass", "pass", "Password")  
	host := flag.String("host", "localhostss", "Host")   
	database := flag.String("db", "dbname", "Database Name")  
	port := flag.Int("port", 3306, "Database Name")  
	output := flag.String("o", os.Getenv("GOPATH") + "/src/api/" + packagename, "Output Path")  
	flag.Parse() 
	
	fmt.Println("hostdb:", *host)
	fmt.Println("Connecting to mysql server " + *host + ":" + strconv.Itoa(*port)) 
	db, err := gorm.Open("mysql", *user+":"+*pass+"@tcp("+ *host +":"+strconv.Itoa(*port)+")/"+ *database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln("Failed to connect database")
	}
	defer db.Close()
	//Get all the tables from Database
	rows, err := db.Raw("SHOW TABLES").Rows()
	defer rows.Close()
	for rows.Next() {
		var table string
		rows.Scan(&table)
		columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(*user, *pass, *host, *port, *database, table)
		if err != nil {
			fmt.Println("Error in selecting column data information from mysql information schema")
			return
		}
		// Generate struct string based on columnDataTypes
		struc, err := db2struct.Generate(*columnDataTypes, table, table, packagename, true, true, true)
		if err != nil {
			fmt.Println("Error in creating struct from json: " + err.Error())
			return
		}
		file, err := os.Create(*output + "/" + table + ".go")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()
		fmt.Fprintf(file, string(struc))
		log.Println("Wrote " + table + ".go to disk")
	}

}