package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/seven320/go_server/controller"
	"github.com/seven320/go_server/db"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world %q", r.URL.Path)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file. %s", err)
	}

	datasource := os.Getenv("DATABASE_DATASOURCE")
	if datasource == "" {
		log.Fatal("Cannot get datasource for database.")
	}

	s := NewServer()
	s.Init(datasource)
	s.Run(os.Getenv("PORT"))
	http.HandleFunc("/", handler) // ハンドラを登録してウェブに表示
	http.ListenAndServe(":1991", nil)
}

type Server struct {
	db     *sqlx.DB
	router *mux.Router
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init(datasource string) {
	cs := db.NewDB(datasource)
	dbcon, err := cs.Open()
	if err != nil {
		log.Fatalf("failed db init. %s", err)
	}
	s.db = dbcon
	s.router = s.Route()
}

func (s *Server) Run(addr string) {
	log.Printf("Listening on port %s", addr)
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", addr),
		handlers.CombinedLoggingHandler(os.Stdout, s.router),
	)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Route() *mux.Router {
	// user 管理やログインに使う
	commonChain := alice.New(
		RecoverMiddleware,
	)

	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/twitter_image").Handler(commonChain.Then(NewPublicHandler()))

	twitterimageController := controller.NewTwitterImage(s.db)

	r.Methods(http.MethodGet).Path("/show_databases").Handler(commonChain.Then(AppHandler{twitterimageController.Show}))
	return r
}

//package  middle ware
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				log.Printf("panic: %s", err)
				http.Error(w, http.StatusText(
					http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// package handler
type AppHandler struct {
	h func(http.ResponseWriter, *http.Request) (int, interface{}, error)
}

func (a AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, res, err := a.h(w, r)
	if err != nil {
		respondErrorJson(w, status, err)
		return
	}
	respondJSON(w, status, res)
	return
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload) //json に変換している
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondErrorJson(w http.ResponseWriter, code int, err error) {
	log.Printf("code=%d, err=%s", code, err)
	if e, ok := err.(*HTTPError); ok {
		respondJSON(w, code, e)
	} else if err != nil {
		he := HTTPError{
			Message: err.Error(),
		}
		respondJSON(w, code, he)
	}
}

//package httputil
type HTTPError struct {
	Message string `json:"message"`
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("message=%v", he.Message)
}

//package sample
type PublicHandler struct{}

func NewPublicHandler() *PublicHandler {
	return &PublicHandler{}
}

func (h *PublicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Message: "Hello from a public endpoint! You don't need to be authenticated to see this.",
	}
	WriteJSON(resp, w, http.StatusOK)
}

type Response struct {
	Message string `json:"message"`
}

func WriteJSON(v interface{}, w http.ResponseWriter, statusCode int) {
	jsonResponse, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, e := w.Write(jsonResponse); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}
