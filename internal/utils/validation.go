package utils

import (
	"regexp"
	"unicode"
)

// ValidateUsername vérifie si le pseudo est valide (lettres, chiffres et underscore uniquement)
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

// ValidateEmail vérifie si l'email respecte le format spécifié
func ValidateEmail(email string) bool {
	// Format: texte@texte_entre_3et6_caracteres.caractere_2a5
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9][a-zA-Z0-9-]{1,4}[a-zA-Z0-9]\.[a-zA-Z]{2,5}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ValidatePassword vérifie si le mot de passe respecte les critères
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}

// GetValidationErrors retourne les messages d'erreur pour chaque champ
func GetValidationErrors(username, email, password string) map[string]string {
	errors := make(map[string]string)

	if !ValidateUsername(username) {
		errors["username"] = "Le pseudo doit contenir entre 3 et 20 caractères et uniquement des lettres, chiffres et underscore"
	}

	if !ValidateEmail(email) {
		errors["email"] = "L'adresse email n'est pas valide"
	}

	if !ValidatePassword(password) {
		errors["password"] = "Le mot de passe doit contenir au moins 8 caractères, une majuscule, une minuscule et un chiffre"
	}

	return errors
}
