package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"hawx.me/code/serve"
)

var (
	port   = flag.String("port", "8080", "")
	socket = flag.String("socket", "", "")
)

type packageConfig struct {
	Prefix, VCS, URL string
}
type Config map[string]packageConfig

func Server(host string, conf Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Println(r.URL)
		
		if r.URL.Path[len(r.URL.Path)-1] == '/' {
			http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusMovedPermanently)
			return
		}

		row, ok := conf[r.URL.Path]
		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if r.FormValue("go-get") != "1" {
			http.Redirect(w, r, "http://godoc.org/"+host+r.URL.Path, http.StatusFound)
			return
		}

		w.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(w, `<meta name="go-import" content="%s%s %s %s">`, host, row.Prefix, row.VCS, row.URL)
	})
}

func DecodeConfig(r io.Reader) (Config, error) {
	config := Config{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		switch len(fields) {
		case 0:
			continue
		case 3:
			config[fields[0]] = packageConfig{fields[0], fields[1], fields[2]}
		default:
			return nil, fmt.Errorf("config malformed: %s", scanner.Text())
		}
	}

	return config, nil
}

const usage = "Usage: vanity [--port PORT | --socket SOCK] HOST CONFIG\n"

func main() {
	flag.Parse()
	argv := flag.Args()

	if len(argv) != 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	configFile, err := os.Open(argv[1])
	if err != nil {
		log.Println(err)
		return
	}

	config, err := DecodeConfig(configFile)
	if err != nil {
		log.Println(err)
		return
	}

	serve.Serve(*port, *socket, Server(argv[0], config))
}
