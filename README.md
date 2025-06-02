# TBC Vclub - Location de Vélos à Bordeaux

## Description
TBC Vclub est une application web de location de vélos pour la ville de Bordeaux. Elle permet aux utilisateurs de localiser les stations de vélos, consulter leur disponibilité en temps réel, et gérer leurs abonnements.

## Fonctionnalités
- Carte interactive des stations de vélos
- Tableau des stations proches
- Actualités de Bordeaux
- Système d'abonnement et de forfaits
- Système de niveau et d'expérience
- Profil utilisateur personnalisé
- Statistiques des stations

## Prérequis
- Go 1.16 ou supérieur
- MySQL 5.7 ou supérieur
- Node.js 14 ou supérieur (pour les assets front-end)

## Installation

1. Cloner le repository
```bash
git clone https://github.com/votre-username/tbcvclub.git
cd tbcvclub
```

2. Installer les dépendances
```bash
go mod tidy
```

3. Configurer la base de données
- Créer une base de données MySQL
- Copier `configs/config.json.example` vers `configs/config.json`
- Modifier les paramètres de connexion dans `configs/config.json`

4. Lancer l'application
```bash
go run cmd/server/main.go
```

L'application sera accessible à l'adresse : http://localhost:8080

## Structure du Projet
```
tbcvclub/
├── cmd/
│   └── server/
│       └── main.go
├── configs/
│   ├── config.go
│   └── config.json
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── middleware/
│   └── utils/
├── static/
├── templates/
└── README.md
```

## API Endpoints

### Authentification
- POST /register - Inscription
- POST /login - Connexion
- POST /logout - Déconnexion

### Stations
- GET /api/stations - Liste des stations
- GET /api/stations/{id} - Détails d'une station
- GET /api/stations/nearby - Stations proches

### Abonnements
- GET /api/subscriptions - Liste des forfaits
- POST /api/subscriptions - Souscrire à un forfait
- GET /api/subscriptions/active - Abonnement actif

### Profil
- GET /api/profile - Informations du profil
- PUT /api/profile - Mise à jour du profil
- GET /api/profile/stats - Statistiques utilisateur

## Système de Niveau et d'Expérience
- Gain d'XP basé sur :
  - Connexions quotidiennes
  - Distance parcourue
  - Temps de location
- Niveaux débloquant des avantages
- Statistiques détaillées

## Contribution
Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou une pull request.

## Licence
Ce projet est sous licence MIT. 