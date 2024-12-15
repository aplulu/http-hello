package main

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

//go:embed index.tmpl.html
var indexTmpl string

type config struct {
	Port           string   `envconfig:"port" default:"8080"`
	JWKSURL        string   `envconfig:"jwks_url" default:""`
	JWKHeaderName  string   `envconfig:"jwk_header_name" default:"Authorization"`
	Title          string   `envconfig:"title" default:"Congratulations"`
	SuccessMessage string   `envconfig:"success_message" default:"Successfully deployed."`
	CustomUserData []string `envconfig:"custom_data" default:"Issuer:.iss,Subject:.sub"`
}

type customUserDataItem struct {
	Label string
	Value string
}

func main() {
	serverCtx := context.Background()

	tmpl := template.Must(template.New("index").Parse(indexTmpl))

	var conf config
	if err := envconfig.Process("", &conf); err != nil {
		log.Fatalf("Failed to process env var: %v", err)
	}

	var ar *jwk.AutoRefresh
	var keySet jwk.Set
	if conf.JWKSURL != "" {
		ar = jwk.NewAutoRefresh(serverCtx)

		var err error
		keySet, err = ar.Fetch(serverCtx, conf.JWKSURL)
		if err != nil {
			log.Fatalf("Failed to fetch JWKS: %v", err)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"title":          conf.Title,
			"successMessage": conf.SuccessMessage,
		}

		if keySet != nil {
			authHeader := r.Header.Get(conf.JWKHeaderName)
			if authHeader != "" {
				http.Error(w, fmt.Sprintf("%s header required.", conf.JWKHeaderName), http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate JWT
			token, err := jwt.Parse(
				[]byte(tokenString),
				jwt.WithKeySet(keySet),
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims := make(map[string]interface{})
			for k, v := range token.PrivateClaims() {
				claims[k] = v
			}

			if sub, ok := token.Get("sub"); ok {
				claims["sub"] = sub
			}
			if iss, ok := token.Get("iss"); ok {
				claims["iss"] = iss
			}
			if exp, ok := token.Get("exp"); ok {
				claims["exp"] = exp
			}

			data["customUserData"] = make([]customUserDataItem, 0)
			for _, v := range conf.CustomUserData {
				s := strings.Split(v, ":")
				if len(s) != 2 {
					continue
				}
				if value, ok := claims[s[1]]; ok {
					data["customUserData"] = append(data["customUserData"].([]customUserDataItem), customUserDataItem{
						Label: s[0],
						Value: fmt.Sprintf("%v", value),
					})
				}
			}
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Printf("Starting server on :%s", conf.Port)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("", conf.Port), nil))
}
