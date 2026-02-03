package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type FuelData struct {
	// Вхідні дані
	Vygilya, Mazut, Gaz float64

	// Склад вугілля
	Hd, Cd, Sd, Nd, Od, Wd, Ad, Vd float64
	// Склад мазуту
	AGM float64

	// Результати розрахунків
	ResultEmVyg, ResultValVukudVyg float64
	ResultEmMazut, ResultValVukudMazut float64
	ResultValVukudGaz float64
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", formHandler)

	fmt.Println("Server started at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	data := FuelData{
		Vygilya: 858613.05,
		Mazut:   88993.41,
		Gaz:     104435.26,
		Hd:      3.50,
		Cd:      52.49,
		Sd:      2.85,
		Nd:      0.97,
		Od:      4.99,
		Wd:      10.00,
		Ad:      25.20,
		Vd:      25.92,
		AGM:     0.15,
	}

	if r.Method == "POST" {
		data.Vygilya, _ = strconv.ParseFloat(r.FormValue("Vygilya"), 64)
		data.Mazut, _ = strconv.ParseFloat(r.FormValue("Mazut"), 64)
		data.Gaz, _ = strconv.ParseFloat(r.FormValue("Gaz"), 64)

		data.ResultEmVyg = pokaznukEmVyg(data.Ad)
		data.ResultValVukudVyg = valVukudVyg(data.ResultEmVyg, data.Vygilya)

		data.ResultEmMazut = pokaznukEmMazut(data.AGM)
		data.ResultValVukudMazut = valVukudMazut(data.ResultEmMazut, data.Mazut)

		data.ResultValVukudGaz = valVukudGaz(data.Gaz)
	}

	tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl.Execute(w, data)
}

// Константи
const (
	teplotaZgoranyaPalyvo = 20.47
	ZolaPalyvo            = 0.8
	GorRechovynyPalyvo    = 1.5
	Efectyvnist           = 0.985
	ZolaMazut             = 1
	teplotaZgoranyaMazut  = 39.49
	GorRechovynMazut      = 0.0
	teplotaZgoranyaGaz    = 33.08
	pokaznukEmGaz         = 0
)

// Функції для розрахунків
func pokaznukEmVyg(a float64) float64 {
	return (math.Pow(10, 6) / teplotaZgoranyaPalyvo) * ZolaPalyvo * (a / (100 - GorRechovynyPalyvo)) * (1 - Efectyvnist)
}

func pokaznukEmMazut(a float64) float64 {
	return (math.Pow(10, 6) / teplotaZgoranyaMazut) * ZolaMazut * (a / (100 - GorRechovynMazut)) * (1 - Efectyvnist)
}

func valVukudVyg(em, vyg float64) float64 {
	return 1e-6 * em * teplotaZgoranyaPalyvo * vyg
}

func valVukudMazut(em, mazut float64) float64 {
	return 1e-6 * em * teplotaZgoranyaMazut * mazut
}

func valVukudGaz(gaz float64) float64 {
	return 1e-6 * pokaznukEmGaz * teplotaZgoranyaGaz * gaz
}
