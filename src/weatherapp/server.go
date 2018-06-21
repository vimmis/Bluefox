/*
	Weatherapp with basic controllers (v1)
*/

package main

import (
	"os"
	"time"
	"log"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/unrolled/render"	
	"net/http"
	"github.com/gorilla/mux" //URL dispatcher and router 
	"github.com/codegangsta/negroni" //HTTP Middleware
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Env Variables
var usr = os.Getenv("MYSQL_USR")
var pwd = os.Getenv("MYSQL_PWD")
var dbname = os.Getenv("MYSQL_DB")
var dbtbl = os.Getenv("MYSQL_TBL")
var mysql_connect = usr+":"+pwd+"@tcp(mysql:3306)/"

// This method helps in configuration of new server and returns it.
func NewServer() *negroni.Negroni {
	formating := render.New(render.Options{
		IndentJSON: true,
	})
	muxr := mux.NewRouter()
	ng := negroni.Classic()
	ng.UseHandler(muxr)
	initRoutes(muxr, formating)
	return ng
}

// Initialize MysqlDB
func init() {
	db, err := sql.Open("mysql", mysql_connect)
	if err != nil {
		log.Fatal("Cant connect to DB: ",err)
	} else {
		defer db.Close()
		 _,err = db.Exec("CREATE DATABASE IF NOT EXISTS "+dbname)
		if err != nil {
		   panic(err)
		}
		_,err = db.Exec("USE "+dbname)
		if err != nil {
			panic(err)
		}		
		_,err = db.Exec("CREATE TABLE IF NOT EXISTS "+dbtbl+" ( id MEDIUMINT NOT NULL AUTO_INCREMENT PRIMARY KEY, Timestamp TIMESTAMP NOT NULL UNIQUE, Tempt FLOAT NOT NULL)")
		if err != nil {
			panic(err)
		}
	}

}

// APIs Routing
func initRoutes(muxr *mux.Router, formating *render.Render) {
	muxr.HandleFunc("/pingLB", pingLBController(formating)).Methods("GET") //Helps in Load balancer(if any) for health check of a server 
	muxr.HandleFunc("/add_temperature_measurement", addTempController(formating)).Methods("POST")
	muxr.HandleFunc("/get_temperature_measurements", getTempController(formating)).Methods("GET")
	muxr.HandleFunc("/get_average_temperature", getAvgTempController(formating)).Methods("GET")

}

//LB Ping
func pingLBController(formating *render.Render) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		formating.JSON(writer, http.StatusOK, struct{ Ping string }{"Ping Version 1!"})
	}
}

// Add Temperature Measurement
func addTempController(formating *render.Render) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		
		var (
			id int
			time_var time.Time
			tempt_var	float64
			requestParams WeatherDataAdd
		)
		
		//Decode Request Body and fetch data
    	_ = json.NewDecoder(req.Body).Decode(&requestParams)		
		var timestamp_string string = requestParams.Timestamp
		var temp_string string = requestParams.Tempt
    	fmt.Println("Add Data- Received", timestamp_string , temp_string)
		
		//Convert string-timestamp to timestamp		
		timestmp, err1 := time.Parse(time.RFC3339, timestamp_string)
		if err1 != nil {
			fmt.Println(err1)
			formating.JSON(writer, http.StatusBadRequest, "Null or Invalid timestamp, expected: 2006-01-02T15:04:05.000Z  format")
			return
		}
		fmt.Println("Converted timestamp: ",timestmp)
		
		//Convert String-temperature to float
		tempt, err2 := strconv.ParseFloat(temp_string, 64)
		if err2 != nil {
			fmt.Println(err2)
			formating.JSON(writer, http.StatusBadRequest, "Null or Invalid temperature or out-of-bounds")
			return
		}
		fmt.Println("Converted temperature: ",tempt)
		
		//Open Mysql
		db, err3 := sql.Open("mysql", mysql_connect+dbname+"?parseTime=true")
		if err3 != nil {
			log.Fatal(err3)
			formating.JSON(writer, http.StatusFailedDependency, "Unable to connect to DB")
			return
		}
		defer db.Close()
		query_stmt := "INSERT INTO "+dbtbl+"(Timestamp, Tempt) VALUES(?, ?)"
		fmt.Println("Adding Query:",query_stmt)
		stmt, err4 := db.Prepare(query_stmt)
		if err4 != nil {
			log.Fatal(err4)
			formating.JSON(writer, http.StatusFailedDependency, "Unable to prepare statement to store data")
			return
		}
		res, err5 := stmt.Exec(timestmp , tempt)
		if err5 != nil {
			log.Fatal(err5)
			formating.JSON(writer, http.StatusFailedDependency, "Unable to store data")
			return
		}
		lastId, err6 := res.LastInsertId()
		if err6 != nil {
			log.Fatal(err6)
			formating.JSON(writer, http.StatusFailedDependency, "Unable to read stored data last ID")
			return
		}
		fmt.Println("MYSQL insert success with ID:", lastId)
		
		//Fetch data again from database for verification
		query_stmt = "select id, Timestamp, Tempt from "+dbtbl+" where id = ?"
		fmt.Println("Fetching Query:",query_stmt)
		
		err7 := db.QueryRow(query_stmt, lastId).Scan(&id, &time_var, &tempt_var)//Expected single row
		if err7 != nil {
			log.Fatal("Fetching failed", err7)
			formating.JSON(writer, http.StatusFailedDependency, "Unable to fetch stored data")
			return
		}
		log.Println(" Added data to mysql: ",id, time_var, tempt_var)
		result := WeatherDataGet {
			Timestamp : time_var,
			Tempt : tempt_var,
		}
		fmt.Println("Weathr Data Ends")
		formating.JSON(writer, http.StatusOK, result)
		
	}
}

