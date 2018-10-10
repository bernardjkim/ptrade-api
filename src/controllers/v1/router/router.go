package router

import (
	"context"
	"log"
	"net/http"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	"github.com/bernardjkim/ptrade-api/pkg/types/routes"
	"github.com/bernardjkim/ptrade-api/src/controllers/v1/sessions"
	"github.com/bernardjkim/ptrade-api/src/controllers/v1/stocks"
	"github.com/bernardjkim/ptrade-api/src/controllers/v1/users"
	"github.com/bernardjkim/ptrade-api/src/system/jwt"

	Users "github.com/bernardjkim/ptrade-api/pkg/types/users"
)

// SubRouter needs to be initialized with a db connection
type SubRouter struct {
	DB *xorm.Engine
}

// Init will initialize this router's routes and database connection.
func (sr *SubRouter) Init(db *xorm.Engine) {
	sr.DB = db
}

// AuthMiddleware handles authentication of requests received by router
func (sr *SubRouter) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check if token is present
		tokenVal := r.Header.Get("Session-Token")
		if len(tokenVal) < 1 {
			log.Println("Ignoring request. No token present.")
			http.Error(w, "No token provided for validation.", http.StatusUnauthorized)
			return
		}

		// get owner of token
		user, err := jwt.GetUserFromToken(sr.DB, tokenVal)
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
func (sr *SubRouter) GetRoutes(DB *xorm.Engine) (SubRoute map[string]routes.SubRoutePackage) {
	var (
		userHandler    users.UserHandler
		stockHandler   stocks.StockHandler
		sessionHandler sessions.SessionHandler
	)

	sessionHandler.Init(sr.DB)
	stockHandler.Init(sr.DB)
	userHandler.Init(sr.DB)

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
				// TODO: not rest standards, what to do about validating auth token???
				routes.Route{"ValidateSession", "GET", "/validate", sessionHandler.Validate},
				routes.Route{"CreateSession", "POST", "", sessionHandler.CreateSession},
				routes.Route{"DeleteSession", "DELETE", "", NotImplemented},
			},
			Middleware: []mux.MiddlewareFunc{},
		},

		"/v1/stocks": routes.SubRoutePackage{
			Routes: routes.Routes{
				routes.Route{"GetStocks", "GET", "", stockHandler.GetStocks},
			},
			Middleware: []mux.MiddlewareFunc{},
		},

		// NOTE: order matters, match /user/.../transactions subroute before /users
		// was not using the assigned middleware
		"/v1/users/{ID:[0-9]+}/stocktransactions": routes.SubRoutePackage{
			Routes: routes.Routes{

				// TODO: currently have users GET/POST transactions directly.
				// maybe want to have user create orders first and later
				// execute transaction
				// routes.Route{"GetUserStockTxns", "GET", "", stockTransactionHandler.GetTransactions},
				// routes.Route{"CreateUserStockTxn", "POST", "", stockTransactionHandler.CreateTransaction},

				routes.Route{"GetUserStockTxn", "GET", "/{txnID:[0-9]+}", NotImplemented},
			},
			Middleware: []mux.MiddlewareFunc{sr.AuthMiddleware},
		},

		"/v1/users/{ID:[0-9]+}/bankingtransactions": routes.SubRoutePackage{
			Routes: routes.Routes{
				// routes.Route{"GetUserBankingTxns", "GET", "", bankingTransactionHandler.GetTransactions},
				// routes.Route{"CreateUserBankingTxn", "POST", "", bankingTransactionHandler.CreateTransaction},
				routes.Route{"GetUserBankingTxn", "GET", "/{txnID:[0-9]+}", NotImplemented},
			},
			Middleware: []mux.MiddlewareFunc{sr.AuthMiddleware},
		},

		"/v1/users": routes.SubRoutePackage{
			Routes: routes.Routes{

				routes.Route{"GetUsers", "GET", "", NotImplemented},
				routes.Route{"CreateUser", "POST", "", userHandler.CreateUser},

				routes.Route{"GetUser", "GET", "/{ID:[0-9]+}", userHandler.GetUser},
			},
			Middleware: []mux.MiddlewareFunc{},
		},
	}
	return
}

// GetIndex will return the index for the api
var GetIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/index.html")
})

// NotImplemented handler is used for API endpoints not yet implemented and will
// return the message "Not Implemented".
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
