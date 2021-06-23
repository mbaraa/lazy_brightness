package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	backlight = "/sys/class/backlight/"
	intel     = "intel_backlight/"
)

var (
	maxBrightnessIntel = getMaxBrits()
	currentBritsIntel  = getCurrentBrits()
)

func getMaxBrits() float32 {
	maxB, _ := os.ReadFile(path.Join(backlight, intel, "max_brightness"))
	maxBInt, _ := strconv.ParseInt(strings.Trim(string(maxB), "\r \n"), 10, 64)
	return float32(maxBInt)
}

func getCurrentBrits() float32 {
	currB, _ := os.ReadFile(path.Join(backlight, intel, "brightness"))
	currBInt, _ := strconv.ParseInt(strings.Trim(string(currB), "\r \n"), 10, 64)
	return float32(currBInt)
}

func decBrits() {
	_5britSteps := maxBrightnessIntel / 50.
	currentBritsIntel = currentBritsIntel - _5britSteps

	newBrit := fmt.Sprint(currentBritsIntel)
	fmt.Println("fuck dec: ", currentBritsIntel)

	err := os.WriteFile(path.Join(backlight, intel, "brightness"), []byte(newBrit), 0644)
	if err != nil {
		fmt.Println("err: ", err)
	}
	currentBritsIntel = getCurrentBrits()
}

func incBrits() {
	_5britSteps := maxBrightnessIntel / 50.
	currentBritsIntel = currentBritsIntel + _5britSteps
	//if currentBritsIntel >= maxBrightnessIntel {
	//	return
	//}

	newBrit := fmt.Sprint(currentBritsIntel)

	fmt.Println("fuck inc: ", currentBritsIntel)

	err := os.WriteFile(path.Join(backlight, intel, "brightness"), []byte(newBrit), 0644)
	if err != nil {
		fmt.Println("err: ", err)
	}
	currentBritsIntel = getCurrentBrits()
}

func brits(brit int) {
	newBrit := (int(maxBrightnessIntel) * brit) / 100
	os.WriteFile(path.Join(backlight, intel, "brightness"), []byte(strconv.Itoa(newBrit)), 0664)
}

// GET /brits/inc
func HandleIncBrightness(res http.ResponseWriter, req *http.Request) {
	incBrits()
}

// GET /brits/dec
func HandleDecBrightness(res http.ResponseWriter, req *http.Request) {
	decBrits()
}

// GET /brits/get
func HandleGetBrightness(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(map[string]interface{}{
		"brightness": getCurrentBrits(),
		"max":        getMaxBrits(),
	})
}

// GET /brits/set?b=someLevel
func HandleSetBrightness(res http.ResponseWriter, req *http.Request) {
	if level, ok := req.URL.Query()["b"]; ok {
		levelInt, _ := strconv.Atoi(level[0])
		if levelInt < 0 || levelInt > 100 {
			res.WriteHeader(400)
			return
		}

		brits(levelInt)

		return
	}
	res.WriteHeader(400)
}

var endPoints = map[string]http.HandlerFunc{
	"GET /dec": HandleDecBrightness,
	"GET /inc": HandleIncBrightness,
	"GET /get": HandleGetBrightness,
	"GET /set": HandleSetBrightness,
}

func getHandler(method, path string) http.HandlerFunc {
	fmt.Println("req: ", method+" "+strings.TrimPrefix(path, "/brits"))
	if hf, ok := endPoints[method+" "+strings.TrimPrefix(path, "/brits")]; ok {
		return hf
	}
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("cyka nahui blyat!"))
	}
}

func main() {
	//fmt.Println(exec.Command("bash", "xbacklight -dec 50").Run())

	//fmt.Println(getCurrentBrits())
	//for {
	//
	//	incBrits()
	//	fmt.Println("brits: ", currentBritsIntel)
	//	fmt.Println("max: " ,maxBrightnessIntel)
	//	var i int
	//	fmt.Scanf("%d", &i)
	//}

	http.HandleFunc("/brits/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		handler := getHandler(req.Method, req.URL.Path)
		handler(res, req)
	})
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)
}
