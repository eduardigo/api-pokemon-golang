package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Pokemon struct {
	Number     int      `json:"number"`
	Name       string   `json:"name"`
	Type1      string `json:"type1"`
	Type2      string `json:"type2"`
	Generation int      `json:"generation"`
	Weight     float32  `json:"weight"`
	Height     float32  `json:"height"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/create-pokemons", createPokemons).Methods("POST")
	router.HandleFunc("/list-pokemons", listPokemons).Methods("GET")
	router.HandleFunc("/get-pokemon-by-number/{number}", getPokemonByNumber).Methods("GET")
	router.HandleFunc("/get-pokemon-by-name/{name}", getPokemonByName).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("pokemons.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func createPokemons(w http.ResponseWriter, r *http.Request) {
	var pokemons []Pokemon
	db := connectDB()
	db.AutoMigrate(&Pokemon{})
	db.Create([]*Pokemon{
		&Pokemon{Number: 1, Name: "Bulbasaur", Type1: "Grass", Type2: "Poison", Generation: 1, Weight: 6.9, Height: 0.7},
		&Pokemon{Number: 2, Name: "Ivysaur", Type1: "Grass", Type2: "Poison", Generation: 1, Weight: 13.0, Height: 1.0},
		&Pokemon{Number: 3, Name: "Venusaur", Type1: "Grass", Type2: "Poison", Generation: 1, Weight: 100.0, Height: 2.0},
	})
	json.NewEncoder(w).Encode(pokemons)
}

func listPokemons(w http.ResponseWriter, r *http.Request) {
	var pokemons []Pokemon
	db := connectDB()
	db.Find(&pokemons)
	json.NewEncoder(w).Encode(pokemons)
}

func getPokemonByNumber(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var pokemon Pokemon
	db := connectDB()
	db.First(&pokemon, params["number"])
	json.NewEncoder(w).Encode(pokemon)
}

func getPokemonByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	var pokemon Pokemon
	db := connectDB()
	db.Where("name LIKE ?", "%"+name+"%").First(&pokemon)
	json.NewEncoder(w).Encode(pokemon)
}
