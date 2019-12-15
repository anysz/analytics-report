package main

import (
	"fmt"
	"log"
	"flag"
	"net/http"
	"github.com/anysz/analytics-report/storage"
)

func main() {
	var addr *string
	// flag
	flag.String(addr, "addr", "127.0.0.1:8081", "[host]:[port]")
	flag.Parse()


	// data storage
	storage.NewFileStorage("meta.json")

	status := storage.Driver.InitConfig()
	if status.IsNew {
		log.Println("Log created.")
	} else {
		log.Println("Log loaded.")
	}


	// router
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}

		if r, e := storage.Driver.GetAllLogDataIds(); e {
			fmt.Fprintf(w,"%#+v\n", r)
		}else{
			fmt.Fprintf(w, "{\"reason\":\"internal_error\"}")
		}
	})


	// init
	if addr != nil {
		 if err := http.ListenAndServe(*addr, mux); err != nil {
		 	panic(err)
		 }
	} else {
		log.Println("Addr is not preserved. (nil)")
	}
}