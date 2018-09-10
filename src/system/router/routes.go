package router

import (
	"net/http"

	"github.com/bkim0128/bjstock-rest-service/pkg/types/routes"

	// AuthHandler "github.com/bkim0128/bjstock-rest-service/src/controllers/auth"
	// PortfolioHandler "github.com/bkim0128/bjstock-rest-service/src/controllers/portfolio"
	// StockHandler "github.com/bkim0128/bjstock-rest-service/src/controllers/stock"

	"github.com/go-xorm/xorm"
)

// Middleware function for Router
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// GetRoutes returns a list of routes handled by this router
func GetRoutes(db *xorm.Engine) routes.Routes {
	// AuthHandler.Init(db)
	// PortfolioHandler.Init(db)
	// StockHandler.Init(db)

	return routes.Routes{

		// Warning: Composite literal uses unkeyed fields.
		// Can remove warnings by including field names (field: value).

		// routes.Route{"AuthStore", "POST", "/auth/login", AuthHandler.Login},
		// routes.Route{"AuthCheck", "GET", "/auth/check", AuthHandler.Check},
		// routes.Route{"AuthSignout", "POST", "/auth/signup", AuthHandler.SignUp},
		// routes.Route{"PortfolioStocks", "GET", "/portfolio/stocks", PortfolioHandler.GetStocks},
		// routes.Route{"PortfolioStocks", "POST", "/portfolio/buy-shares", PortfolioHandler.BuyShares},
		// routes.Route{"StockList", "GET", "/stock/list", StockHandler.GetStockList},
	}
}
