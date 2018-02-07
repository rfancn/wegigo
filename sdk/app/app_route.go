package app

import "net/http"

type IAppRoute interface {
	GetMethod() string
	GetUrl() string
	GetHandler() func(http.ResponseWriter, *http.Request)
}

type AppRoute struct {
	Method  string
	Url     string
	Handler func(http.ResponseWriter, *http.Request)
}

func NewAppRoute(method string, url string, handler func(http.ResponseWriter, *http.Request)) *AppRoute {
	return &AppRoute{Method: method, Url: url, Handler:handler}
}

func (r *AppRoute) GetMethod() string {
	return r.Method
}

func (r *AppRoute) GetUrl() string {
	return r.Url
}

func (r *AppRoute) GetHandler() func(http.ResponseWriter, *http.Request) {
	return r.Handler
}