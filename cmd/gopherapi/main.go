package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/friendsofgo/gopherapi/cmd/sample-data"
	gopher "github.com/friendsofgo/gopherapi/pkg"
	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"
	"github.com/friendsofgo/gopherapi/pkg/log/logrus"
	"github.com/friendsofgo/gopherapi/pkg/modifying"
	"github.com/friendsofgo/gopherapi/pkg/removing"
	"github.com/friendsofgo/gopherapi/pkg/server"
	"github.com/friendsofgo/gopherapi/pkg/storage/inmem"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	var (
		hostName, _     = os.Hostname()
		defaultServerID = fmt.Sprintf("%s-%s", os.Getenv("GOPHERAPI_NAME"), hostName)
		defaultHost     = os.Getenv("GOPHERAPI_SERVER_HOST")
		defaultPort, _  = strconv.Atoi(os.Getenv("GOPHERAPI_SERVER_PORT"))
	)

	host := flag.String("host", defaultHost, "define host of the server")
	port := flag.Int("port", defaultPort, "define port of the server")
	serverID := flag.String("server-id", defaultServerID, "define server identifier")
	withData := flag.Bool("withData", false, "initialize the api with some gophers")
	flag.Parse()

	var gophers map[string]gopher.Gopher
	if *withData {
		gophers = sample.Gophers
	}

	logger := logrus.NewLogger()

	repo := inmem.NewRepository(gophers)
	fetchingService := fetching.NewService(repo, logger)
	addingService := adding.NewService(repo)
	modifyingService := modifying.NewService(repo)
	removingService := removing.NewService(repo)

	httpAddr := fmt.Sprintf("%s:%d", *host, *port)

	s := server.New(
		*serverID,
		fetchingService,
		addingService,
		modifyingService,
		removingService,
	)

	fmt.Println("The gopher server is on tap now:", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, s.Router()))
}
