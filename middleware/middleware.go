package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

type Stack []Middleware

// Method untuk menambah middleware
func (s Stack) With(mw ...Middleware) Stack {
	return append(s, mw...)
}

// Method untuk apply ke handler
func (s Stack) Then(handler http.HandlerFunc) http.Handler {
	if len(s) == 0 {
		return handler
	}
	return Chain(handler, s...)
}
