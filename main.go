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

type studentSurvey struct {
	Name            string
	Email           string
	Age             string
	Role            string
	FavoriteFeature string
	Recommendation  []string
	Comments        string
}

type studentNewsletter struct {
	Email string
}

type contactUs struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/demosurveysitedb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func init() {
	tmpl = template.Must(template.ParseGlob("./template/*.html"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
	return
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "contact.html", nil)
		return
	}
	userMessage := contactUs{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}
	if r.FormValue("submit") == "Send" {
		_, err := db.Exec("INSERT INTO messages (name, email, subject, message) VALUES(?, ?, ?, ?)", userMessage.Name, userMessage.Email, userMessage.Subject, userMessage.Message)
		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: "Message Sent Succesfully"})
		}
	}

}

func formHandler(w http.ResponseWriter, r *http.Request) {
	db = getMySQLDB()
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "form.html", nil)
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
			return
		} else {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: "Congratulations!!! Survey Complete"})
			return
		}
	}
	fmt.Println(student)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	db := getMySQLDB()
	if r.Method == http.MethodPost {
		newsletter := studentNewsletter{
			Email: r.FormValue("email"),
		}
		_, err := db.Exec("insert into newsletter (email) values(?)", newsletter.Email)
		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: "Failed to subscribe: " + err.Error()})
		}
		tmpl.Execute(w, struct {
			Success bool
			Message string
		}{Success: true, Message: "You have Subscribed for our routine Newsletters"})
	}
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/survey", formHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.ListenAndServe(":8080", nil)
}
