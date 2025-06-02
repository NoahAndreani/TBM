// État de l'application
let stations = [];
let currentView = 'list';
let searchQuery = '';

// Chargement initial
document.addEventListener('DOMContentLoaded', () => {
    loadStations();
    setupEventListeners();
});

// Configuration des écouteurs d'événements
function setupEventListeners() {
    // Gestion de la recherche
    const searchInput = document.getElementById('searchStation');
    const searchButton = document.getElementById('searchButton');

    searchInput.addEventListener('input', (e) => {
        searchQuery = e.target.value.toLowerCase();
        filterAndDisplayStations();
    });

    searchButton.addEventListener('click', () => {
        filterAndDisplayStations();
    });

    // Gestion des vues (liste/grille)
    const viewButtons = document.querySelectorAll('[data-view]');
    viewButtons.forEach(button => {
        button.addEventListener('click', (e) => {
            const newView = e.target.closest('button').dataset.view;
            changeView(newView);
        });
    });
}

// Chargement des stations
async function loadStations() {
    try {
        const response = await fetch('/api/stations');
        stations = await response.json();
        
        // Chargement des statistiques pour chaque station
        await Promise.all(stations.map(async (station) => {
            try {
                const statsResponse = await fetch(`/api/stations/${station.id}/stats`);
                station.stats = await statsResponse.json();
            } catch (error) {
                console.error(`Erreur lors du chargement des stats pour la station ${station.id}:`, error);
                station.stats = null;
            }
        }));

        filterAndDisplayStations();
    } catch (error) {
        console.error('Erreur lors du chargement des stations:', error);
        displayError();
    }
}

// Filtrage et affichage des stations
function filterAndDisplayStations() {
    const filteredStations = stations.filter(station => 
        station.name.toLowerCase().includes(searchQuery) ||
        station.address.toLowerCase().includes(searchQuery)
    );

    const container = document.getElementById('stationsList');
    container.innerHTML = '';

    if (filteredStations.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    Aucune station ne correspond à votre recherche.
                </div>
            </div>
        `;
        return;
    }

    const template = document.getElementById(
        currentView === 'list' ? 'stationListTemplate' : 'stationGridTemplate'
    );

    filteredStations.forEach(station => {
        const element = template.content.cloneNode(true);
        
        // Remplissage des informations de base
        element.querySelector('.card-title').textContent = station.name;
        element.querySelector('.address').textContent = station.address;

        // Disponibilité des vélos
        const ratio = station.availableBikes / station.totalSlots;
        const progressBar = element.querySelector('.progress-bar');
        progressBar.style.width = `${ratio * 100}%`;
        progressBar.className = `progress-bar ${getStatusClass(ratio)}`;

        element.querySelector('.bikes-info').textContent = 
            `${station.availableBikes} vélos disponibles sur ${station.totalSlots} emplacements`;

        // Statistiques
        if (station.stats) {
            element.querySelector('.popularity').textContent = 
                `Popularité: ${Math.round(station.stats.popularityScore * 100)}%`;
            
            const peakHours = station.stats.peakHours
                .map(hour => `${hour}h`)
                .join(', ');
            element.querySelector('.peak-hours').textContent = 
                `Heures de pointe: ${peakHours}`;
        }

        // Badge de statut
        const statusBadge = element.querySelector('.status-badge');
        statusBadge.className = `status-badge badge ${getStatusClass(ratio)}`;
        statusBadge.textContent = getStatusText(ratio);

        container.appendChild(element);
    });
}

// Changement de vue (liste/grille)
function changeView(newView) {
    const buttons = document.querySelectorAll('[data-view]');
    buttons.forEach(button => {
        button.classList.toggle('active', button.dataset.view === newView);
    });

    currentView = newView;
    filterAndDisplayStations();
}

// Utilitaires
function getStatusClass(ratio) {
    if (ratio > 0.3) return 'bg-success';
    if (ratio > 0) return 'bg-warning';
    return 'bg-danger';
}

function getStatusText(ratio) {
    if (ratio > 0.3) return 'Bien disponible';
    if (ratio > 0) return 'Peu disponible';
    return 'Non disponible';
}

// Affichage des erreurs
function displayError() {
    const container = document.getElementById('stationsList');
    container.innerHTML = `
        <div class="col-12">
            <div class="alert alert-danger">
                <i class="bi bi-exclamation-triangle-fill"></i>
                Une erreur est survenue lors du chargement des stations.
                <button type="button" class="btn btn-danger btn-sm ms-3" onclick="loadStations()">
                    Réessayer
                </button>
            </div>
        </div>
    `;
}

// Rafraîchissement automatique toutes les 30 secondes
setInterval(loadStations, 30000); 