//Get Temperature Measurements
func getTempController(formating *render.Render) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
	
		var (
			id int
			time_var time.Time
			tempt_var	float64
			start_string string
			stop_string string
			measurements_array [] WeatherDataGet
		)
		
		param1, err1 := req.URL.Query()["start_timestamp"]    
		if !err1 || len(param1) < 1 {
			log.Println("Url Param 'start_timestamp' is missing", err1)
			formating.JSON(writer, http.StatusBadRequest, "start time missing")
			return
		}
		param2, err2 := req.URL.Query()["stop_timestamp"]    
		if !err2 || len(param2) < 1 {
			log.Println("Url Param 'stop_timestamp' is missing", err2)
			formating.JSON(writer, http.StatusBadRequest, "stop_timestamp missing")
			return
		}
		
		start_string = param1[0]
		stop_string = param2[0]
		fmt.Println("Recieved params: ", start_string,stop_string )
				
		//Convert string timestamp to timestamp		
		timestmp_start, err3 := time.Parse(time.RFC3339, start_string)
		if err3 != nil {
			fmt.Println(err3)
			formating.JSON(writer, http.StatusBadRequest, "Invalid start timestamp, expected: 2006-01-02T15:04:05.000Z  format")
			return
		}
		timestmp_stop, err4 := time.Parse(time.RFC3339, stop_string)
		if err4 != nil {
			fmt.Println(err4)
			formating.JSON(writer, http.StatusBadRequest, "Invalid stop timestamp, expected: 2006-01-02T15:04:05.000Z  format")
			return
		}
		fmt.Println("Converted timestamps: ",timestmp_start,timestmp_stop )
		
		//Open Mysql
		db, err5 := sql.Open("mysql", mysql_connect+dbname+"?parseTime=true")
		if err5 != nil {
			log.Fatal("Connection failed", err5)
			formating.JSON(writer, http.StatusFailedDependency, "Unable to connect mysql")
			return
		}
		defer db.Close()
		
		//Fetch data from ranges(inclusive)
		query_stmt := "select id, Timestamp, Tempt from "+dbtbl+" where Timestamp between ? AND ?"
		fmt.Println("Fetching Query:",query_stmt)
		rows, _ := db.Query(query_stmt, timestmp_start, timestmp_stop)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&id, &time_var, &tempt_var)
			log.Println("Fetching: ",id, time_var, tempt_var)
			temp:= WeatherDataGet {
					Timestamp : time_var,
					Tempt : tempt_var,
			}
			measurements_array = append(measurements_array, temp)
		}
		fmt.Println("Weather Data Ends")
		if len(measurements_array) > 0 {
			formating.JSON(writer, http.StatusOK, measurements_array)
		}else {
			formating.JSON(writer, http.StatusOK,"null")
		}
	}
}
 
//Get Temperature Average
func getAvgTempController(formating *render.Render) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
	
		var (
			id int
			time_var time.Time
			tempt_var	float64
			start_string string
			stop_string string
		)
		
		param1, err1 := req.URL.Query()["start_timestamp"]    
		if !err1 || len(param1) < 1 {
			log.Println("Url Param 'start_timestamp' is missing", err1)
			formating.JSON(writer, http.StatusBadRequest, "start time missing")
			return
		}
		
		param2, err2 := req.URL.Query()["stop_timestamp"]    
		if !err2 || len(param2) < 1 {
			log.Println("Url Param 'stop_timestamp' is missing", err2)
			formating.JSON(writer, http.StatusBadRequest, "stop_timestamp missing")
			return
		}
		start_string = param1[0]
		stop_string = param2[0]
		fmt.Println("Recieved params: ", start_string,stop_string )
		
		//Convert string timestamp to timestamp		
		timestmp_start, err3 := time.Parse(time.RFC3339, start_string)
		if err3 != nil {
			fmt.Println(err3)
			formating.JSON(writer, http.StatusBadRequest, "Invalid start_timestamp, expected: 2006-01-02T15:04:05.000Z  format")
			return
		}
		timestmp_stop, err4 := time.Parse(time.RFC3339, stop_string)
		if err4 != nil {
			fmt.Println(err4)
			formating.JSON(writer, http.StatusBadRequest, "Invalid stop_timestamp, expected: 2006-01-02T15:04:05.000Z  format")
			return
		}
		fmt.Println("Converted timestamps: ",timestmp_start,timestmp_stop )
		
		//Open Mysql
		db, err5 := sql.Open("mysql", mysql_connect+dbname+"?parseTime=true")
		if err5 != nil {
			log.Fatal("Connection failed", err5)
			formating.JSON(writer, http.StatusFailedDependency, "Unable to connect mysql")
			return
		}
		
		//Fetch data for the ranges(inclusive)
		defer db.Close()
		query_stmt := "select id, Timestamp, Tempt from "+dbtbl+" where Timestamp between ? AND ?"
		fmt.Println("Fetching Query:",query_stmt)
		rows, _ := db.Query(query_stmt, timestmp_start, timestmp_stop)
		defer rows.Close()
		
		count_rows:=0.0 //counts rows fetched
		sum_temp:=0.0 //sums temperatures fetched
		for rows.Next() {
			rows.Scan(&id, &time_var, &tempt_var)
			log.Println("Fetching: ",id, time_var, tempt_var)
			count_rows= count_rows+1
			sum_temp = sum_temp+tempt_var			
		}
		avg:=0.0
		if count_rows > 0 {
			avg=sum_temp/count_rows
		}
		fmt.Println("Average:",avg)
		fmt.Println("Weather Data Ends")
		formating.JSON(writer, http.StatusOK, avg)
	}
}
