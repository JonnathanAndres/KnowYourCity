// The main binary is the KnowYourCity server.
package main

import (
	"log"
	"flag"
	"net/http"


	"gopkg.in/mgo.v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var hashKey = securecookie.GenerateRandomKey(32)
var blockKey = securecookie.GenerateRandomKey(32)
var sc = securecookie.New(hashKey, blockKey)

// Configures variables
var (
	serveraddr string
)

// main
func main() {
	serveraddr = *flag.String("srv", ":9004", "-srv=addr; set server listening address.Default :9004")
	flag.Parse()

	dbSession, err := dialDB()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	setAccountRoutes(router, dbSession)

	// TODO: Switch to TLS in prod(samson)
	http.ListenAndServe(serveraddr, router)
}

// dialDB
func dialDB() (*mgo.Session, error) {
	return mgo.Dial("localhost")
}

// setAccountRoutes
// REST good practics: trailing slash denotes a directory, while the lack of it denotes a file/resource
func setAccountRoutes(router *mux.Router, dbSession *mgo.Session) {
	router.HandleFunc("/account/register", withCORS(withClient(DbSession(db, SecretCode(Register))))).Methods("POST")
	router.HandleFunc("/account/login", withCORS(withClient(DbSession(db, Login)))).Methods("POST")
	router.HandleFunc("/account/password/forgot", withCORS(withClient(DbSession(db, ForgotPass)))).Methods("POST")
	router.HandleFunc("/page/account/password/reset", withCORS(withClient(DbSession(db, ResetPassPage)))).Methods("POST")
	router.HandleFunc("/account/password/reset", withCORS(withClient(DbSession(db, ResetPass)))).Methods("POST")
	router.HandleFunc("/account/2factor", withCORS(withClient(DbSession(db, L2Factor)))).Methods("POST")
	router.HandleFunc("/account/logout", withCORS(withClient(DbSession(db, Logout)))).Methods("POST")

	router.HandleFunc("/keys/code/get", withCORS(withClient(DbSession(dbSession, GetSecretCode))))
	router.HandleFunc("/keys/code/set", withCORS(withClient(DbSession(dbSession, SetSecretCode))))
}


/*
*/
