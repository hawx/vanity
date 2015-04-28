package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/hawx/serve"
)

var (
	port   = flag.String("port", "8080", "")
	socket = flag.String("socket", "", "")
	host   = flag.String("host", "", "")
)

type configRow struct {
	Prefix, Vcs, Url string
}
type Config []configRow

func (c Config) Len() int           { return len(c) }
func (c Config) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Config) Less(i, j int) bool { return c[i].Prefix > c[j].Prefix }

func Server(host string, conf Config) http.Handler {
	sort.Sort(conf)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", 405)
			return
		}

		if host == "" {
			host = r.Host
		}

		log.Println(r.URL)

		if r.FormValue("go-get") != "1" {
			http.Redirect(w, r, "http://godoc.org/"+host+r.URL.Path, http.StatusFound)
			return
		}

		for _, row := range conf {
			if strings.HasPrefix(r.URL.Path, row.Prefix) {
				fmt.Fprintf(w, `<meta name="go-import" content="%s%s %s %s">`, host, row.Prefix, row.Vcs, row.Url)
				return
			}
		}
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
			config = append(config, configRow{fields[0], fields[1], fields[2]})
		default:
			return nil, fmt.Errorf("config malformed: %s", scanner.Text())
		}
	}

	return config, nil
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprint(os.Stderr, "Usage: vanity [--port PORT | --socket SOCK] CONFIG\n")
		os.Exit(1)
	}

	configFile, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	config, err := DecodeConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	serve.Serve(*port, *socket, Server(*host, config))
}
