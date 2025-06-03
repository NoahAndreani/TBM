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

// Fonction pour charger la liste des utilisateurs
async function loadUsers() {
    try {
        const response = await fetch('/admin/users', {
            headers: {
                'Authorization': `Bearer ${window.authToken}`
            }
        });

        if (!response.ok) {
            throw new Error('Erreur lors du chargement des utilisateurs');
        }

        const users = await response.json();
        displayUsers(users);
    } catch (error) {
        console.error('Erreur:', error);
        alert('Erreur lors du chargement des utilisateurs');
    }
}

// Fonction pour afficher les utilisateurs dans le tableau
function displayUsers(users) {
    const tbody = document.getElementById('usersList');
    tbody.innerHTML = '';

    users.forEach(user => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${user.id}</td>
            <td>${user.username}</td>
            <td>${user.email}</td>
            <td>${user.role}</td>
            <td>${user.level}</td>
            <td>${new Date(user.created_at).toLocaleString()}</td>
            <td>${new Date(user.last_connection).toLocaleString()}</td>
            <td>
                ${user.role !== 'admin' ? `
                    <button class="btn btn-danger btn-sm" onclick="deleteUser(${user.id})">
                        <i class="bi bi-trash"></i> Supprimer
                    </button>
                ` : ''}
            </td>
        `;
        tbody.appendChild(tr);
    });
}

// Fonction pour supprimer un utilisateur
async function deleteUser(userId) {
    if (!confirm('Êtes-vous sûr de vouloir supprimer cet utilisateur ?')) {
        return;
    }

    try {
        const response = await fetch(`/admin/users/${userId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${window.authToken}`
            }
        });

        if (!response.ok) {
            throw new Error('Erreur lors de la suppression de l\'utilisateur');
        }

        // Recharger la liste des utilisateurs
        loadUsers();
    } catch (error) {
        console.error('Erreur:', error);
        alert('Erreur lors de la suppression de l\'utilisateur');
    }
}

// Initialisation au chargement de la page
document.addEventListener('DOMContentLoaded', async () => {
    try {
        // Attendre que l'authentification soit prête
        await waitForAuth();

        // Vérifier si l'utilisateur est connecté
        if (!window.authToken) {
            window.location.href = '/';
            return;
        }

        // Charger les utilisateurs
        loadUsers();

        // Afficher le menu utilisateur
        document.getElementById('userMenu').classList.remove('d-none');
    } catch (error) {
        console.error('Erreur:', error);
        window.location.href = '/';
    }
}); 