package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hculpan/lordsofcrime/entity"
	"github.com/hculpan/lordsofcrime/templates"
	"github.com/hculpan/lordsofcrime/utils"
)

func createAccountHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Create account requested: %s -> %s\n", req.Method, req.URL.Path)

	switch req.Method {
	case "GET":
		if err := executeTemplate("create_account.gohtml", nil, w, req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		p1 := req.FormValue("inputPassword")
		p2 := req.FormValue("confirmPassword")
		username := req.FormValue("inputEmail")
		fullname := req.FormValue("inputFullName")
		displayname := req.FormValue("inputDisplayName")

		switch {
		case p1 == "":
			if err := executeTemplate("create_account.gohtml", &templates.TemplateData{ErrorText: "Password cannot be empty"}, w, req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case p1 != p2:
			if err := executeTemplate("create_account.gohtml", &templates.TemplateData{ErrorText: "Passwords do not match"}, w, req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			if user, err := entity.AddNewUser(fullname, displayname, username, p1); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				token, err := utils.CreateToken(*user)
				if err != nil {
					fmt.Printf("ERROR: %+v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: time.Now().Add(3 * time.Hour)})
				http.Redirect(w, req, "/index.html", http.StatusSeeOther)
			}
		}
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Login requested: %s -> %s\n", req.Method, req.URL.Path)

	if req.FormValue("inputEmail") != "" {
		username := req.FormValue("inputEmail")
		password := req.FormValue("inputPassword")
		if user, err := utils.Authenticate(username, password); err == nil {
			fmt.Printf("User logged in: %s/%s\n", user.FullName, user.Username)
			token, err := utils.CreateToken(*user)
			if err != nil {
				fmt.Printf("ERROR: %+v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: time.Now().Add(3 * time.Hour)})
				http.Redirect(w, req, "/index.html", http.StatusSeeOther)
			}
		} else {
			fmt.Printf("Authentication failed: %s\n", username)
			if err2 := executeTemplate("login.gohtml", &templates.TemplateData{ErrorText: "Invalid username/password"}, w, req); err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		if err := executeTemplate("login.gohtml", nil, w, req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	})
	if err := executeTemplateNoUserInfo("logout.gohtml", nil, "", w, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func recoverAccountHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Recover account requested: %s -> %s\n", req.Method, req.URL.Path)

	switch req.Method {
	case "GET":
		if err := executeTemplate("recover_account.gohtml", nil, w, req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		username := req.FormValue("inputEmail")

		user := entity.FindUserByUsername(username)
		if user.ID > 0 {
			if err := executeTemplate("message.gohtml", &templates.TemplateData{ErrorText: "Email with a new password has been sent."}, w, req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {

		}
	}
}
