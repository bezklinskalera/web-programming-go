package main

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type PageData struct {
	// Завдання 1
	StrymKZ, Napryga, FictTimeKZ, PotTP, RozNav, Tm float64
	RozStrymNormAv, RozStrymAv, EcoPerer, Ss       float64

	// Завдання 2
	PotKZ2, Napruga2, SNomt2 float64
	XC, XT, SumaOpir, PochStrym float64

	// Завдання 3
	Umax3, UVn3, RSn3, XSn3, RSmin3, XSmin3 float64
	XT3, XSh3, XShmin3, ZSh3, ZShmin3 float64
	ISH3, ISH23, ISHmin3, ISHmin23, KPr float64
	RShn, XShn, ZShn, RShnmin, XShnmin, ZShnmin float64
	IShn3, IShn23, IShnmin3, IShnmin23 float64
	Ll3, Rl0, Xl0, Rl3, Xl3, RSumn, XSumn, ZSumn float64
	RSumminn, XSumminn, ZSumminn float64
	Ln3, Ln23, Lnmin3, Lnmin23 float64

	Calculated bool
}

func main() {
	tmpl := template.Must(template.ParseFiles("template.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			StrymKZ:  2.5, Napryga: 10, FictTimeKZ: 2.5, PotTP: 2000, RozNav: 1300, Tm: 4000,
			PotKZ2:   200, Napruga2: 10.5, SNomt2: 6.3,
			Umax3:    11.1, UVn3: 115, RSn3: 10.65, XSn3: 24.02, RSmin3: 34.88, XSmin3: 65.68,
			Ll3:      12.37, Rl0: 0.64, Xl0: 0.363,
		}

		if r.Method == http.MethodPost {
			r.ParseForm()

			// Завдання 1
			data.StrymKZ = parseFormFloat(r, "strymKZ")
			data.Napryga = parseFormFloat(r, "napryga")
			data.FictTimeKZ = parseFormFloat(r, "fictTimeKZ")
			data.PotTP = parseFormFloat(r, "potTP")
			data.RozNav = parseFormFloat(r, "rozNav")
			data.Tm = parseFormFloat(r, "Tm")

			data.RozStrymNormAv = rozStrymNormAv(data.RozNav, data.Napryga)
			data.RozStrymAv = rozStrymAv(data.RozNav, data.Napryga)
			data.Ss = ss(data.StrymKZ, data.FictTimeKZ)
			data.EcoPerer = ecoPerer(data.RozStrymNormAv)

			// Завдання 2
			data.Napruga2 = parseFormFloat(r, "napruga2")
			data.PotKZ2 = parseFormFloat(r, "potKZ2")
			data.SNomt2 = parseFormFloat(r, "sNomt2")

			data.XC = OpirXc(data.Napruga2, data.PotKZ2)
			data.XT = OpirXt(data.Napruga2, data.SNomt2)
			data.SumaOpir = sumaOpir(data.Napruga2, data.PotKZ2, data.SNomt2)
			data.PochStrym = pochStrym(data.Napruga2, data.PotKZ2, data.SNomt2)

			// Завдання 3
			data.Umax3 = parseFormFloat(r, "umax3")
			data.UVn3 = parseFormFloat(r, "uVn3")
			data.RSn3 = parseFormFloat(r, "rSn3")
			data.XSn3 = parseFormFloat(r, "xSn3")
			data.RSmin3 = parseFormFloat(r, "rSmin3")
			data.XSmin3 = parseFormFloat(r, "xSmin3")

			data.XT3 = reaktOpir(data.Umax3, data.UVn3, data.SNomt2)
			data.XSh3 = xSh3(data.XSn3, data.Umax3, data.UVn3, data.SNomt2)
			data.XShmin3 = xShmin3(data.XSmin3, data.Umax3, data.UVn3, data.SNomt2)
			data.ZSh3 = zSh3(data.RSn3, data.XSn3, data.Umax3, data.UVn3, data.SNomt2)
			data.ZShmin3 = zShmin3(data.RSmin3, data.XSmin3, data.Umax3, data.UVn3, data.SNomt2)
			data.ISH3 = iSh3(data.UVn3, data.RSn3, data.XSn3, data.Umax3, data.UVn3, data.SNomt2)
			data.ISH23 = iSh23(data.UVn3, data.RSn3, data.XSn3, data.Umax3, data.UVn3, data.SNomt2)
			data.ISHmin3 = iSHmin3(data.UVn3, data.RSmin3, data.XSmin3, data.Umax3, data.UVn3, data.SNomt2)
			data.ISHmin23 = iSHmin23(data.UVn3, data.RSmin3, data.XSmin3, data.Umax3, data.UVn3, data.SNomt2)
			data.KPr = kPr(data.Umax3, data.UVn3)

			// Тут можна додати всі інші обчислення 3-го завдання...
			data.Calculated = true
		}

		tmpl.Execute(w, data)
	})

	port := "8094"
	log.Printf("Server started! Open in browser: http://localhost:%s/\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func parseFormFloat(r *http.Request, key string) float64 {
	v := r.FormValue(key)
	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return val
}

