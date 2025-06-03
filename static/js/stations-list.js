// État de l'application
let stations = [];
let currentView = 'list';
let searchQuery = '';
let stationMap = null;

// Chargement initial
document.addEventListener('DOMContentLoaded', () => {
    loadStations();
    setupEventListeners();

    // Ouvre automatiquement le détail si paramètre dans l'URL
    const params = new URLSearchParams(window.location.search);
    const stationId = params.get('station');
    const showDetails = params.get('showDetails');
    if (stationId && showDetails) {
        // Attendre que les stations soient chargées
        const interval = setInterval(() => {
            if (stations && stations.length > 0) {
                const station = stations.find(s => String(s.id) === String(stationId));
                if (station) {
                    showStationDetails(station);
                }
                clearInterval(interval);
            }
        }, 100);
    }
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

    // Gestion du modal de détails
    const detailsModal = document.getElementById('stationDetailsModal');
    detailsModal.addEventListener('shown.bs.modal', () => {
        if (stationMap) {
            stationMap.invalidateSize();
        }
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
                    <i class="bi bi-info-circle me-2"></i>
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
        element.querySelector('.address span').textContent = station.address;

        // Disponibilité des vélos
        const ratio = station.availableBikes / station.totalSlots;
        const progressBar = element.querySelector('.progress-bar');
        progressBar.style.width = `${ratio * 100}%`;
        progressBar.className = `progress-bar ${getStatusClass(ratio)}`;

        element.querySelector('.bikes-info span').textContent = 
            `${station.availableBikes} vélos disponibles sur ${station.totalSlots} emplacements`;

        // Statistiques
        if (station.stats) {
            element.querySelector('.popularity span').textContent = 
                `Popularité: ${Math.round(station.stats.popularityScore * 100)}%`;
            
            const peakHours = station.stats.peakHours
                .map(hour => `${hour}h`)
                .join(', ');
            element.querySelector('.peak-hours span').textContent = 
                `Heures de pointe: ${peakHours}`;
        }

        // Badge de statut
        const statusBadge = element.querySelector('.status-badge');
        statusBadge.className = `status-badge badge ${getStatusClass(ratio)}`;
        statusBadge.textContent = getStatusText(ratio);

        // Ajout du gestionnaire de clic pour les détails
        const detailsButton = element.querySelector('.view-details');
        if (detailsButton) {
            detailsButton.addEventListener('click', () => showStationDetails(station));
        }

        container.appendChild(element);
    });
}

// Fonction utilitaire pour obtenir ou générer un % batterie persistant pour une station
function getStationBatteryLevel(stationId) {
    const key = `station_battery_${stationId}`;
    let value = sessionStorage.getItem(key);
    if (value === null) {
        value = Math.floor(Math.random() * (80 - 60 + 1)) + 60;
        sessionStorage.setItem(key, value);
    }
    return value;
}

// Affichage des détails d'une station
function showStationDetails(station) {
    const modal = document.getElementById('stationDetailsModal');
    
    // Remplissage des informations de base
    modal.querySelector('.modal-title').textContent = station.name;
    modal.querySelector('.station-address').textContent = station.address;

    // Types de vélos
    const electricBikesInfo = modal.querySelector('.electric-bikes-info .text-muted');
    const batteryLevel = getStationBatteryLevel(station.id); // Utilisation du niveau persistant
    electricBikesInfo.innerHTML = `
        <div class="d-flex justify-content-between">
            <span>Nombre de vélos: ${station.electricBikes}</span>
            ${station.stats ? `<span>Utilisation: ${Math.round(station.stats.averageUsage * 100)}%</span>` : ''}
        </div>
        <div class="mt-2">
            <div class="d-flex justify-content-between align-items-center">
                <span>Niveau de batterie moyen:</span>
                <div class="progress" style="width: 100px; height: 8px;">
                    <div class="progress-bar bg-success" role="progressbar" style="width: ${batteryLevel}%"></div>
                </div>
                <span class="ms-2">${batteryLevel}%</span>
            </div>
        </div>
    `;

    const classicBikesInfo = modal.querySelector('.classic-bikes-info .text-muted');
    classicBikesInfo.innerHTML = `
        <div class="d-flex justify-content-between">
            <span>Nombre de vélos: ${station.classicBikes}</span>
        </div>
    `;

    // Statistiques - on ne les affiche que si elles existent
    const statsSection = modal.querySelector('.station-stats');
    if (station.stats && station.stats.peakHours && station.stats.peakHours.length > 0) {
        const peakHours = station.stats.peakHours
            .map(hour => `${hour}h`)
            .join(', ');
        modal.querySelector('.peak-hours span').textContent = 
            `Heures de pointe: ${peakHours}`;
        statsSection.style.display = 'block';
    } else {
        statsSection.style.display = 'none';
    }

    // État de la station
    const statusInfo = modal.querySelector('.status-info');
    const isOperational = station.status === 'operational';
    statusInfo.innerHTML = `
        <i class="bi ${isOperational ? 'bi-check-circle text-success' : 'bi-x-circle text-danger'} me-2"></i>
        <span>La station est ${isOperational ? 'en service' : 'hors service'}</span>
    `;

    // Ajout du bouton de redirection vers la carte
    const mapButtonContainer = modal.querySelector('.map-button-container');
    if (!mapButtonContainer) {
        const mapButton = document.createElement('div');
        mapButton.className = 'mt-4 text-center map-button-container';
        mapButton.innerHTML = `
            <a href="/?station=${encodeURIComponent(station.id)}&showPopup=true" 
               class="btn btn-primary">
                <i class="bi bi-map me-2"></i>
                Voir sur la carte
            </a>
        `;
        modal.querySelector('.station-details').appendChild(mapButton);
    }

    // Affichage du modal
    const modalInstance = new bootstrap.Modal(modal);
    modalInstance.show();
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