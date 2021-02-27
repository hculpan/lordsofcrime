package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/hculpan/lordsofcrime/templates"
	"github.com/hculpan/lordsofcrime/utils"
)

var templateList *template.Template

var fileserver = http.FileServer(http.Dir("./resources"))

// SetupRoutes set up the http routes and handlers
func SetupRoutes() {
	LoadTemplates()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login.html", loginHandler)
	http.HandleFunc("/create_account.html", createAccountHandler)
	http.HandleFunc("/logout.html", logoutHandler)
	http.HandleFunc("/recover_account.html", recoverAccountHandler)
	http.Handle("/resources/", http.StripPrefix("/resources", fileserver))
}

func executeTemplate(name string, templateData *templates.TemplateData, w http.ResponseWriter, req *http.Request) error {
	cookies := map[string]string{}
	for _, v := range req.Cookies() {
		cookies[v.Name] = v.Value
	}

	if templateData == nil {
		templateData = &templates.TemplateData{Cookies: cookies}
	} else {
		templateData.Cookies = cookies
	}

	if token, err := req.Cookie("token"); err == nil {
		if claims, err := utils.DecodeToken(token.Value); err == nil {
			templateData.UserDisplayName = claims.DisplayName
		}
	} else if err.Error() != "http: named cookie not present" {
		fmt.Printf("Error decoding token: %v\n", err)
	}

	return templateList.ExecuteTemplate(w, name, templateData)
}

func executeTemplateNoUserInfo(name string, data interface{}, errorText string, w http.ResponseWriter, req *http.Request) error {
	cookies := map[string]string{}
	for _, v := range req.Cookies() {
		cookies[v.Name] = v.Value
	}

	templateData := templates.TemplateData{
		ErrorText: errorText,
		Cookies:   cookies,
		Data:      data,
	}

	return templateList.ExecuteTemplate(w, name, templateData)
}

// LoadTemplates load the templates
func LoadTemplates() {
	templateList = template.Must(template.ParseGlob("./templates/*.gohtml"))
}
