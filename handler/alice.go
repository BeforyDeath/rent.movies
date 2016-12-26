package handler

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	a "github.com/justinas/alice"
)

type alice struct {
	base a.Chain
}

func (m *alice) NewChain(constructors ...a.Constructor) {
	m.base = a.New(constructors...)
}

func (m alice) AddChain(h http.HandlerFunc, c ...a.Constructor) httprouter.Handle {
	b := m.base.Extend(a.New(c...))
	return m.stripParams(b.ThenFunc(h))
}

func (m alice) stripParams(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}
