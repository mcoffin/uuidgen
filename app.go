package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"code.google.com/p/go-uuid/uuid"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const (
	CountKey = "n"
)

type UUID uuid.UUID

func (self UUID) MarshalJSON() ([]byte, error) {
	realUUID := uuid.UUID(self)
	return []byte(fmt.Sprintf("\"%s\"", realUUID.String())), nil
}

func Generate(res http.ResponseWriter, req *http.Request) {
	nStr := req.PostFormValue(CountKey)
	n, err := strconv.Atoi(nStr)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, "[")
	for i := 0; i < n; i++ {
		u := uuid.NewRandom()
		data, _ := json.Marshal(UUID(u))
		res.Write(data)
		if i < n-1 {
			fmt.Fprint(res, ",")
		}
	}
	fmt.Fprint(res, "]")
}

func main() {
	bind := flag.String("http", ":8080", "http binding")
	flag.Parse()

	fs := http.FileServer(http.Dir("static"))

	r := mux.NewRouter()
	r.Handle("/", fs)
	r.HandleFunc("/uuids", Generate).Methods("POST")

	n := negroni.Classic()
	n.UseHandler(r)
	http.Handle("/", n)
	http.ListenAndServe(*bind, nil)
}
