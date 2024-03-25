package src

import (
	goproxy "github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var stripHeaders = []string{
	"X-Proxy-Login",
	"X-Proxy-Password",
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

const TorDefaultGateway = "http://127.0.0.1:9050"

func GetProxy() *goproxy.ProxyHttpServer {
	var validatorFunc ValidatorProxyFunc = EmptyValidator
	var ok bool

	validator := GetConfig("VALIDATOR")
	if validator != "" && validator != "no" {
		if validatorFunc, ok = ValidatorsAvailable[validator]; !ok {
			log.Fatalln("Validator " + validator + " not found")
		} else {
			log.Println("Using validator " + validator)
		}
	}

	proxy := goproxy.New()

	topLevelProxy := GetConfig("PROXY")
	if GetConfigBool("TOR") {
		log.Println("Using Tor as proxy")
		topLevelProxy = TorDefaultGateway
	}

	if topLevelProxy != "" {
		log.Println("Using top level proxy " + topLevelProxy)
		proxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
			return url.Parse(topLevelProxy)
		}
		proxy.ConnectDial = proxy.NewConnectDialToProxy(topLevelProxy)
	}

	proxy.OnRequest().DoFunc(handleHTTP)

	var handler goproxy.FuncHttpsHandler = func(r *http.Request, host string) (*http.Request, *goproxy.ConnectAction, string) {
		login, password := parseCredentials(r)

		for _, header := range stripHeaders {
			r.Header.Del(header)
		}

		// Check credentials via validator
		allowed, _ := validatorFunc(login, password)
		if !allowed {
			log.Println("Bad credentials for " + host)
			return r, goproxy.RejectConnect, host
		}

		// Separate host and get necessary proxy
		for domain, target := range LuaDomainRoutes {
			if domain != host {
				continue
			}
			if target == "TOR" {
				target = TorDefaultGateway
			}

			proxy.ConnectDial = proxy.NewConnectDialToProxy(target)
		}

		if len(LuaAllowedPorts) > 0 {
			ok = false

			for _, port := range LuaAllowedPorts {
				if strings.TrimSpace(port) == r.URL.Port() {
					ok = true
					break
				}
			}

			if !ok {
				log.Println("Port disallowed: " + r.URL.Port())
				return r, goproxy.RejectConnect, host
			}
		}

		for user, ports := range LuaUserPorts {
			for _, port := range ports {
				if strings.TrimSpace(port) != r.URL.Port() {
					continue
				}
				if user == login {
					continue
				}

				log.Println("User " + user + " cannot access port " + port)
				return r, goproxy.RejectConnect, host
			}
		}

		return r, nil, host
	}

	proxy.OnRequest().HandleConnect(handler)
	return proxy
}

func handleHTTP(r *http.Request) (*http.Request, *http.Response) {
	// Only https scheme is allowed, rewrite another
	if r.URL.Scheme != "http" && r.URL.Scheme != "https" {
		redirect := "https://" + r.URL.Host + r.URL.Path
		response := goproxy.NewResponse(r,
			goproxy.ContentTypeText, http.StatusTemporaryRedirect,
			"Protocol "+r.URL.Scheme+" is not enabled")
		response.Header.Set("Location", redirect)
		return r, response
	}

	return r, nil
}

func parseCredentials(r *http.Request) (string, string) {
	login, password := r.Header.Get("X-Proxy-Login"), r.Header.Get("X-Proxy-Password")
	if login != "" {
		return login, password
	}

	auth := strings.Split(r.Header.Get("Proxy-Authorization"), " ")
	if len(auth) >= 2 {
		cred := strings.Split(auth[1], ":")
		if len(cred) >= 2 {
			return cred[0], cred[1]
		}

		return auth[1], ""
	}

	return "", ""
}
