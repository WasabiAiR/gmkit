package metrics

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

// SanitizeRoute cleans up a route for the Gorilla Mux to be used as a metric name
func SanitizeRoute(method, route string) string {
	// replace slashes with underscores
	route = strings.Replace(route, "/", "_", -1)

	// strip out the curly braces
	route = strings.Replace(route, "{", "", -1)
	route = strings.Replace(route, "}", "", -1)

	// prefix the route with the http verb
	if !strings.HasPrefix(route, "_") {
		route = "_" + route
	}

	return strings.ToUpper(method) + route
}

// InstrumentRouter walks all the handlers on a mux.Router and wraps each route's
// handler with a metrics Handler for reporting sanitized route timing metrics
func InstrumentRouter(r *mux.Router) error {
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, _ := route.GetMethods()
		tmpl, _ := route.GetPathTemplate()

		if len(methods) == 0 {
			return nil
		}

		sort.Strings(methods)
		route.Handler(Handler(route.GetHandler(), SanitizeRoute(strings.Join(methods, "_"), tmpl)))
		return nil
	})

	return fmt.Errorf("instrumenting router: %w", err)
}
