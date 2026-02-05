package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Структура для електроприймачів
type EP struct {
	N       int
	Pn      float64
	PnSum   float64
	Kv      float64
	PnKv    float64
	Tg      float64
	PnKvTg  float64
	Pn2     float64
	GroupI  float64
}

// Дані для шаблону
type FuelData struct {
	EPs        []EP
	SumPn      float64
	SumPnKv    float64
	SumPnKvTg  float64
	SumPn2     float64
	Ne         float64
	Kp         float64
	Pp         float64
	Qp         float64
	Sp         float64
	Ip         float64
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// реєструємо функцію add1 для шаблону
	funcMap := template.FuncMap{
		"add1": func(i int) int { return i + 1 },
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := FuelData{
			EPs: []EP{
				{N: 4, Pn: 20, Kv: 0.15, Tg: 1.33},
				{N: 2, Pn: 14, Kv: 0.12, Tg: 1},
				{N: 1, Pn: 36, Kv: 0.3, Tg: 1.52},
			},
		}

		// Підсумкові розрахунки
		for i := range data.EPs {
			data.EPs[i].PnSum = float64(data.EPs[i].N) * data.EPs[i].Pn
			data.EPs[i].PnKv = data.EPs[i].PnSum * data.EPs[i].Kv
			data.EPs[i].PnKvTg = data.EPs[i].PnKv * data.EPs[i].Tg
			data.EPs[i].Pn2 = float64(data.EPs[i].N) * data.EPs[i].Pn * data.EPs[i].Pn
			data.EPs[i].GroupI = data.EPs[i].PnKvTg / 0.38 // приклад
			data.SumPn += data.EPs[i].PnSum
			data.SumPnKv += data.EPs[i].PnKv
			data.SumPnKvTg += data.EPs[i].PnKvTg
			data.SumPn2 += data.EPs[i].Pn2
		}

		data.Ne = 10
		data.Kp = 0.85
		data.Pp = 100
		data.Qp = 50
		data.Sp = 120
		data.Ip = 180

		tmpl := template.Must(template.New("template.html").Funcs(funcMap).ParseFiles("template.html"))
		tmpl.Execute(w, data)
	})

	fmt.Println("Server started at http://localhost:8091")
	http.ListenAndServe(":8091", nil)
}
