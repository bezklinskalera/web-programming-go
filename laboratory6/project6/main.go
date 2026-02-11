package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type SolarData struct {
	// Вхідні дані
	SerDPot, SerKvadrVid, SerKvadrVidZmen, Vartist float64

	// Результати обчислень
	ChastkaEn, W1, Prub1, W2, Sh1                 float64
	ChastkaEn2, W3, Prub2, W4, Sh2, GPrub        float64
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", formHandler)

	fmt.Println("Server started at http://localhost:8091")
	http.ListenAndServe(":8091", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	data := SolarData{
		SerDPot:        5,
		SerKvadrVid:    1,
		SerKvadrVidZmen: 0.25,
		Vartist:        7,
	}

	if r.Method == "POST" {
		data.SerDPot, _ = strconv.ParseFloat(r.FormValue("SerDPot"), 64)
		data.SerKvadrVid, _ = strconv.ParseFloat(r.FormValue("SerKvadrVid"), 64)
		data.SerKvadrVidZmen, _ = strconv.ParseFloat(r.FormValue("SerKvadrVidZmen"), 64)
		data.Vartist, _ = strconv.ParseFloat(r.FormValue("Vartist"), 64)

		// Обчислення
		data.ChastkaEn = chastkaEn(data.SerDPot, data.SerKvadrVid)
		data.W1 = w1(data.SerDPot, data.SerKvadrVid) * 0.01
		data.W2 = w2(data.SerDPot, data.SerKvadrVid) * 0.01
		data.Prub1 = data.W1 * data.Vartist
		data.Sh1 = data.W2 * data.Vartist

		data.ChastkaEn2 = chastkaEn(data.SerDPot, data.SerKvadrVidZmen)
		data.W3 = w1(data.SerDPot, data.SerKvadrVidZmen) * 0.01
		data.W4 = w2(data.SerDPot, data.SerKvadrVidZmen) * 0.01
		data.Prub2 = data.W3 * data.Vartist
		data.Sh2 = data.W4 * data.Vartist
		data.GPrub = data.Prub2 - data.Sh2
	}

	tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl.Execute(w, data)
}

// Функції розрахунків
func w1(vSerDPot, vSerKvadrVid float64) float64 {
	return vSerDPot * 24 * chastkaEn(vSerDPot, vSerKvadrVid)
}

func w2(vSerDPot, vSerKvadrVid float64) float64 {
	return vSerDPot * 24 * (100 - chastkaEn(vSerDPot, vSerKvadrVid))
}

func chastkaEn(vSerDPot, vSerKvadrVid float64) float64 {
	mean := vSerDPot
	stddev := vSerKvadrVid
	lowerBound := 4.75
	upperBound := 5.25

	return 0.5 * (erf((upperBound-mean)/(math.Sqrt(2)*stddev)) - erf((lowerBound-mean)/(math.Sqrt(2)*stddev)))
}

// Реалізація ерф функції
func erf(x float64) float64 {
	a1 := 0.254829592
	a2 := -0.284496736
	a3 := 1.421413741
	a4 := -1.453152027
	a5 := 1.061405429
	p := 0.3275911

	sign := 1.0
	if x < 0 {
		sign = -1
	}
	x = math.Abs(x)

	t := 1.0 / (1.0 + p*x)
	y := 1.0 - ((((a5*t+a4)*t+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)

	return sign * y * 100
}
