package main

import (
	"encoding/json"
	"maps"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParams(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, err
	}
	return id, nil
}

type envelop map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelop, headers http.Header) error {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	json = append(json, '\n')

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}
