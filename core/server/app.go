package server

import (
	mapFilesHttp "GOSecretProject/core/mapfiles/delivery/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"net/http"
)

type App struct {

}

func NewApp() *App {
		return &App{
	}
}

func (app *App) StartRouter() {
	router := mux.NewRouter()

	commonRouter := router.PathPrefix("/api").Subrouter()

	mapFilesHttp.RegisterHTTPEndpoints(commonRouter)

	http.Handle("/", router)

	port := 8001
	golog.Infof("Server started at port :%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port),
		nil)

	if err != nil {
		golog.Error("Server haha failed: ", err)
	}
}