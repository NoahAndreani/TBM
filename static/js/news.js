// Fonction pour charger les actualités
async function loadNews() {
    try {
        const response = await fetch('/api/news');
        const news = await response.json();
        displayNews(news);
    } catch (error) {
        console.error('Erreur lors du chargement des actualités:', error);
        displayError();
    }
}

// Fonction pour afficher les actualités
function displayNews(news) {
    const newsContainer = document.getElementById('news');
    newsContainer.innerHTML = '';

    if (!news || news.length === 0) {
        newsContainer.innerHTML = '<p class="text-muted">Aucune actualité disponible</p>';
        return;
    }

    news.forEach(item => {
        const article = document.createElement('article');
        article.className = 'news-item mb-3';
        
        const date = new Date(item.date);
        const formattedDate = new Intl.DateTimeFormat('fr-FR', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        }).format(date);

        article.innerHTML = `
            <h6 class="news-title">${item.title}</h6>
            <p class="news-content">${item.content}</p>
            <div class="news-meta">
                <small class="text-muted">
                    <i class="bi bi-clock"></i> ${formattedDate}
                    ${item.source ? `<span class="ms-2"><i class="bi bi-link-45deg"></i> ${item.source}</span>` : ''}
                </small>
            </div>
        `;

        if (item.link) {
            article.querySelector('.news-title').innerHTML = `
                <a href="${item.link}" target="_blank" rel="noopener noreferrer">${item.title}</a>
            `;
        }

        newsContainer.appendChild(article);
    });
}

// Fonction pour afficher une erreur
function displayError() {
    const newsContainer = document.getElementById('news');
    newsContainer.innerHTML = `
        <div class="alert alert-warning" role="alert">
            <i class="bi bi-exclamation-triangle-fill"></i>
            Impossible de charger les actualités
            <button type="button" class="btn btn-sm btn-warning ms-2" onclick="loadNews()">
                Réessayer
            </button>
        </div>
    `;
}

// Chargement initial des actualités
loadNews();

// Rafraîchissement des actualités toutes les 5 minutes
setInterval(loadNews, 300000); 