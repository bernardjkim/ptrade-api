package router

import (
	"context"
	"log"
	"net/http"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	"github.com/bkim0128/bjstock-rest-service/pkg/types/routes"
	"github.com/bkim0128/bjstock-rest-service/src/system/jwt"

	Users "github.com/bkim0128/bjstock-rest-service/pkg/types/users"
	SessionHandler "github.com/bkim0128/bjstock-rest-service/src/controllers/v1/sessions"
	StockHandler "github.com/bkim0128/bjstock-rest-service/src/controllers/v1/stocks"
	TransactionHandler "github.com/bkim0128/bjstock-rest-service/src/controllers/v1/transactions"
	UserHandler "github.com/bkim0128/bjstock-rest-service/src/controllers/v1/users"
)

var db *xorm.Engine

// AuthMiddleware handles authentication of requests received by router
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check if token is present
		tokenVal := r.Header.Get("X-App-Token")
		if len(tokenVal) < 1 {
			log.Println("Ignoring request. No token present.")
			http.Error(w, "No token provided for validation.", http.StatusUnauthorized)
			return
		}

		// get owner of token
		user, err := jwt.GetUserFromToken(db, tokenVal)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// pass on user id to next handler
		ctx := context.WithValue(r.Context(), Users.UserIDKey, user.ID)

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRoutes returns mappings from names to subroute packages.
func GetRoutes(DB *xorm.Engine) (SubRoute map[string]routes.SubRoutePackage) {
	db = DB

	SessionHandler.Init(db)
	TransactionHandler.Init(db)
	StockHandler.Init(db)
	UserHandler.Init(db)

	/* ROUTES */

	// Warning: Composite literal uses unkeyed fields.
	// Can remove warnings by including field names (field: value).
	SubRoute = map[string]routes.SubRoutePackage{
		"/v1": routes.SubRoutePackage{
			Routes: routes.Routes{
				routes.Route{"GetIndex", "GET", "", GetIndex},
			},
			Middleware: []mux.MiddlewareFunc{},
		},

		"/v1/sessions": routes.SubRoutePackage{
			Routes: routes.Routes{
				routes.Route{"CreateSession", "POST", "", SessionHandler.CreateSession},
				routes.Route{"DeleteSession", "DELETE", "", NotImplemented},
			},
			Middleware: []mux.MiddlewareFunc{},
		},
		"/v1/stocks": routes.SubRoutePackage{
			Routes: routes.Routes{
				routes.Route{"GetStocks", "GET", "", StockHandler.GetStocks},
			},
			Middleware: []mux.MiddlewareFunc{},
		},

		// NOTE: order matters, match /user/.../transactions subroute before /users
		// was not using the assigned middleware
		"/v1/users/{ID:[0-9]+}/transactions": routes.SubRoutePackage{
			Routes: routes.Routes{

				// TODO: currently have users GET/POST transactions directly.
				// maybe want to have user create orders first and later
				// execute transaction
				routes.Route{"GetUserTxns", "GET", "", TransactionHandler.GetTransactions},
				routes.Route{"CreateUserTxn", "POST", "", TransactionHandler.CreateTransaction},

				routes.Route{"GetUserTxn", "GET", "/{txnID:[0-9]+}", NotImplemented},
			},
			Middleware: []mux.MiddlewareFunc{AuthMiddleware},
		},
		"/v1/users": routes.SubRoutePackage{
			Routes: routes.Routes{

				routes.Route{"GetUsers", "GET", "", NotImplemented},
				routes.Route{"CreateUser", "POST", "", UserHandler.CreateUser},

				routes.Route{"GetUser", "GET", "/{ID:[0-9]+}", UserHandler.GetUser},
			},
			Middleware: []mux.MiddlewareFunc{},
		},
	}
	return
}

// NotImplemented handler is used for API endpoints not yet implemented and will
// return the message "Not Implemented".
var GetIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
})

// NotImplemented handler is used for API endpoints not yet implemented and will
// return the message "Not Implemented".
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
