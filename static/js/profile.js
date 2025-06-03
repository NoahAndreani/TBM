// Fonctions d'authentification
function getAuthToken() {
    return localStorage.getItem('token');
}

function addAuthHeader() {
    return {
        'Authorization': `Bearer ${getAuthToken()}`
    };
}

// Fonction pour vérifier si l'authentification est prête
function waitForAuth() {
    return new Promise((resolve) => {
        // Si authToken est déjà défini, on résout immédiatement
        if (typeof window.authToken !== 'undefined') {
            resolve();
            return;
        }

        // Sinon, on vérifie toutes les 100ms jusqu'à ce que ce soit prêt
        const checkAuth = setInterval(() => {
            if (typeof window.authToken !== 'undefined') {
                clearInterval(checkAuth);
                resolve();
            }
        }, 100);

        // On arrête de vérifier après 5 secondes
        setTimeout(() => {
            clearInterval(checkAuth);
            resolve();
        }, 5000);
    });
}

// Chargement initial
document.addEventListener('DOMContentLoaded', async () => {
    try {
        // Attendre que l'authentification soit prête
        await waitForAuth();

        // Récupération du token d'authentification
        if (!window.authToken) {
            console.log('Aucun token trouvé, redirection vers la page d\'accueil');
            window.location.href = '/';
            return;
        }

        // Tenter de charger le profil
        const response = await fetch('/api/profile', {
            headers: window.addAuthHeader()
        });

        if (!response.ok) {
            // Si la requête échoue (token invalide ou expiré)
            console.log('Token invalide ou expiré, redirection vers la page d\'accueil');
            localStorage.removeItem('token');
            window.authToken = null;
            window.location.href = '/';
            return;
        }

        // Si on arrive ici, l'utilisateur est authentifié
        loadProfile();
        setupEventListeners();
    } catch (error) {
        console.error('Erreur de vérification d\'authentification:', error);
        window.location.href = '/';
    }
});

// Configuration des écouteurs d'événements
function setupEventListeners() {
    document.getElementById('editProfileForm').addEventListener('submit', (e) => {
        e.preventDefault();
        updateProfile();
    });
}

// Fonction pour charger le profil de l'utilisateur
async function loadProfile() {
    try {
        const response = await fetch('/api/profile', {
            headers: {
                'Authorization': `Bearer ${window.authToken}`
            }
        });

        if (!response.ok) {
            throw new Error('Erreur lors du chargement du profil');
        }

        const user = await response.json();
        displayProfile(user);

        // Afficher le bouton d'administration si l'utilisateur est admin
        if (user.role === 'admin') {
            document.getElementById('adminSection').classList.remove('d-none');
        }

        // Chargement des statistiques
        const statsResponse = await fetch('/api/profile/stats', {
            headers: window.addAuthHeader()
        });

        if (statsResponse.ok) {
            const stats = await statsResponse.json();
            displayStats(stats);
        }
    } catch (error) {
        console.error('Erreur:', error);
        showErrorMessage('Impossible de charger les informations du profil');
    }
}

// Fonction pour afficher les informations du profil
function displayProfile(user) {
    document.getElementById('username').textContent = user.username;
    document.getElementById('email').textContent = user.email;
    document.getElementById('currentLevel').textContent = user.level;
    document.getElementById('experiencePoints').textContent = `${user.experience} XP`;
    
    // Calculer la progression vers le niveau suivant
    const expForNextLevel = user.level * 1000;
    const progress = (user.experience / expForNextLevel) * 100;
    document.getElementById('experienceBar').style.width = `${progress}%`;
    document.getElementById('experienceBar').setAttribute('aria-valuenow', progress);
    
    // Afficher les statistiques
    document.getElementById('totalDistance').textContent = `${user.total_distance.toFixed(2)} km`;
    document.getElementById('totalTime').textContent = formatTime(user.total_ride_time);
    document.getElementById('consecutiveDays').textContent = user.consecutive_days;

    // Afficher les dates
    document.getElementById('lastConnection').textContent = new Date(user.last_connection).toLocaleString();
    document.getElementById('memberSince').textContent = new Date(user.created_at).toLocaleString();

    // Pré-remplissage du formulaire de modification
    document.getElementById('editUsername').value = user.username;
    document.getElementById('editEmail').value = user.email;

    // Afficher l'expérience restante pour le niveau suivant
    document.getElementById('experienceToNext').textContent = 
        `${expForNextLevel - user.experience} XP pour le niveau suivant`;

    // Afficher les récompenses du niveau
    displayLevelRewards(user.level);
}

