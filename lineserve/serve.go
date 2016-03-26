/*
This binary is to serve callback url for LINE authorization.
It also reverse proxy to Microsoft Exchange server.
*/
package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/kataras/iris"
)

const usoftIP = "https://127.0.0.1:444"

var usoftES *httputil.ReverseProxy

func serve(rw http.ResponseWriter, req *http.Request) {
	usoftES.ServeHTTP(rw, req)
}

func initMicrosoftExchangeServerProxy() {
	url, _ := url.Parse(usoftIP)
	usoftES = httputil.NewSingleHostReverseProxy(url)
	usoftES.Transport = &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout:   time.Minute,
		ResponseHeaderTimeout: time.Minute,
		//ExpectContinueTimeout: time.Minute,
		DisableCompression: true,
	}
}

func main() {
	initMicrosoftExchangeServerProxy()
	usoftHandler := iris.ToHandlerFunc(usoftES.ServeHTTP)
	iris.Get("/ping", func(c *iris.Context) {
		c.HTML("<b>pong pong</b>")
	})
	iris.Get("/line/callback", func(c *iris.Context) {
		c.HTML("<p><b>token:</b></p><p><b>secret:</b></p>")
	})
	iris.Get("/ews/*anything", usoftHandler)
	iris.Get("/owa/*anything", usoftHandler)
	iris.Get("/ecp/*anything", usoftHandler)
	iris.Post("/ews/*anything", usoftHandler)
	iris.Post("/owa/*anything", usoftHandler)
	iris.Post("/ecp/*anything", usoftHandler)
	iris.Head("/ews/*anything", usoftHandler)
	iris.Head("/owa/*anything", usoftHandler)
	iris.Head("/ecp/*anything", usoftHandler)
	iris.Delete("/ews/*anything", usoftHandler)
	iris.Delete("/owa/*anything", usoftHandler)
	iris.Delete("/ecp/*anything", usoftHandler)
	iris.Put("/ews/*anything", usoftHandler)
	iris.Put("/owa/*anything", usoftHandler)
	iris.Put("/ecp/*anything", usoftHandler)
	iris.Connect("/ews/*anything", usoftHandler)
	iris.Connect("/owa/*anything", usoftHandler)
	iris.Connect("/ecp/*anything", usoftHandler)
	iris.Get("/ews", usoftHandler)
	iris.Get("/owa", usoftHandler)
	iris.Get("/ecp", usoftHandler)
	iris.Post("/ews", usoftHandler)
	iris.Post("/owa", usoftHandler)
	iris.Post("/ecp", usoftHandler)
	iris.Head("/ews", usoftHandler)
	iris.Head("/owa", usoftHandler)
	iris.Head("/ecp", usoftHandler)
	iris.Delete("/ews", usoftHandler)
	iris.Delete("/owa", usoftHandler)
	iris.Delete("/ecp", usoftHandler)
	iris.Put("/ews", usoftHandler)
	iris.Put("/owa", usoftHandler)
	iris.Put("/ecp", usoftHandler)
	iris.Connect("/ews", usoftHandler)
	iris.Connect("/owa", usoftHandler)
	iris.Connect("/ecp", usoftHandler)
	iris.Get("/Autodiscover/*anything", usoftHandler)
	iris.Post("/Autodiscover/*anythingws", usoftHandler)
	iris.Delete("/Autodiscover/*anything", usoftHandler)
	iris.Put("/Autodiscover/*anything", usoftHandler)
	iris.Connect("/Autodiscover/*anything", usoftHandler)
	iris.Get("/Autodiscover", usoftHandler)
	iris.Post("/Autodiscover", usoftHandler)
	iris.Delete("/Autodiscover", usoftHandler)
	iris.Put("/Autodiscover", usoftHandler)
	iris.Connect("/Autodiscover", usoftHandler)
	iris.ListenTLS(
		":443",
		"certs/c1/fullchain1.pem",
		"certs/c1/privkey1.pem")
	//		"certs/c2/fullchain1.pem",
	//		"certs/c2/privkey.pem")
	//log.Fatal(http.ListenAndServe(":8088", iris.Serve()))
}
