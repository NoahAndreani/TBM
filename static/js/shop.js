// État de l'application
let currentSubscription = null;
let selectedPlan = null;

// Chargement initial
document.addEventListener('DOMContentLoaded', () => {
    loadSubscriptionPlans();
    if (isAuthenticated) {
        loadActiveSubscription();
    }
    updateAuthUI();
    setupEventListeners();
});

// Configuration des écouteurs d'événements
function setupEventListeners() {
    // Gestion de la confirmation d'abonnement
    document.getElementById('confirmSubscription').addEventListener('click', () => {
        subscribeToSelectedPlan();
    });
}

// Chargement des forfaits disponibles
async function loadSubscriptionPlans() {
    try {
        const response = await fetch('/api/subscriptions');
        const plans = await response.json();
        displaySubscriptionPlans(plans);
    } catch (error) {
        console.error('Erreur lors du chargement des forfaits:', error);
        displayError();
    }
}

// Chargement de l'abonnement actif
async function loadActiveSubscription() {
    try {
        const response = await fetch('/api/subscriptions/active', {
            headers: addAuthHeader()
        });
        
        if (response.ok) {
            currentSubscription = await response.json();
            displayActiveSubscription();
        }
    } catch (error) {
        console.error('Erreur lors du chargement de l\'abonnement actif:', error);
    }
}

// Affichage des forfaits
function displaySubscriptionPlans(plans) {
    const container = document.getElementById('subscriptionPlans');
    const template = document.getElementById('planTemplate');
    container.innerHTML = '';

    plans.forEach(plan => {
        const element = template.content.cloneNode(true);
        
        element.querySelector('.card-title').textContent = plan.name;
        element.querySelector('.price .h2').textContent = plan.price;
        element.querySelector('.description').textContent = plan.description;

        const benefitsList = element.querySelector('.benefits');
        plan.benefits.forEach(benefit => {
            const li = document.createElement('li');
            li.innerHTML = `<i class="bi bi-check-circle-fill text-success me-2"></i>${benefit}`;
            benefitsList.appendChild(li);
        });

        const subscribeBtn = element.querySelector('.subscribe-btn');
        subscribeBtn.disabled = !isAuthenticated || (currentSubscription && currentSubscription.isActive);
        subscribeBtn.addEventListener('click', () => {
            selectedPlan = plan;
            showConfirmationModal(plan);
        });

        container.appendChild(element);
    });
}

// Affichage de l'abonnement actif
function displayActiveSubscription() {
    const container = document.getElementById('activeSubscription');
    if (!currentSubscription || !currentSubscription.isActive) {
        container.classList.add('d-none');
        return;
    }

    container.classList.remove('d-none');
    document.getElementById('subType').textContent = currentSubscription.type;
    
    const endDate = new Date(currentSubscription.endDate);
    document.getElementById('subExpiration').textContent = 
        endDate.toLocaleDateString('fr-FR', { 
            day: '2-digit',
            month: '2-digit',
            year: 'numeric'
        });

    // Calcul du temps restant
    const now = new Date();
    const startDate = new Date(currentSubscription.startDate);
    const totalDuration = endDate - startDate;
    const remainingDuration = endDate - now;
    const progress = ((totalDuration - remainingDuration) / totalDuration) * 100;

    const progressBar = document.getElementById('subProgress');
    progressBar.style.width = `${progress}%`;
    progressBar.setAttribute('aria-valuenow', progress);

    // Affichage du temps restant
    const daysLeft = Math.ceil(remainingDuration / (1000 * 60 * 60 * 24));
    document.getElementById('subTimeLeft').textContent = 
        `${daysLeft} jour${daysLeft > 1 ? 's' : ''} restant${daysLeft > 1 ? 's' : ''}`;
}

// Affichage du modal de confirmation
function showConfirmationModal(plan) {
    document.getElementById('confirmPlanName').textContent = plan.name;
    document.getElementById('confirmPlanPrice').textContent = plan.price;
    
    let duration;
    switch (plan.type) {
        case 'daily':
            duration = '24 heures';
            break;
        case 'weekly':
            duration = '7 jours';
            break;
        case 'monthly':
            duration = '30 jours';
            break;
        case 'yearly':
            duration = '1 an';
            break;
    }
    document.getElementById('confirmPlanDuration').textContent = duration;

    const modal = new bootstrap.Modal(document.getElementById('confirmSubscriptionModal'));
    modal.show();
}

// Souscription à un forfait
async function subscribeToSelectedPlan() {
    if (!selectedPlan) return;

    try {
        const response = await fetch('/api/subscriptions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...addAuthHeader()
            },
            body: JSON.stringify({
                plan_id: selectedPlan.id
            })
        });

        if (!response.ok) {
            throw new Error('Erreur lors de la souscription');
        }

        currentSubscription = await response.json();
        
        // Fermeture du modal
        const modal = bootstrap.Modal.getInstance(document.getElementById('confirmSubscriptionModal'));
        modal.hide();

        // Mise à jour de l'interface
        displayActiveSubscription();
        loadSubscriptionPlans(); // Pour mettre à jour les boutons

        // Message de succès
        showSuccessMessage('Votre abonnement a été activé avec succès !');
    } catch (error) {
        console.error('Erreur lors de la souscription:', error);
        showErrorMessage('Une erreur est survenue lors de la souscription. Veuillez réessayer.');
    }
}

// Mise à jour de l'interface selon l'état de connexion
function updateAuthUI() {
    const loginPrompt = document.getElementById('loginPrompt');
    loginPrompt.classList.toggle('d-none', isAuthenticated);

    if (isAuthenticated) {
        loadActiveSubscription();
    } else {
        document.getElementById('activeSubscription').classList.add('d-none');
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

// Affichage des erreurs
function displayError() {
    const container = document.getElementById('subscriptionPlans');
    container.innerHTML = `
        <div class="col-12">
            <div class="alert alert-danger">
                <i class="bi bi-exclamation-triangle-fill"></i>
                Une erreur est survenue lors du chargement des forfaits.
                <button type="button" class="btn btn-danger btn-sm ms-3" onclick="loadSubscriptionPlans()">
                    Réessayer
                </button>
            </div>
        </div>
    `;
} 