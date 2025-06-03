package main

import (
	"log"
	"net/http"
	"tbcvclub/configs"
	"tbcvclub/internal/database"
	"tbcvclub/internal/handlers"
	"tbcvclub/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialisation de la configuration
	if err := configs.LoadConfig(); err != nil {
		log.Fatalf("Erreur chargement config: %v", err)
	}

	// Initialisation de la base de données
	if err := database.InitDB(configs.AppConfig.Database.Path); err != nil {
		log.Fatalf("Erreur initialisation DB: %v", err)
	}
	defer database.Close()

	r := mux.NewRouter()

	// Configuration CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Routes statiques
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes publiques
	r.HandleFunc("/", handlers.Home)
	r.HandleFunc("/login", handlers.Login)
	r.HandleFunc("/register", handlers.Register)
	r.HandleFunc("/api/stations", handlers.GetStations)
	r.HandleFunc("/api/news", handlers.GetNews)

	// Routes d'authentification
	r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")
	r.HandleFunc("/api/auth/register", handlers.Register).Methods("POST")

	// API Routes protégées
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Routes des stations
	api.HandleFunc("/stations", handlers.GetStations).Methods("GET")
	api.HandleFunc("/stations/{id}", handlers.GetStationByID).Methods("GET")
	api.HandleFunc("/stations/nearby", handlers.GetNearbyStations).Methods("GET")
	api.HandleFunc("/stations/{id}/stats", handlers.GetStationStats).Methods("GET")

	// Routes des abonnements
	api.HandleFunc("/subscriptions", handlers.GetSubscriptionPlans).Methods("GET")
	api.HandleFunc("/subscriptions", handlers.Subscribe).Methods("POST")
	api.HandleFunc("/subscriptions/active", handlers.GetActiveSubscription).Methods("GET")

	// Routes du profil
	api.HandleFunc("/profile", handlers.GetProfile).Methods("GET")
	api.HandleFunc("/profile", handlers.UpdateProfile).Methods("PUT")
	api.HandleFunc("/profile/stats", handlers.GetProfileStats).Methods("GET")
	api.HandleFunc("/profile/ride", handlers.AddRideStats).Methods("POST")

	// Pages principales
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	}).Methods("GET")

	r.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/stations.html")
	}).Methods("GET")

	r.HandleFunc("/shop", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/shop.html")
	}).Methods("GET")

	r.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/profile.html")
	}).Methods("GET")

	// Handler 404 personnalisé
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
	})

	// Démarrage du serveur
	addr := configs.AppConfig.Server.Host + ":" + configs.AppConfig.Server.Port
	log.Printf("Serveur démarré sur http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, c.Handler(r)))
}
