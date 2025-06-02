// Gestion de l'état de connexion
let isAuthenticated = false;
let authToken = localStorage.getItem('authToken');
let currentUser = null;

// Initialisation des modals
let loginModal;
let registerModal;

// Vérifie si l'utilisateur est déjà connecté au chargement de la page
if (authToken) {
    fetchUserProfile();
}

// Validation côté client
const validators = {
    username: (value) => {
        const regex = /^[a-zA-Z0-9_]+$/;
        if (value.length < 3 || value.length > 20) {
            return "Le pseudo doit contenir entre 3 et 20 caractères";
        }
        if (!regex.test(value)) {
            return "Le pseudo ne peut contenir que des lettres, chiffres et underscore";
        }
        return null;
    },
    email: (value) => {
        const regex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9][a-zA-Z0-9-]{1,4}[a-zA-Z0-9]\.[a-zA-Z]{2,5}$/;
        if (!regex.test(value)) {
            return "L'adresse email n'est pas valide";
        }
        return null;
    },
    password: (value) => {
        if (value.length < 8) {
            return "Le mot de passe doit contenir au moins 8 caractères";
        }
        if (!/[A-Z]/.test(value)) {
            return "Le mot de passe doit contenir au moins une majuscule";
        }
        if (!/[a-z]/.test(value)) {
            return "Le mot de passe doit contenir au moins une minuscule";
        }
        if (!/[0-9]/.test(value)) {
            return "Le mot de passe doit contenir au moins un chiffre";
        }
        return null;
    }
};

// Initialisation des formulaires d'authentification
document.addEventListener('DOMContentLoaded', () => {
    // Initialisation des modals
    loginModal = new bootstrap.Modal(document.getElementById('loginModal'));
    registerModal = new bootstrap.Modal(document.getElementById('registerModal'));

    // Gestion du formulaire de connexion
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }

    // Gestion du formulaire d'inscription
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }

    // Mise à jour initiale de l'interface
    updateUIAuth();
});

// Fonction de connexion
async function handleLogin(event) {
    event.preventDefault();
    
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;
    const errorDiv = document.getElementById('loginError');
    
    try {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();

        if (response.ok) {
            // Stockage du token
            localStorage.setItem('token', data.token);
            // Fermeture du modal
            loginModal.hide();
            // Mise à jour de l'interface
            updateAuthUI(true);
            // Redirection vers la page d'accueil
            window.location.href = '/';
        } else {
            errorDiv.textContent = data.error || 'Erreur de connexion';
            errorDiv.classList.remove('d-none');
        }
    } catch (error) {
        errorDiv.textContent = 'Erreur de connexion au serveur';
        errorDiv.classList.remove('d-none');
    }
}

// Fonction d'inscription
async function handleRegister(event) {
    event.preventDefault();
    
    const username = document.getElementById('registerUsername').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;
    const passwordConfirm = document.getElementById('registerPasswordConfirm').value;
    const errorDiv = document.getElementById('registerError');

    // Validation côté client
    if (password !== passwordConfirm) {
        errorDiv.textContent = 'Les mots de passe ne correspondent pas';
        errorDiv.classList.remove('d-none');
        return;
    }

    if (!/^[a-zA-Z0-9]{3,20}$/.test(username)) {
        errorDiv.textContent = 'Le nom d\'utilisateur doit contenir entre 3 et 20 caractères alphanumériques';
        errorDiv.classList.remove('d-none');
        return;
    }

    if (!/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/.test(password)) {
        errorDiv.textContent = 'Le mot de passe doit contenir au moins 8 caractères, une majuscule, une minuscule et un chiffre';
        errorDiv.classList.remove('d-none');
        return;
    }

    try {
        const response = await fetch('/api/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, password })
        });

        const data = await response.json();

        if (response.ok) {
            // Fermeture du modal d'inscription
            registerModal.hide();
            // Ouverture du modal de connexion
            loginModal.show();
        } else {
            errorDiv.textContent = data.error || 'Erreur lors de l\'inscription';
            errorDiv.classList.remove('d-none');
        }
    } catch (error) {
        errorDiv.textContent = 'Erreur de connexion au serveur';
        errorDiv.classList.remove('d-none');
    }
}

// Fonction de déconnexion
function logout() {
    localStorage.removeItem('token');
    updateAuthUI(false);
    window.location.href = '/';
}

// Mise à jour de l'interface utilisateur en fonction de l'état de connexion
function updateAuthUI(isLoggedIn) {
    const authButtons = document.getElementById('authButtons');
    const userMenu = document.getElementById('userMenu');
    
    if (isLoggedIn) {
        authButtons.classList.add('d-none');
        userMenu.classList.remove('d-none');
    } else {
        authButtons.classList.remove('d-none');
        userMenu.classList.add('d-none');
    }
}

// Vérification de l'état de connexion au chargement de la page
document.addEventListener('DOMContentLoaded', function() {
    const token = localStorage.getItem('token');
    updateAuthUI(!!token);
});

// Gestion des erreurs
function showFieldError(fieldId, message) {
    const field = document.getElementById(fieldId);
    if (!field) return;

    field.classList.add('is-invalid');
    
    // Création ou mise à jour du message d'erreur
    let errorDiv = field.nextElementSibling;
    if (!errorDiv || !errorDiv.classList.contains('invalid-feedback')) {
        errorDiv = document.createElement('div');
        errorDiv.className = 'invalid-feedback';
        field.parentNode.insertBefore(errorDiv, field.nextSibling);
    }
    errorDiv.textContent = message;
}

function showFormError(form, message) {
    const existingAlert = form.querySelector('.alert');
    if (existingAlert) {
        existingAlert.remove();
    }

    const alert = document.createElement('div');
    alert.className = 'alert alert-danger mt-3';
    alert.textContent = message;
    form.appendChild(alert);
}

function clearErrors(form) {
    // Suppression des messages d'erreur
    form.querySelectorAll('.alert').forEach(alert => alert.remove());
    
    // Réinitialisation des champs invalides
    form.querySelectorAll('.is-invalid').forEach(field => {
        field.classList.remove('is-invalid');
        const errorDiv = field.nextElementSibling;
        if (errorDiv && errorDiv.classList.contains('invalid-feedback')) {
            errorDiv.remove();
        }
    });
}

// Récupération du profil utilisateur
async function fetchUserProfile() {
    try {
        const response = await fetch('/api/profile', {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        if (!response.ok) {
            throw new Error('Session expirée');
        }

        currentUser = await response.json();
        isAuthenticated = true;
        updateUIAuth();
    } catch (error) {
        console.error('Erreur lors de la récupération du profil:', error);
        logout();
    }
}

// Mise à jour de l'interface utilisateur selon l'état de connexion
function updateUIAuth() {
    const authButtons = document.getElementById('authButtons');
    const userMenu = document.getElementById('userMenu');
    const userMenuUsername = document.getElementById('userMenuUsername');

    if (authButtons && userMenu) {
        if (isAuthenticated && currentUser) {
            authButtons.classList.add('d-none');
            userMenu.classList.remove('d-none');
            if (userMenuUsername) {
                userMenuUsername.textContent = currentUser.username;
            }
        } else {
            authButtons.classList.remove('d-none');
            userMenu.classList.add('d-none');
        }
    }
} 