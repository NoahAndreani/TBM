package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Générer le hash du mot de passe "root"
	password := []byte("root")
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Ouvrir la base de données
	db, err := sql.Open("sqlite3", "data/tbcvclub.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Mettre à jour le mot de passe de l'admin
	_, err = db.Exec("UPDATE users SET hashed_password = ? WHERE username = 'admin'", string(hashedPassword))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Le mot de passe de l'admin a été mis à jour avec succès !")
}
