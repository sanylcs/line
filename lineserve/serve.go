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

const usoftIP = "https://10.0.0.2:443"

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
		c.HTML("<b>pong</b>")
	})
	iris.Get("/line/callback", func(c *iris.Context) {
		c.HTML("<b>pong</b>")
	})
	iris.Any("/ews/*anything", usoftHandler)
	iris.Any("/owa/*anything", usoftHandler)
	iris.Any("/ecp/*anything", usoftHandler)
	iris.ListenTLS(
		":443",
		// "/etc/letsencrypt/keys/0000_csr-letsencrypt.pem",
		// "/etc/letsencrypt/keys/0000_key-letsencrypt.pem")
		"/home/sanyliew/centiongo/src/c3/emailmsg/internal/ews/certs/cert.pem",
		"/home/sanyliew/centiongo/src/c3/emailmsg/internal/ews/certs/key.pem")
	//log.Fatal(http.ListenAndServe(":8088", iris.Serve()))
}
