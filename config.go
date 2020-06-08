package main

import (
	"os"
	"strconv"
	"strings"
)

// Config ...
type Config struct {
	Http struct {
		Port int
	}
	/*==============================================*/
	Jwt struct {
		Private string
		Public  string
	}
	/*==============================================*/
	Mongodb struct {
		Host     string
		Port     string
		Database string
		Username string
		Password string
	}
	/*==============================================*/
	Mssql struct {
		Host     string
		Database string
		Username string
		Password string
		Port     int
	}
	Mysql struct {
		Host     string
		Database string
		Username string
		Password string
		Port     string
	}
	Redisdb struct {
		Addr     string
		Password string
		DB       int
	}
}

func (c *Config) Read() {
	// receive env from kubernetes
	// port
	if val, ok := os.LookupEnv("HTTP_PORT"); ok {
		if Num, err := strconv.Atoi(val); err == nil {
			c.Http.Port = Num
		}
	}
	// mssql
	if val, ok := os.LookupEnv("MSSQL_HOST"); ok {
		c.Mssql.Host = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MSSQL_DB"); ok {
		c.Mssql.Database = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MSSQL_USERNAME"); ok {
		c.Mssql.Username = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MSSQL_PASSWORD"); ok {
		c.Mssql.Password = strings.TrimSpace(val)
	}
	// mongo
	if val, ok := os.LookupEnv("MONGO_HOST"); ok {
		c.Mongodb.Host = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MONGO_PORT"); ok {
		c.Mongodb.Port = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MONGO_DATABASE"); ok {
		c.Mongodb.Database = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MONGO_USERNAME"); ok {
		c.Mongodb.Username = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MONGO_PASSWORD"); ok {
		c.Mongodb.Password = strings.TrimSpace(val)
	}
	// mysql
	if val, ok := os.LookupEnv("MYSQL_HOST"); ok {
		c.Mysql.Host = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MYSQL_DB"); ok {
		c.Mysql.Database = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MYSQL_USERNAME"); ok {
		c.Mysql.Username = strings.TrimSpace(val)
	}
	if val, ok := os.LookupEnv("MYSQL_PASSWORD"); ok {
		c.Mysql.Password = strings.TrimSpace(val)
	}
	return
}
