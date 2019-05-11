package handlers

import (
	"fmt"
	"log"
	"net/http"

	helpers "../helpers"
	repos "../repos"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Handlers

// GET calls
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginPageHandler")
	var body, _ = helpers.LoadFile("templates/login.html")
	fmt.Fprintf(w, body)
}

// POST calls
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler")
	name := r.FormValue("name")
	password := r.FormValue("password")
	redirectTarget := "/"
	log.Printf("Value of name is %s and password is %s", name, password)
	if !helpers.IsEmpty(name) && !helpers.IsEmpty(password) {
		// Perform a login check with the db
		_userIsValid := repos.UserIsValid(name, password)

		if _userIsValid {
			SetCookie(name, w)
			redirectTarget = "/index"
		} else {
			redirectTarget = "/register"
		}
	}
	http.Redirect(w, r, redirectTarget, 302)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterPageHandler")
	var body, _ = helpers.LoadFile("templates/register.html")
	fmt.Fprintf(w, body)
}

// POST
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterHandler")
	r.ParseForm()

	uName := r.FormValue("username")
	email := r.FormValue("email")
	pwd := r.FormValue("password")
	confirmPwd := r.FormValue("confirmPassword")

	_uName, _email, _pwd, _confirmPwd := false, false, false, false
	_uName = helpers.IsEmpty(uName)
	_email = helpers.IsEmpty(email)
	_pwd = helpers.IsEmpty(pwd)
	_confirmPwd = helpers.IsEmpty(confirmPwd)

	if _uName || _email || _pwd || _confirmPwd {
		fmt.Fprintf(w, "\nOne of Field is coming null")
	} else {
		fmt.Fprintf(w, "\nUsername for Register : ", uName)
		fmt.Fprintf(w, "\nEmail for Register : ", email)
		fmt.Fprintf(w, "\nPassword for Register : ", pwd)
		fmt.Fprintf(w, "\nConfirm Password for Register : ", confirmPwd)
	}
}

// GET
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("IndexPageHandler")
	userName := GetUserName(r)
	if !helpers.IsEmpty(userName) {
		var indexBody, _ = helpers.LoadFile("templates/index.html")
		fmt.Fprintf(w, indexBody, userName)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// POST
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LogoutHandler")
	ClearCookie(w)
	http.Redirect(w, r, "/", 302)
}

// Cookie

func SetCookie(userName string, w http.ResponseWriter) {
	log.Println("SetCookie")
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("cookie", value); err == nil {
		cookie := &http.Cookie{
			Name:  "cookie",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func ClearCookie(w http.ResponseWriter) {
	log.Println("ClearCookie")
	cookie := &http.Cookie{
		Name:   "cookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func GetUserName(r *http.Request) (userName string) {
	log.Println("GetUserName")
	if cookie, err := r.Cookie("cookie"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			//
			userName = cookieValue["name"]
		}
	}
	return userName
}
