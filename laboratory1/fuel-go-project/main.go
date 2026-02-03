package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type FuelData struct {
	// Завдання 1
	H, C, S, N, O, W, A float64
	Krs, Krg             float64
	Hd, Cd, Sd, Nd, Od, Ad float64
	Hgor, Cgor, Sgor, Ngor, Ogor float64
	NRob, NSyha, NGor float64

	// Завдання 2
	C2, H2, O2, S2, W2, A2, V2 float64
	ResultC2, ResultH2, ResultO2, ResultS2, ResultW2, ResultA2, ResultV2 float64
	ResultTeplota2 float64
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Головна сторінка
	http.HandleFunc("/", formHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	data := FuelData{}

	if r.Method == "POST" {
		// Завдання 1
		data.H, _ = strconv.ParseFloat(r.FormValue("h"), 64)
		data.C, _ = strconv.ParseFloat(r.FormValue("c"), 64)
		data.S, _ = strconv.ParseFloat(r.FormValue("s"), 64)
		data.N, _ = strconv.ParseFloat(r.FormValue("n"), 64)
		data.O, _ = strconv.ParseFloat(r.FormValue("o"), 64)
		data.W, _ = strconv.ParseFloat(r.FormValue("w"), 64)
		data.A, _ = strconv.ParseFloat(r.FormValue("a"), 64)

		data.Krs = 100 / (100 - data.W)
		data.Krg = 100 / (100 - data.W - data.A)

		data.Hd = data.Krs * data.H
		data.Cd = data.Krs * data.C
		data.Sd = data.Krs * data.S
		data.Nd = data.Krs * data.N
		data.Od = data.Krs * data.O
		data.Ad = data.Krs * data.A

		data.Hgor = data.Krg * data.H
		data.Cgor = data.Krg * data.C
		data.Sgor = data.Krg * data.S
		data.Ngor = data.Krg * data.N
		data.Ogor = data.Krg * data.O

		data.NRob = (339*data.C + 1030*data.H - 108.8*(data.O-data.S) - 25*data.W) / 1000
		data.NSyha = (data.NRob + 0.025*data.W) * 100 / (100 - data.W)
		data.NGor = (data.NRob + 0.025*data.W) * 100 / (100 - data.W - data.A)

		// Завдання 2
		data.C2, _ = strconv.ParseFloat(r.FormValue("c2"), 64)
		data.H2, _ = strconv.ParseFloat(r.FormValue("h2"), 64)
		data.O2, _ = strconv.ParseFloat(r.FormValue("o2"), 64)
		data.S2, _ = strconv.ParseFloat(r.FormValue("s2"), 64)
		data.W2, _ = strconv.ParseFloat(r.FormValue("w2"), 64)
		data.A2, _ = strconv.ParseFloat(r.FormValue("a2"), 64)
		data.V2, _ = strconv.ParseFloat(r.FormValue("v2"), 64)

		data.ResultC2 = data.C2 * (100 - data.W2 - data.A2) / 100
		data.ResultH2 = data.H2 * (100 - data.W2 - data.A2) / 100
		data.ResultO2 = data.O2 * (100 - data.W2 - data.A2) / 100
		data.ResultS2 = data.S2 * (100 - data.W2 - data.A2) / 100
		data.ResultA2 = data.A2 * (100 - data.W2) / 100
		data.ResultV2 = data.V2 * (100 - data.W2) / 100

		teplota := 40.4
		data.ResultTeplota2 = teplota * ((100 - data.W2 - data.A2) / 100) - 0.025*data.W2
	}

	tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl.Execute(w, data)
}
