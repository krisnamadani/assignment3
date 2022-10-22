package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Data struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Status struct {
	Status      Data   `json:"status"`
	StatusWater string `json:"status_water"`
	StatusWind  string `json:"status_wind"`
}

func task() {
	for {
		var status Status

		var min int = 1
		var max int = 100

		status.Status.Water = rand.Intn(max-min) + min
		status.Status.Wind = rand.Intn(max-min) + min

		file_json, _ := os.Create("status.json")

		byteValue, _ := json.Marshal(status)

		file_json.Write(byteValue)

		time.Sleep(15 * time.Second)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("template.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl.Execute(w, nil)
}

func status(w http.ResponseWriter, r *http.Request) {
	var status Status

	file_json, _ := os.Open("status.json")
	defer file_json.Close()
	byte_value, _ := ioutil.ReadAll(file_json)
	json.Unmarshal(byte_value, &status)

	check_water := status.Status.Water
	check_wind := status.Status.Wind

	var status_water, status_wind string

	switch {
	case check_water < 5:
		status_water = "Aman"
	case check_water >= 6 && check_water <= 8:
		status_water = "Siaga"
	case check_water > 8:
		status_water = "Bahaya"
	}

	switch {
	case check_wind < 6:
		status_wind = "Aman"
	case check_wind >= 7 && check_wind <= 15:
		status_wind = "Siaga"
	case check_wind > 15:
		status_wind = "Bahaya"
	}

	status.StatusWater = status_water
	status.StatusWind = status_wind

	json.NewEncoder(w).Encode(status)
}

func main() {
	go task()

	http.HandleFunc("/", index)
	http.HandleFunc("/status", status)

	fmt.Println("server berjalan di http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
