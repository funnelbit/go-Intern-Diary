package web

//go:generate go-assets-builder --package=web --output=./templates-gen.go --strip-prefix="/templates/" --variable=Templates ../templates

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/justinas/nosurf"

	"github.com/hatena/go-Intern-Diary/service"
)

type Server interface {
	Handler() http.Handler
}

const sessionKey = "DIARY_SESSION"

var templates map[string]*template.Template

func init() {
	var err error
	templates, err = loadTemplates()
	if err != nil {
		panic(err)
	}
}

func loadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)
	bs, err := ioutil.ReadAll(Templates.Files["main.tmpl"])
	if err != nil {
		return nil, err
	}
	mainTmpl := template.Must(template.New("main.tmpl").Parse(string(bs)))
	for fileName, file := range Templates.Files {
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		mainTmpl := template.Must(mainTmpl.Clone())
		templates[fileName] = template.Must(mainTmpl.New(fileName).Parse(string(bs)))
	}
	return templates, nil
}

func NewServer(app service.DiaryApp) Server {
	return &server{app: app}
}

type server struct {
	app service.DiaryApp
}

func (s *server) Handler() http.Handler {
	router := httptreemux.New()
	handle := func(method, path string, handler http.Handler) {
		router.UsingContext().Handler(method, path, csrfMiddleware(loggingMiddleware(headerMiddleware(handler))))
	}

	handle("GET", "/signup", s.willSignupHandler())
	handle("POST", "/signup", s.signupHandler())

	handle("GET", "/signin", s.willSigninHandler())
	handle("POST", "/signin", s.signinHandler())

	handle("GET", "/signout", s.willSignoutHandler())
	handle("POST", "/signout", s.signoutHandler())

	handle("GET", "/diaries", s.showDiariesHandler())
	handle("GET", "/create_diary", s.willCreateDiary())
	handle("POST", "/create_diary", s.CreateDiary())

	handle("GET", "/diary/:id", s.showDiaryHandler())

	return router
}

var csrfMiddleware = func(next http.Handler) http.Handler {
	return nosurf.New(next)
}

var csrfToken = func(r *http.Request) string {
	return nosurf.Token(r)
}

func (s *server) willSignupHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, "signup.tmpl", nil)
	})
}

func (s *server) willSigninHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, r, "signin.tmpl", nil)
	})
}

func (s *server) signupHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, password := r.FormValue("name"), r.FormValue("password")
		if err := s.app.CreateNewUser(name, password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := s.app.FindUserByName(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(24 * time.Hour)
		token, err := s.app.CreateNewToken(user.ID, expiresAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    sessionKey,
			Value:   token,
			Expires: expiresAt,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (s *server) signinHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, password := r.FormValue("name"), r.FormValue("password")
		ok, err := s.app.LoginUser(name, password)
		if err != nil || !ok {
			http.Error(w, "user not found or invalid password", http.StatusBadRequest)
			return
		}

		user, err := s.app.FindUserByName(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(24 * time.Hour)
		token, err := s.app.CreateNewToken(user.ID, expiresAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    sessionKey,
			Value:   token,
			Expires: expiresAt,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (s *server) willSignoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionKey)
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		s.renderTemplate(w, r, "signout.tmpl", nil)
	})
}

func (s *server) signoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    sessionKey,
			Value:   "",
			Expires: time.Unix(0, 0),
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (s *server) showDiariesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionKey)
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		user, err := s.app.FindUserByToken(cookie.Value)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		diaries, err := s.app.FindDiariesByUserID(user.ID)
		var el error
		if err != nil || diaries == nil {
			el = errors.New("作品を読み込めませんでした")
		}

		s.renderTemplate(w, r, "diaries.tmpl", map[string]interface{}{
			"User":    user,
			"Diaries": diaries,
			"Error":   el,
		})
	})
}

func (s *server) willCreateDiary() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionKey)
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		user, err := s.app.FindUserByToken(cookie.Value)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		}

		s.renderTemplate(w, r, "new_diary.tmpl", map[string]interface{}{
			"User": user,
		})
	})
}

func (s *server) CreateDiary() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionKey)
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		user, err := s.app.FindUserByToken(cookie.Value)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		}

		name := r.FormValue("name")

		diary, err := s.app.CreateNewDiary(user, name)
		if err != nil || diary == nil {
			http.Redirect(w, r, "/create_diary", http.StatusSeeOther)
		}

		http.Redirect(w, r, "/diaries", http.StatusSeeOther)
	})
}

func (s *server) showDiaryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionKey)
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		user, err := s.app.FindUserByToken(cookie.Value)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		}

		diaryID, err := strconv.ParseUint(s.getParams(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid entry id", http.StatusBadRequest)
			return
		}

		diary, err := s.app.FindDiaryByID(diaryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		s.renderTemplate(w, r, "diary.tmpl", map[string]interface{}{
			"Diary":    diary,
			"Articles": []string{},
		})
	})
}

func (s *server) renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["CSRFToken"] = csrfToken(r)
	err := templates[tmpl].ExecuteTemplate(w, "main.tmpl", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) getParams(r *http.Request, name string) string {
	return httptreemux.ContextParams(r.Context())[name]
}
