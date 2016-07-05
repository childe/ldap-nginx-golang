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
	configArg          string
	host               string
	port               int
	ldapserver         string
	basedn             string
	binddn             string
	bindpw             string
	filter             string
	insecureSkipVerify bool
}{
	configArg: "config.json",
}

func ladpAuth(username string, password string) bool {
	l, err := ldap.DialTLS("tcp", options.ldapserver, &tls.Config{InsecureSkipVerify: options.insecureSkipVerify})
	if err != nil {
		log.Printf("connecting ldap server failed: %s\n", err)
		return false
	}

	bindErr := l.Bind(options.binddn, options.bindpw)

	if bindErr != nil {
		log.Printf("bind failed: %s", bindErr)
		return false
	}

	filterString := fmt.Sprintf(options.filter, username)
	log.Println(filterString)

	searchRequest := ldap.NewSearchRequest(
		options.basedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filterString,
		nil,
		nil,
	)

	searchResult, searchErr := l.Search(searchRequest)
	if searchErr != nil {
		log.Printf("search failed: %s\n", searchErr)
		return false
	}

	if len(searchResult.Entries) != 1 {
		log.Println("search result is not only one(more or less)")
		return false
	}

	userDN := searchResult.Entries[0].DN
	log.Println(userDN)

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
		return
	}

	authorizationValue := string(authorizationBytes)
	log.Printf("authorization: %s", string(authorizationBytes))

	userANDpw := strings.SplitN(authorizationValue, ":", 2)
	if len(userANDpw) != 2 {
		log.Println("Authenticate Value Format Error")
		w.Header().Add("WWW-Authenticate", "Basic realm=\"\"")
		w.WriteHeader(401)
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
	flag.StringVar(&options.host, "host", "127.0.0.1", "ip/host that bind to, default 127.0.0.1")
	flag.IntVar(&options.port, "port", 8080, "port that bind to, default 8080")
	flag.StringVar(&options.ldapserver, "ldapserver", "", "required. ldapserver")
	flag.StringVar(&options.basedn, "basedn", "", "required. basedn.")
	flag.StringVar(&options.binddn, "binddn", "", "required if search action need bind first")
	flag.StringVar(&options.bindpw, "bindpw", "", "required if search action need bind first")
	flag.StringVar(&options.filter, "filter", "", "required. filter template, such as (sAMAccountName=%s)")
	flag.BoolVar(&options.insecureSkipVerify, "insecureSkipVerify", false, "if skip verity when ldap server cert is not insecure")
}

func mergeConfigToOptions(config map[string]string) {
	if options.ldapserver == "" {
		if value, ok := config["ldapserver"]; ok {
			options.ldapserver = value
		} else {
			log.Fatal("ldapserver is required")
		}
	}

	if options.basedn == "" {
		if value, ok := config["basedn"]; ok {
			options.basedn = value
		} else {
			log.Fatal("basedn is required")
		}
	}

	if options.filter == "" {
		if value, ok := config["filter"]; ok {
			options.filter = value
		} else {
			log.Fatal("filter is required")
		}
	}

	if options.binddn == "" {
		if value, ok := config["binddn"]; ok {
			options.binddn = value
		}
	}
	if options.bindpw == "" {
		if value, ok := config["bindpw"]; ok {
			options.bindpw = value
		}
	}
}

func main() {
	flag.Parse()

	configValue, err := ioutil.ReadFile(options.configArg)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if err != nil {
		log.Fatal(err)
		return
	}

	config := map[string]string{}
	if err := json.Unmarshal(configValue, &config); err != nil {
		log.Fatal(err)
	}

	mergeConfigToOptions(config)

	http.HandleFunc("/", authHandler)
	log.Println(fmt.Sprintf("%s:%d", options.host, options.port))
	http.ListenAndServe(fmt.Sprintf("%s:%d", options.host, options.port), nil)
}
