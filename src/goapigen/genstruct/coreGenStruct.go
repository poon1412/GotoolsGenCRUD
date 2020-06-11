package genstruct

import (
	"io/ioutil"
	"strings"
)

const (
	golangByteArray  = "[]byte"
	gureguNullInt    = "null.Int"
	sqlNullInt       = "sql.NullInt64"
	golangInt        = "int"
	golangInt64      = "int64"
	gureguNullFloat  = "null.Float"
	sqlNullFloat     = "sql.NullFloat64"
	golangFloat      = "float"
	golangFloat32    = "float32"
	golangFloat64    = "float64"
	gureguNullString = "null.String"
	sqlNullString    = "sql.NullString"
	gureguNullTime   = "null.Time"
	golangTime       = "time.Time"
)

type MainStruct struct {
	Version      string            `json:"version"`
	Revision     int               `json:"revision"`
	TemplatePath string            `json:"-"`
	OutputPath   string            `json:"-"`
	Modules      map[string]Module `json:"modules,omitempty"`
}

type Module struct {
	Name               string              `json:"name"`
	Mode               string              `json:"mode"`
	InjectionMode      string     			`json:"injection_mode"`
	UnpublishField       []string            `json:"unpublish_field"`
	UneditableField      []string            `json:"uneditable_field"`
	FilterList         []string            `json:"filter_list"`
	FilterDetail       []string            `json:"filter_detail"`
	UniqueKey          []string            `json:"unique_key"`
	Validator          map[string]string   `json:"validator"`
	CacheGroupBackend  []string            `json:"cache_group_backend"`
	CacheGroupFrontend []string            `json:"cache_group_frontend"`
	CacheDuration      int                 `json:"cache_duration"`
	Plugin           map[string]map[string]string `json:"plugin"`
	Relation           map[string][]string `json:"relation"`
	Columns            map[string]Column   `json:"columns"`
}

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	NullAble string `json:"nullable"`
	Key      string `json:"key"`
	Default  string `json:"default"`
}



type Plugin struct {
	Version      string            `json:"version"`
	Name     string               `json:"name"`
	Schema   []Schema	`json:"schema"`
}
 
type Schema struct {
	Provider string `json:"provider"`
	Type string `json:"type"`
	Columns            map[string]Column   `json:"columns"`
}


func title(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	s = strings.Replace(s, " ", "", -1)
	return s
}

func readTemplate(file string) string {
	template, err := ioutil.ReadFile(TemplatePath + "/" + file)
	check(err)
	return string(template)
}

func defultStructType(mysqlType string) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint", "bigint":
		return "0"
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
		return "\"\""
	case "date", "datetime", "time", "timestamp":
		return "\"\""
	case "decimal", "double":
		return "0"
	case "float":
		return "0"
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return "\"\""
	}
	return ""
}

func convertStructType(mysqlType string) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		return golangInt
	case "bigint":
		return golangInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
		return "string"
	case "date", "datetime", "time", "timestamp":
		return golangTime
	case "decimal", "double":
		return golangFloat64
	case "float":
		return golangFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray
	}
	return ""
}

func mysqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt
	case "bigint", "timestamp":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "time":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	case "date", "datetime":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return golangTime
	case "decimal", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case "float":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray
	}
	return ""
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
