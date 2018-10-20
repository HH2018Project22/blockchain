package main

import "flag"

var (
	host = "0.0.0.0"
	port = 8080
)

func init() {
	flag.StringVar(&host, "host", host, "HTTP host")
	flag.IntVar(&port, "port", port, "HTTP port")
}

func main() {

	flag.Parse()

}
