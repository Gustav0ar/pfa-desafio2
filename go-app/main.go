package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

var port string
var DB *sql.DB

type Course struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func AllCourses() ([]Course, error) {
	if DB == nil {
		return nil, nil
	}

	rows, err := DB.Query("SELECT * FROM FULL_CYCLE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course

	for rows.Next() {
		var course Course

		err := rows.Scan(&course.Id, &course.Title)
		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := AllCourses()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err != nil {
		fmt.Fprintf(w, "<h3>Error while obtaining courses</h3>")
	}

	html := `<h2>Cursos</h2>
						<br/>
           <ul>`

	for _, course := range courses {
		html += fmt.Sprintf("<li>%d - %s</li>", course.Id, course.Title)
	}

	html += "</ul>"

	fmt.Fprintf(w, html)
}

func handleRequests() {
	http.HandleFunc("/", getAllCourses)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func main() {
	var err error

	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = "127.0.0.1"
	}
	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort == "" {
		mysqlPort = "3306"
	}
	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	DB, err = sql.Open("mysql", fmt.Sprintf("root:root@tcp(%s:%s)/fullcycle", mysqlHost, mysqlPort))

	if err != nil {
		panic(err.Error())
	}

	retryCount := 30
	for {
		err := DB.Ping()
		if err != nil {
			if retryCount == 0 {
				log.Fatalf("Not able to establish connection to database")
			}

			log.Printf(fmt.Sprintf("Could not connect to database. Wait 2 seconds. %d retries left...", retryCount))
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	defer DB.Close()

	handleRequests()
}
