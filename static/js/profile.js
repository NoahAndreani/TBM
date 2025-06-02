// Chargement initial
document.addEventListener('DOMContentLoaded', async () => {
    try {
        // Attendre la vérification de l'authentification
        if (!authToken) {
            window.location.href = '/';
            return;
        }

        // Tenter de charger le profil
        const response = await fetch('/api/profile', {
            headers: addAuthHeader()
        });

        if (!response.ok) {
            // Si la requête échoue (token invalide ou expiré)
            window.location.href = '/';
            return;
        }

        // Si on arrive ici, l'utilisateur est authentifié
        isAuthenticated = true;
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

// Chargement du profil
async function loadProfile() {
    try {
        const response = await fetch('/api/profile', {
            headers: addAuthHeader()
        });

        if (!response.ok) {
            throw new Error('Erreur lors du chargement du profil');
        }

        const profile = await response.json();
        displayProfile(profile);

        // Chargement des statistiques
        const statsResponse = await fetch('/api/profile/stats', {
            headers: addAuthHeader()
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

// Affichage des informations du profil
function displayProfile(profile) {
    document.getElementById('username').textContent = profile.username;
    document.getElementById('email').textContent = profile.email;

    const lastConnection = new Date(profile.last_connection);
    document.getElementById('lastConnection').textContent = lastConnection.toLocaleDateString('fr-FR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });

    const memberSince = new Date(profile.created_at);
    document.getElementById('memberSince').textContent = memberSince.toLocaleDateString('fr-FR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
    });

    // Pré-remplissage du formulaire de modification
    document.getElementById('editUsername').value = profile.username;
    document.getElementById('editEmail').value = profile.email;
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
                ...addAuthHeader()
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