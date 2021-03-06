package main

import (
	"os"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  "{\"time\":\"${time_rfc3339_nano}\", \"remote_ip\":\"${remote_ip}\", \"host\":\"${host}\", \"method\":\"${method}\", \"uri\":\"${uri}\", \"status\":${status}, latency:${latency}, \"latency_human\":\"${latency_human}\", \"bytes_in\":${bytes_in}, \"bytes_out\":${bytes_out}}\n",
		Output:  os.Stdout,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderAccessControlAllowOrigin},
	}))
	e.Use(middleware.Recover())

	dbh := DBHandler{}
	err := dbh.initDB()
	if err != nil {
		e.Logger.Panic(err)
	}
	e.Logger.Print("DB handler initiated")

	e.File("/favicon.ico", "images/favicon.ico")

	// Indicators
	e.GET("/indicators", dbh.getIndicators)
	e.GET("/indicators/:indicator", dbh.getIndicator)

	// Repository Group
	r := e.Group("/repo/:owner/:repo")

	// Repository
	r.GET("", dbh.getRepository)

	// Health Group
	h := r.Group("/health")

	// Health Indicators
	h.GET("/docs", dbh.getRepositoryHealthDocs)
	h.GET("/response_times", dbh.getRepositoryHealthResponseTimes)
	h.GET("/pull_requests", dbh.getRepositoryHealthPullRequests)

	// data, err := json.MarshalIndent(e.Routes(), "", "  ")
	// if err != nil {
	// 	e.Logger.Fatal(err)
	// }
	// ioutil.WriteFile("routes.json", data, 0644)
	// e.Logger.Print("routes.json file created")

	e.Logger.Printf("%d routes created", len(e.Routes()))

	e.Server.Addr = ":8080"

	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
