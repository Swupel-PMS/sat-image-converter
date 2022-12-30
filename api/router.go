package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"net/url"
	_ "sat-api/docs"
)

type Module interface {
	Base() string
	Convert() http.HandlerFunc
	Info() http.HandlerFunc
}

func router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

//	@title			Sat Image API
//	@version		1.0
//	@description	API for converting geo data to image
//	@termsOfService	https://swupel-pms.vercel.app/

//	@contact.name	Francisco Cardoso
//	@contact.url	https://swupel-pms.vercel.app/
//	@contact.email	swupelpms@gmail.com

//	@license.name	Copyright 2022 swupelpms
//	@license.url	mailto:swupelpms@gmail.com

//	@host		localhost:10001
//	@BasePath	/

// Start - start api
func Start(ctx context.Context, module Module) error {
	r := router()
	sat, err := url.JoinPath(module.Base(), "sat")
	if err != nil {
		return err
	}
	r.Post(sat, module.Convert())
	info, err := url.JoinPath(module.Base(), "info")
	if err != nil {
		return err
	}
	r.Get(info, module.Info())
	r.Get("/swagger/*", httpSwagger.Handler())
	s := http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", "10001"), // TODO: change to flag or env
	}
	go func() {
		<-ctx.Done()
		s.Close()
	}()
	return s.ListenAndServe()
}
