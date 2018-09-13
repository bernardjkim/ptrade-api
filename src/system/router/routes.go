package router

import (
	"net/http"

	"github.com/go-xorm/xorm"

	"github.com/bernardjkim/ptrade-api/pkg/types/routes"
)

// Middleware function for Router
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// GetRoutes returns a list of routes handled by this router
func GetRoutes(db *xorm.Engine) routes.Routes {

	// Warning: Composite literal uses unkeyed fields.
	// Can remove warnings by including field names (field: value).
	return routes.Routes{
		routes.Route{"GetIndex", "GET", "/", GetIndex},
	}
}

// GetIndex handler will serve the index for the api
var GetIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
})

// NotImplemented handler is used for API endpoints not yet implemented and will
// return the message "Not Implemented".
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
