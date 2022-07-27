package service

import (
	"net/http"
	"strconv"

	"github.com/Blackmamba23/cars-fleet-service/internal/fleet"
	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/context"
)

// CarsFleetService will keep a handle on the saved items repository and implement
// the gizmo/server.JSONService interface.
type CarsFleetService struct {
	fleetRepo fleet.Repository
}

// NewCarsFleetService will attempt to instantiate a new repository and service.
func NewCarsFleetService(fleetRepo fleet.Repository) (*CarsFleetService, error) {
	return &CarsFleetService{
		fleetRepo,
	}, nil
}

// Prefix is to implement gizmo/server.Service interface. The string will be prefixed to all endpoint
// routes.
func (s *CarsFleetService) Prefix() string {
	return "/svc/fleet/v1"
}

// Middleware provides a hook to add service-wide http.Handler middleware to the service.
// In this example we are using it to add GZIP compression to our responses.
// This method helps satisfy the server.Service interface.
func (s *CarsFleetService) Middleware(h http.Handler) http.Handler {
	// wrap the response with our GZIP Middleware
	return context.ClearHandler(gziphandler.GzipHandler(h))
}

// JSONMiddleware provides a hook to add service-wide middleware for how JSONEndpoints
// should behave. In this example, we’re using the hook to check for a header to
// identify and authorize the user. This method helps satisfy the server.JSONService interface.
func (s *CarsFleetService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (code int, res interface{}, err error) {

		// wrap our endpoint with an auth check and call it
		code, res, err = authCheck(j)(r)

		// if the endpoint returns an unexpected error, return a generic message
		// and log it.
		if err != nil && code != http.StatusUnauthorized {
			// LogWithFields will add all the request context values
			// to the structured log entry along some other request info
			server.LogWithFields(r).WithField("error", err).Error("unexpected service error")
			return http.StatusServiceUnavailable, nil, ServiceUnavailableErr
		}

		return code, res, err
	}
}

// idKey is a type to use as a key for storing data in the request context.
type idKey int

// userIDKey can be used to store/retrieve a user ID in a request context.
const userIDKey idKey = 0

// authCheck is a JSON middleware to check the request for a valid API_TOKEN
// header and set it into the request context. If the header is invalid
// or does not exist, a 401 response will be returned.
func authCheck(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (code int, res interface{}, err error) {
		// check for User ID header injected by API Gateway
		idStr := r.Header.Get("API_TOKEN")
		// verify it's an int
		id, err := strconv.ParseUint(idStr, 10, 64)
		// reject request if bad/no user ID
		if err != nil || id == 0 {
			return http.StatusUnauthorized, nil, UnauthErr
		}
		// set the ID in context if we're good
		context.Set(r, userIDKey, id)

		return j(r)
	}
}

// Endpoints returns the endpoints for our stream service.
func (s *CarsFleetService) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/cars": {
			"GET": server.JSONToHTTP(s.GetCarsByName).ServeHTTP,
		},
		"/car": {
			"GET": server.JSONToHTTP(s.GetCarByName).ServeHTTP,
		},
	}
}

type (
	// Response is a generic struct for responding with a simple JSON message.
	jsonResponse struct {
		Message string `json:"message"`
	}
	// jsonErr is a tiny helper struct to make displaying errors in JSON better.
	jsonErr struct {
		Err string `json:"error"`
	}
)

func (e *jsonErr) Error() string {
	return e.Err
}

var (
	// ServiceUnavailableErr is a global error that will get returned when we are experiencing
	// technical issues.
	ServiceUnavailableErr = &jsonErr{"sorry, this service is currently unavailable"}
	// UnauthErr is a global error returned when the user does not supply the proper
	// authorization headers.
	UnauthErr = &jsonErr{"please include a valid API_TOKEN header in the request"}
)
