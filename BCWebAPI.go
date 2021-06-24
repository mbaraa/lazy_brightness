package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type BCWebAPI struct {
	endPoints map[string]http.HandlerFunc
	bc        *BrightnessController
}

func NewBCWebAPI(bc *BrightnessController) *BCWebAPI {
	return (&BCWebAPI{bc: bc}).initEndPoints()
}

func (bapi *BCWebAPI) initEndPoints() *BCWebAPI {
	bapi.endPoints = map[string]http.HandlerFunc{
		"GET /dec": bapi.handleDecBrightness,
		"GET /inc": bapi.handleIncBrightness,
		"GET /get": bapi.handleGetBrightness,
		"GET /set": bapi.handleSetBrightness,
	}
	return bapi
}

// ServeHTTP responds with a valid handler otherwise responses with a 404 status code
//
func (bapi *BCWebAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if handler, ok := bapi.endPoints[req.Method+" "+strings.TrimPrefix(req.URL.Path, "/brits")]; ok {
		handler(res, req)
		return
	}
	if req.Method != "OPTIONS" {
		res.WriteHeader(404)
	}
}

// GET /brits/inc
func (bapi *BCWebAPI) handleIncBrightness(res http.ResponseWriter, req *http.Request) {
	bapi.bc.IncBrits()
}

// GET /brits/dec
func (bapi *BCWebAPI) handleDecBrightness(res http.ResponseWriter, req *http.Request) {
	bapi.bc.DecBrits()
}

// GET /brits/get
func (bapi *BCWebAPI) handleGetBrightness(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(map[string]interface{}{
		"brightness": bapi.bc.GetCurrentBrits(),
		"max":        bapi.bc.GetMaxBrits(),
	})
}

// GET /brits/set?b=someLevel ∋ [0,100]
func (bapi *BCWebAPI) handleSetBrightness(res http.ResponseWriter, req *http.Request) {
	if level, ok := req.URL.Query()["b"]; ok {
		levelInt, _ := strconv.Atoi(level[0])
		if levelInt < 0 || levelInt > 100 {
			res.WriteHeader(400)
			return
		}

		bapi.bc.SetBrightness(levelInt)

		return
	}
	res.WriteHeader(400)
}
