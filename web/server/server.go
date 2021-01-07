package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

//блокирующая функция
func Start(address string, router *mux.Router) {
	srv := &http.Server{
		Addr: address,

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
	defer srv.Close()
}
