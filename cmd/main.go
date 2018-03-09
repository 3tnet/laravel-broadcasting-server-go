package main

import (
	"github.com/ty666/go-socket.io"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"flag"
	"fmt"
	"os"
	"github.com/3tnet/laravel-broadcasting-server-go"
	"github.com/rs/cors"
	"github.com/3tnet/laravel-broadcasting-server-go/subscriber"
)

const Version = "0.0.2"

func arg(argName *string, envKey, defaultVal string) {
	if *argName == "" {
		*argName = os.Getenv(envKey)
		if *argName == "" {
			*argName = defaultVal
		}
	}
}

func main() {
	var (
		host              = flag.String("host", "", "Laravel broadcasting server host")
		authHost          = flag.String("auth_host", "", "Auth host")
		authEndpoint      = flag.String("auth_endpoint", "", "Auth endpoint")
		corsAllowedOrigin = flag.String("cors_allowed_origin", "", "Cors header allowedOrigins")
	)

	flag.Usage = usage
	flag.Parse()

	arg(host, "HOST", ":9999")
	arg(authHost, "AUTH_HOST", "http://localhost:8000")
	arg(authEndpoint, "AUTH_ENDPOINT", "/broadcasting/auth")
	arg(corsAllowedOrigin, "CORS_ALLOWED_ORIGIN", "http://localhost:8000")

	logger := log.New(os.Stdout, "[broadcasting] ", log.LstdFlags)
	logger.Println("version " + Version + " starting")

	ioServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	errLogger := log.New(os.Stderr, "[broadcasting] ", log.Llongfile|log.LstdFlags)
	s := broadcasting.NewServer(ioServer, *authHost, *authEndpoint, logger, errLogger)
	s.Listen(subscriber.NewHttpSubscriber(r))

	r.Handle("/socket.io/", ioServer)

	logger.Println("serving on " + *host)
	logger.Println("auth host:" + *authHost)
	logger.Println("auth endpoint:" + *authEndpoint)
	logger.Println("cors allowed origin:" + *corsAllowedOrigin)

	log.Fatal(http.ListenAndServe(*host, cors.New(cors.Options{
		AllowedOrigins:   []string{*corsAllowedOrigin},
		AllowCredentials: true,
	}).Handler(r)))

}

func usage() {
	fmt.Fprint(os.Stderr, "image server version: "+Version+"\nOptions:\n")
	flag.PrintDefaults()
}
