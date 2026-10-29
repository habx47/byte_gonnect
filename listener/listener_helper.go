package listener

import (
	"errors"
	"fmt"
	"strings"

	"github.com/habx47/byte_gonnect/router"
)

var GETRsp_root = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 25\r\nConnection: close\r\n\r\nHello, from the root page")
var GETRsp_hello = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 26\r\nConnection: close\r\n\r\nHello, from the hello page")

var CreateRouteErr = errors.New("failed to create a route. missing '/'")

func createRoutes(endpointResps router.EndpointResps, routes router.Routes) error {
	for key := range endpointResps {
		_, route, found := strings.Cut(key, "/")
		if !found {
			return fmt.Errorf("%w", CreateRouteErr)
		}

		route = "/" + route

		routes[route] = true
	}
	return nil
}
