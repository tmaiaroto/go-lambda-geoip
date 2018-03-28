package main

import (
	"github.com/oschwald/geoip2-golang"
	"github.com/fatih/structs"
	aegis "github.com/tmaiaroto/aegis/framework"
	"net"
	"context"
	"net/url"
	"log"
	"errors"
)

var db *geoip2.Reader

func main() {
	// Load the database in main() in order to benefit from Lambda container re-use.
	data, err := Asset("GeoLite2-City.mmdb")
	if err == nil {
		db, _ = geoip2.FromBytes(data)
		defer db.Close()

		// Handle API requests
		router := aegis.NewRouter(fallThrough)
		router.Handle("GET", "/", root)
		// router.Listen()

		// Handle RPCs
		rpcRouter := aegis.NewRPCRouter()
		rpcRouter.Handle("lookup", lookupProcedure)

		handlers := aegis.Handlers{
			Router:    router,
			RPCRouter: rpcRouter,
		}
		handlers.Listen()
	} else {
		log.Println("Could not load GeoLite2-City.mmdb. Is it included in the binary?", err)
	}
}

func fallThrough(ctx context.Context, evt *aegis.APIGatewayProxyRequest, res *aegis.APIGatewayProxyResponse, params url.Values) error {
    res.StatusCode = 404
    return nil
}

func root(ctx context.Context, evt *aegis.APIGatewayProxyRequest, res *aegis.APIGatewayProxyResponse, params url.Values) error {
	record, err := lookup(evt.RequestContext.Identity.SourceIP)
	res.JSON(200, record)
    return err
}

func lookupProcedure(ctx context.Context, evt *map[string]interface{}) (map[string]interface{}, error) {
	var resp map[string]interface{}
	if evt != nil {
		derefEvt := *evt
		record, err := lookup(derefEvt["ipAddress"].(string))
		if err == nil {
			resp = structs.Map(record)
			return resp, err
		}
		return resp, err
	}
	return resp, errors.New("no IP address passed to procedure")
}

func lookup(ipAddress string) (*geoip2.City, error) {
	parsedIP := net.ParseIP(ipAddress)
	return db.City(parsedIP)
}