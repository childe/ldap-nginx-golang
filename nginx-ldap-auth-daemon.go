package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/ldap.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var options = &struct {
	configArg string
}{
	configArg: "config.json",
}

var config map[string]string

func ladpAuth(username string, password string) bool {
	var l *ldap.Conn
	var err error
	l, err = ldap.DialTLS("tcp", config["ldapserver"], &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Println(err)
	}

	err = l.Bind(config["binddn"], config["bindpw"])

	if err != nil {
		log.Println(err)
	}

	filterString := fmt.Sprintf(config["filter"], username)
	log.Println(filterString)

	searchRequest := ldap.NewSearchRequest(
		config["basedn"],
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filterString,
		nil,
		nil,
	)

	searchResult, searchErr := l.Search(searchRequest)
	if searchErr != nil {
		log.Println(searchErr)
		return false
	}

	if len(searchResult.Entries) != 1 {
		log.Println("search result is not only one(more or less)")
		return false
	}

	userDN := searchResult.Entries[0].DN
	log.Println(userDN)

	//attributes := sr.Entries[0].Attributes
	//for _, attr := range attributes {
	//log.Printf("%s: %s\n", attr.Name, attr.Values)
	//}

	err = l.Bind(userDN, password)
	if err != nil {
		log.Println(err)
		return false
	} else {
		return true
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	authorization := r.Header.Get("Authorization")
	log.Printf("original authorization: %s", authorization)

	if authorization == "" {
		w.Header().Add("WWW-Authenticate", "Basic realm=\"\"")
		w.WriteHeader(401)
		//fmt.Fprintf(w, "")
		return
	}

	authorizationBytes, err := base64.StdEncoding.DecodeString(authorization[len("Basic "):])

	if err != nil {
		log.Println(err)
		w.Header().Add("WWW-Authenticate", "Basic realm=\"\"")
		w.WriteHeader(401)
		log.Println(err)
		return
	}

	authorizationValue := string(authorizationBytes)
	log.Printf("authorization: %s", string(authorizationBytes))

	userANDpw := strings.SplitN(authorizationValue, ":", 2)
	if len(userANDpw) != 2 {
		w.Header().Add("WWW-Authenticate", "Basic realm=\"\"")
		w.WriteHeader(401)
		log.Println("Authenticate Value Format Error")
		return
	}

	username := userANDpw[0]
	password := userANDpw[1]

	log.Printf("username: %s", username)

	if ladpAuth(username, password) {
		w.WriteHeader(200)
	} else {
		w.Header().Add("WWW-Authenticate", "Basic realm=\"\"")
		w.WriteHeader(401)
		log.Println("Auth Failed")
	}
	return
}

func init() {
	flag.StringVar(&options.configArg, "config", options.configArg, "path to ldap-auth-daemon configuration file")
}

func main() {
	configValue, err := ioutil.ReadFile(options.configArg)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if err != nil {
		log.Fatal(err)
		return
	}
	if err := json.Unmarshal(configValue, &config); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", authHandler)
	http.ListenAndServe(":8080", nil)
}
