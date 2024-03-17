package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql" // mysql
)

var tmpl *template.Template

var db *sql.DB

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/demosurveysitedb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func init() {
	tmpl = template.Must(template.ParseFiles("./template/form.html"))
}

type studentSurvey struct {
	Name            string
	Email           string
	Age             string
	Role            string
	FavoriteFeature string
	Recommendation  []string
	Comments        string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	student := studentSurvey{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Age:             r.FormValue("age"),
		Role:            r.FormValue("current_role"),
		FavoriteFeature: r.FormValue("favorite_features"),
		Recommendation:  r.Form["improvements"],
		Comments:        r.FormValue("comment"),
	}

	recommendations := ""
	if len(student.Recommendation) > 0 {
		recommendations = strings.Join(student.Recommendation, ", ")
	}

	if r.FormValue("submit") == "Submit" {
		age, _ := strconv.Atoi(student.Age)
		_, err := db.Exec("insert into student (name, email, age, current_role, favorite_features, improvements, comment) values(?, ?, ?, ?, ?, ?, ?)", student.Name, student.Email, age, student.Role, student.FavoriteFeature, recommendations, student.Comments)
		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: "Congratulations!!! Survey Complete"})
		}
	}
	fmt.Println(student)
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", formHandler)
	http.ListenAndServe(":8080", nil)
}