// --- Функції розрахунків (всі JS -> Go) ---
func rozStrymNormAv(rozNav, napryga float64) float64 {
	return (rozNav / 2) / (math.Sqrt(3) * napryga)
}

func rozStrymAv(rozNav, napryga float64) float64 {
	return 2 * ((rozNav / 2) / (math.Sqrt(3) * napryga))
}

func ecoPerer(rozStrymNormAv float64) float64 {
	return rozStrymNormAv / 1.4
}

func ss(strymKZ, fictTimeKZ float64) float64 {
	return (strymKZ * math.Sqrt(fictTimeKZ)) / 92
}

func OpirXc(napruga2, potKZ2 float64) float64 {
	return (napruga2 * napruga2) / potKZ2
}

func OpirXt(napruga2, sNomt2 float64) float64 {
	return (napruga2 / 100) * (napruga2 * napruga2) / sNomt2
}

func sumaOpir(napruga2, potKZ2, sNomt2 float64) float64 {
	return OpirXc(napruga2, potKZ2) + OpirXt(napruga2, sNomt2)
}

func pochStrym(napruga2, potKZ2, sNomt2 float64) float64 {
	return napruga2 / (math.Sqrt(3) * sumaOpir(napruga2, potKZ2, sNomt2))
}

func reaktOpir(umax3, uVn3, sNomt2 float64) float64 {
	return (umax3 * uVn3 * uVn3) / (100 * sNomt2)
}

func xSh3(xSn3, umax3, uVn3, sNomt2 float64) float64 {
	return xSn3 + reaktOpir(umax3, uVn3, sNomt2)
}

func xShmin3(xSmin3, umax3, uVn3, sNomt2 float64) float64 {
	return xSmin3 + reaktOpir(umax3, uVn3, sNomt2)
}

func zSh3(rSh3, xSn3, umax3, uVn3, sNomt2 float64) float64 {
	x := xSh3(xSn3, umax3, uVn3, sNomt2)
	return math.Sqrt(rSh3*rSh3 + x*x)
}

func zShmin3(rSmin3, xSmin3, umax3, uVn3, sNomt2 float64) float64 {
	x := xShmin3(xSmin3, umax3, uVn3, sNomt2)
	return math.Sqrt(rSmin3*rSmin3 + x*x)
}

func iSh3(vUVn3, rSh3, xSn3, umax3, uVn3, sNomt2 float64) float64 {
	return (vUVn3 * 1000) / (math.Sqrt(3) * zSh3(rSh3, xSn3, umax3, uVn3, sNomt2))
}

func iSh23(vUVn3, rSh3, xSn3, umax3, uVn3, sNomt2 float64) float64 {
	return iSh3(vUVn3, rSh3, xSn3, umax3, uVn3, sNomt2) * math.Sqrt(3) / 2
}

func iSHmin3(vUVn3, rSmin3, xSmin3, umax3, uVn3, sNomt2 float64) float64 {
	return (vUVn3 * 1000) / (math.Sqrt(3) * zShmin3(rSmin3, xSmin3, umax3, uVn3, sNomt2))
}

func iSHmin23(vUVn3, rSmin3, xSmin3, umax3, uVn3, sNomt2 float64) float64 {
	return iSHmin3(vUVn3, rSmin3, xSmin3, umax3, uVn3, sNomt2) * math.Sqrt(3) / 2
}

func kPr(vUmax3, vUVn3 float64) float64 {
	return (vUmax3 * vUmax3) / (vUVn3 * vUVn3)
}
