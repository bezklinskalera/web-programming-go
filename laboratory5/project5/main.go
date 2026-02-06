package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Дані для шаблону
type PageData struct {
	// Вхідні дані
	ChVidEl110, ChVidPl110, ChVidT110, ChVidVV10, ChVidPr10       string
	TrVidEl110, TrVidPl110, TrVidT110, TrVidVV10, TrVidPr10       string
	ChVidSek                                                       string
	ChVid2, TrVid2, SerChas2                                       string
	ZbutkiAv, ZbutkiPl                                             string

	// Розраховані результати
	ChastotaVidOdnok, TryvVidOdnok, KAva, KPlan                   string
	ChDvaKola, ChDvaKolaSekVum, MatAv, MatPl, MatZb               string
}

// Перетворення рядка у float64, з fallback 0
func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	fmt.Println("Server started at http://localhost:8095")
	log.Fatal(http.ListenAndServe(":8095", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ChVidEl110: "0.01", ChVidPl110: "0.07", ChVidT110: "0.015", ChVidVV10: "0.02", ChVidPr10: "0.18",
		TrVidEl110: "30", TrVidPl110: "10", TrVidT110: "100", TrVidVV10: "15", TrVidPr10: "2",
		ChVidSek: "0.02",
		ChVid2: "0.01", TrVid2: "0.045", SerChas2: "0.004",
		ZbutkiAv: "23.6", ZbutkiPl: "17.6",
	}

	// Якщо це POST — беремо значення з форми
	if r.Method == http.MethodPost {
		data.ChVidEl110 = r.FormValue("ChVidEl110")
		data.ChVidPl110 = r.FormValue("ChVidPl110")
		data.ChVidT110 = r.FormValue("ChVidT110")
		data.ChVidVV10 = r.FormValue("ChVidVV10")
		data.ChVidPr10 = r.FormValue("ChVidPr10")

		data.TrVidEl110 = r.FormValue("TrVidEl110")
		data.TrVidPl110 = r.FormValue("TrVidPl110")
		data.TrVidT110 = r.FormValue("TrVidT110")
		data.TrVidVV10 = r.FormValue("TrVidVV10")
		data.TrVidPr10 = r.FormValue("TrVidPr10")

		data.ChVidSek = r.FormValue("ChVidSek")
		data.ChVid2 = r.FormValue("ChVid2")
		data.TrVid2 = r.FormValue("TrVid2")
		data.SerChas2 = r.FormValue("SerChas2")
		data.ZbutkiAv = r.FormValue("ZbutkiAv")
		data.ZbutkiPl = r.FormValue("ZbutkiPl")

		// Перетворюємо у float64 для обчислень
		chVidEl110 := parseFloat(data.ChVidEl110)
		chVidPl110 := parseFloat(data.ChVidPl110)
		chVidT110 := parseFloat(data.ChVidT110)
		chVidVV10 := parseFloat(data.ChVidVV10)
		chVidPr10 := parseFloat(data.ChVidPr10)

		trVidEl110 := parseFloat(data.TrVidEl110)
		trVidPl110 := parseFloat(data.TrVidPl110)
		trVidT110 := parseFloat(data.TrVidT110)
		trVidVV10 := parseFloat(data.TrVidVV10)
		trVidPr10 := parseFloat(data.TrVidPr10)

		chVidSek := parseFloat(data.ChVidSek)
		chVid2 := parseFloat(data.ChVid2)
		trVid2 := parseFloat(data.TrVid2)
		serChas2 := parseFloat(data.SerChas2)
		zbutkiAv := parseFloat(data.ZbutkiAv)
		zbutkiPl := parseFloat(data.ZbutkiPl)

		// Обчислення
		resultChastotaVidOdnok := chVidEl110 + chVidPl110 + chVidT110 + chVidVV10 + chVidPr10
		resultTryvVidOdnok := (chVidEl110*trVidEl110 + chVidPl110*trVidPl110 + chVidT110*trVidT110 +
			chVidVV10*trVidVV10 + chVidPr10*trVidPr10) / resultChastotaVidOdnok
		resultKAva := resultChastotaVidOdnok * resultTryvVidOdnok / 8760
		resultKPlan := 1.2 * (43.0 / 8760)
		resultChDvaKola := 2 * resultChastotaVidOdnok * (resultKAva + resultKPlan)
		resultChDvaKolaSekVum := resultChDvaKola + chVidSek
		resultMatAv := chVid2 * trVid2 * 5120 * 6451
		resultMatPl := serChas2 * 5120 * 6451
		resultMatZb := zbutkiAv*resultMatAv + zbutkiPl*resultMatPl

		// Форматування
		data.ChastotaVidOdnok = fmt.Sprintf("%.3f", resultChastotaVidOdnok)
		data.TryvVidOdnok = fmt.Sprintf("%.1f", resultTryvVidOdnok)
		data.KAva = fmt.Sprintf("%.5f", resultKAva)
		data.KPlan = fmt.Sprintf("%.5f", resultKPlan)
		data.ChDvaKola = fmt.Sprintf("%.5f", resultChDvaKola)
		data.ChDvaKolaSekVum = fmt.Sprintf("%.4f", resultChDvaKolaSekVum)
		data.MatAv = fmt.Sprintf("%.0f", resultMatAv)
		data.MatPl = fmt.Sprintf("%.0f", resultMatPl)
		data.MatZb = fmt.Sprintf("%.0f", resultMatZb)
	}

	tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl.Execute(w, data)
}
