// Gestion de l'état de connexion
let isAuthenticated = false;
let authToken = localStorage.getItem('authToken');
let currentUser = null;

// Vérifie si l'utilisateur est déjà connecté au chargement de la page
if (authToken) {
    fetchUserProfile();
}

// Gestion du formulaire de connexion
document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            throw new Error('Identifiants invalides');
        }

        const data = await response.json();
        handleSuccessfulAuth(data);
        
        // Fermeture du modal
        const modal = bootstrap.Modal.getInstance(document.getElementById('loginModal'));
        modal.hide();
    } catch (error) {
        showError('loginForm', error.message);
    }
});

// Gestion du formulaire d'inscription
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const username = document.getElementById('registerUsername').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, email, password })
        });

        if (!response.ok) {
            throw new Error('Erreur lors de l\'inscription');
        }

        const data = await response.json();
        handleSuccessfulAuth(data);
        
        // Fermeture du modal
        const modal = bootstrap.Modal.getInstance(document.getElementById('registerModal'));
        modal.hide();
    } catch (error) {
        showError('registerForm', error.message);
    }
});

// Gestion de la déconnexion
function logout() {
    localStorage.removeItem('authToken');
    isAuthenticated = false;
    currentUser = null;
    updateUIAuth();
    window.location.href = '/';
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

// Gestion de l'authentification réussie
function handleSuccessfulAuth(data) {
    authToken = data.token;
    currentUser = data.user;
    isAuthenticated = true;
    localStorage.setItem('authToken', authToken);
    updateUIAuth();
}

// Mise à jour de l'interface utilisateur selon l'état de connexion
function updateUIAuth() {
    const authButtons = document.getElementById('authButtons');
    const userMenu = document.getElementById('userMenu');

    if (isAuthenticated && currentUser) {
        authButtons.classList.add('d-none');
        userMenu.classList.remove('d-none');
    } else {
        authButtons.classList.remove('d-none');
        userMenu.classList.add('d-none');
    }
}

// Affichage des erreurs dans les formulaires
function showError(formId, message) {
    const form = document.getElementById(formId);
    let errorDiv = form.querySelector('.alert');
    
    if (!errorDiv) {
        errorDiv = document.createElement('div');
        errorDiv.className = 'alert alert-danger mt-3';
        form.appendChild(errorDiv);
    }
    
    errorDiv.textContent = message;
}

// Ajout du token d'authentification à toutes les requêtes API
function addAuthHeader(headers = {}) {
    if (authToken) {
        return {
            ...headers,
            'Authorization': `Bearer ${authToken}`
        };
    }
    return headers;
} 