// Fonction pour formater le temps en heures et minutes
function formatTime(minutes) {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins}min`;
}

// Affichage des statistiques
function displayStats(stats) {
    // Niveau et expérience
    document.getElementById('currentLevel').textContent = stats.level;
    document.getElementById('experiencePoints').textContent = `${stats.experience} XP`;

    const experienceProgress = (stats.experience / stats.experience_to_next) * 100;
    const experienceBar = document.getElementById('experienceBar');
    experienceBar.style.width = `${experienceProgress}%`;
    experienceBar.setAttribute('aria-valuenow', experienceProgress);

    document.getElementById('experienceToNext').textContent = 
        `${stats.experience_to_next - stats.experience} XP pour le niveau suivant`;

    // Statistiques générales
    document.getElementById('totalDistance').textContent = 
        `${stats.total_distance.toFixed(1)} km`;
    
    const hours = Math.floor(stats.total_ride_time / 60);
    const minutes = stats.total_ride_time % 60;
    document.getElementById('totalTime').textContent = 
        `${hours}h${minutes > 0 ? ` ${minutes}min` : ''}`;

    document.getElementById('consecutiveDays').textContent = 
        stats.consecutive_days;

    // Récompenses du niveau
    displayLevelRewards(stats.level);
}

// Affichage des récompenses du niveau
function displayLevelRewards(level) {
    const rewards = getLevelRewards(level);
    const rewardsList = document.getElementById('levelRewards');
    rewardsList.innerHTML = '';

    rewards.forEach(reward => {
        const li = document.createElement('li');
        li.className = 'mb-2';
        li.innerHTML = `
            <i class="bi bi-check-circle-fill text-success me-2"></i>
            ${reward}
        `;
        rewardsList.appendChild(li);
    });
}

// Obtention des récompenses pour un niveau
function getLevelRewards(level) {
    const rewards = [
        'Accès aux vélos standards',
        'Bonus XP de base'
    ];

    if (level >= 5) {
        rewards.push('Bonus XP x1.5');
    }
    if (level >= 10) {
        rewards.push('Accès aux vélos électriques');
        rewards.push('Bonus XP x2');
    }
    if (level >= 15) {
        rewards.push('Temps de location gratuit +15min');
    }
    if (level >= 20) {
        rewards.push('Accès prioritaire aux stations');
        rewards.push('Bonus XP x3');
    }

    return rewards;
}

// Mise à jour du profil
async function updateProfile() {
    const username = document.getElementById('editUsername').value;
    const email = document.getElementById('editEmail').value;
    const password = document.getElementById('editPassword').value;

    try {
        const response = await fetch('/api/profile', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                ...window.addAuthHeader()
            },
            body: JSON.stringify({
                username,
                email,
                ...(password && { password })
            })
        });

        if (!response.ok) {
            throw new Error('Erreur lors de la mise à jour du profil');
        }

        const updatedProfile = await response.json();
        displayProfile(updatedProfile);

        // Fermeture du modal
        const modal = bootstrap.Modal.getInstance(document.getElementById('editProfileModal'));
        modal.hide();

        // Réinitialisation du champ mot de passe
        document.getElementById('editPassword').value = '';

        showSuccessMessage('Profil mis à jour avec succès');
    } catch (error) {
        console.error('Erreur:', error);
        showErrorMessage('Impossible de mettre à jour le profil');
    }
}

// Affichage des messages
function showSuccessMessage(message) {
    const container = document.createElement('div');
    container.className = 'alert alert-success alert-dismissible fade show';
    container.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    document.querySelector('.container').insertBefore(container, document.querySelector('.container').firstChild);
}

function showErrorMessage(message) {
    const container = document.createElement('div');
    container.className = 'alert alert-danger alert-dismissible fade show';
    container.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    document.querySelector('.container').insertBefore(container, document.querySelector('.container').firstChild);
} 