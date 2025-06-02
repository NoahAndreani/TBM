// Initialisation de la carte centrée sur Bordeaux
const map = L.map('map').setView([44.837789, -0.57918], 13);

// Ajout de la couche OpenStreetMap
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors'
}).addTo(map);

// Stockage des marqueurs des stations
let stationMarkers = [];

// Icônes personnalisées pour les stations
const stationIcons = {
    available: L.icon({
        iconUrl: '/static/images/marker-green.svg',
        iconSize: [25, 41],
        iconAnchor: [12, 41],
        popupAnchor: [1, -34]
    }),
    few: L.icon({
        iconUrl: '/static/images/marker-orange.svg',
        iconSize: [25, 41],
        iconAnchor: [12, 41],
        popupAnchor: [1, -34]
    }),
    empty: L.icon({
        iconUrl: '/static/images/marker-red.svg',
        iconSize: [25, 41],
        iconAnchor: [12, 41],
        popupAnchor: [1, -34]
    })
};

// Fonction pour charger les stations
async function loadStations() {
    try {
        const response = await fetch('/api/stations');
        const stations = await response.json();
        
        // Suppression des marqueurs existants
        stationMarkers.forEach(marker => marker.remove());
        stationMarkers = [];

        // Ajout des nouveaux marqueurs
        stations.forEach(station => {
            const icon = getStationIcon(station);
            const marker = L.marker([station.latitude, station.longitude], { icon })
                .bindPopup(createStationPopup(station))
                .addTo(map);
            stationMarkers.push(marker);
        });

        // Mise à jour de la liste des stations proches
        updateNearbyStations(stations);
    } catch (error) {
        console.error('Erreur lors du chargement des stations:', error);
    }
}

// Fonction pour déterminer l'icône en fonction de la disponibilité
function getStationIcon(station) {
    const ratio = station.availableBikes / station.totalSlots;
    if (ratio > 0.3) return stationIcons.available;
    if (ratio > 0) return stationIcons.few;
    return stationIcons.empty;
}

// Fonction pour créer le contenu du popup d'une station
function createStationPopup(station) {
    return `
        <div class="station-popup">
            <h5>${station.name}</h5>
            <p>${station.address}</p>
            <div class="station-status">
                <strong>Vélos disponibles:</strong> ${station.availableBikes} / ${station.totalSlots}
            </div>
            <div class="station-status">
                <strong>État:</strong> ${station.status === 'operational' ? 'En service' : 'Hors service'}
            </div>
        </div>
    `;
}

// Fonction pour mettre à jour la liste des stations proches
function updateNearbyStations(stations) {
    const nearbyContainer = document.getElementById('nearbyStations');
    nearbyContainer.innerHTML = '';

    // Tri des stations par distance si la géolocalisation est disponible
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(position => {
            const userLat = position.coords.latitude;
            const userLng = position.coords.longitude;

            stations.forEach(station => {
                station.distance = getDistance(userLat, userLng, station.latitude, station.longitude);
            });

            stations.sort((a, b) => a.distance - b.distance);
            displayNearbyStations(stations.slice(0, 5));
        });
    } else {
        displayNearbyStations(stations.slice(0, 5));
    }
}

// Fonction pour afficher les stations proches
function displayNearbyStations(stations) {
    const nearbyContainer = document.getElementById('nearbyStations');
    stations.forEach(station => {
        const ratio = station.availableBikes / station.totalSlots;
        const statusClass = ratio > 0.3 ? 'success' : (ratio > 0 ? 'warning' : 'danger');

        const stationElement = document.createElement('a');
        stationElement.href = '#';
        stationElement.className = `list-group-item list-group-item-${statusClass}`;
        stationElement.innerHTML = `
            <div class="d-flex w-100 justify-content-between">
                <h6 class="mb-1">${station.name}</h6>
                <small>${station.distance ? Math.round(station.distance * 1000) + 'm' : ''}</small>
            </div>
            <p class="mb-1">${station.availableBikes} vélos disponibles</p>
        `;

        stationElement.addEventListener('click', () => {
            map.setView([station.latitude, station.longitude], 16);
            const marker = stationMarkers.find(m => 
                m.getLatLng().lat === station.latitude && 
                m.getLatLng().lng === station.longitude
            );
            if (marker) marker.openPopup();
        });

        nearbyContainer.appendChild(stationElement);
    });
}

// Fonction pour calculer la distance entre deux points (formule de Haversine)
function getDistance(lat1, lon1, lat2, lon2) {
    const R = 6371; // Rayon de la Terre en km
    const dLat = toRad(lat2 - lat1);
    const dLon = toRad(lon2 - lon1);
    const a = 
        Math.sin(dLat/2) * Math.sin(dLat/2) +
        Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * 
        Math.sin(dLon/2) * Math.sin(dLon/2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
    return R * c;
}

function toRad(value) {
    return value * Math.PI / 180;
}

// Chargement initial des stations
loadStations();

// Rafraîchissement des stations toutes les 30 secondes
setInterval(loadStations, 30000); 