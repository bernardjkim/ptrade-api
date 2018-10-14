package router

import (
	"net/http"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	"github.com/bernardjkim/ptrade-api/pkg/types/routes"
	V1 "github.com/bernardjkim/ptrade-api/src/controllers/v1/router"
)

// Router is a wrapper for a mux router.
type Router struct {
	Router *mux.Router
}

// Init will initialize this router's routes and database connection.
func (r *Router) Init(db *xorm.Engine) {
	r.Router.Use(Middleware)

	// TODO: do some routes need to be part of the base routes?

	baseRoutes := GetRoutes(db)
	for _, route := range baseRoutes {
		r.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// Serve stylesheets
	r.Router.PathPrefix("/stylesheets/").Handler(http.StripPrefix("/stylesheets/",
		http.FileServer(http.Dir("./docs/stylesheets/"))))

	// Serve javascript files
	r.Router.PathPrefix("/javascripts/").Handler(http.StripPrefix("/javascripts/",
		http.FileServer(http.Dir("./docs/javascripts/"))))

	// Serve image files
	r.Router.PathPrefix("/images/").Handler(http.StripPrefix("/images/",
		http.FileServer(http.Dir("./docs/images/"))))

	var (
		v1SubRouter V1.SubRouter
	)

	v1SubRouter.Init(db)
	v1SubRoutes := v1SubRouter.GetRoutes(db)
	for name, pack := range v1SubRoutes {
		r.AttachSubRouterWithMiddleware(name, pack.Routes, pack.Middleware)
	}

}

// AttachSubRouterWithMiddleware will attach the subrouter to the given path
// using the given middleware function
func (r *Router) AttachSubRouterWithMiddleware(path string, subroutes routes.Routes, middleware []mux.MiddlewareFunc) (SubRouter *mux.Router) {

	SubRouter = r.Router.PathPrefix(path).Subrouter()

	// chain middleware functions
	for _, mw := range middleware {
		SubRouter.Use(mw)
	}

	// attach routes to sub router
	for _, route := range subroutes {
		SubRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

	}
	return
}

// NewRouter simply returns a new Rounter
func NewRouter() (r Router) {
	r.Router = mux.NewRouter().StrictSlash(true)
	return
}
