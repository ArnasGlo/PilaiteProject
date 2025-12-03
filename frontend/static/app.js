// base API url: if frontend served from same origin, use relative paths
const apiBase = '';

// helpers
async function doFetch(path, opts = {}) {
    const url = apiBase + path;
    const defaultHeaders = { 'Accept': 'application/json' };
    if (opts.body && !(opts.body instanceof FormData)) {
        defaultHeaders['Content-Type'] = 'application/json';
        opts.body = JSON.stringify(opts.body);
    }
    opts.headers = Object.assign(defaultHeaders, opts.headers || {});
    const res = await fetch(url, opts);
    const contentType = res.headers.get('content-type') || '';
    const body = contentType.includes('application/json') ? await res.json().catch(()=>null) : await res.text().catch(()=>null);
    return { res, body };
}

async function loadSpots(){
    console.log('loadSpots() called');
    try {
        const {res, body} = await doFetch('/spots/');

        if (!res.ok) {
            console.error('Failed to load spots. Status:', res.status);
            return;
        }
        renderSpots(body);
    } catch (error) {
        console.error('Error in loadSpots:', error);
    }
}

function renderSpots(spots) {
    const container = document.querySelector('.blocks-container');
    container.innerHTML = ''; // Clear existing

    spots.forEach((spot, index) => {
        const li = document.createElement('li');
        li.className = 'spot-card';
        li.innerHTML = `
            <a href="spot-details.html?id=${spot.id}" class="card-link">
                <div class="card-image">
                    <img src="${spot.image_url || 'images/placeholder.jpg'}" alt="${spot.name}">
                    <span class="category-tag">${spot.category}</span>
                </div>
                <div class="card-content">
                    <h2 class="card-title">${spot.name}</h2>
                    <div class="card-address">
                        <span class="location-icon">📍</span>
                        <span>${spot.address}</span>
                    </div>
                </div>
            </a>
        `;
        container.appendChild(li);
    });
}

const dropdownLinks = document.querySelectorAll('.dropdown-content a');
if (dropdownLinks.length > 0) {
    dropdownLinks.forEach(link => {
        link.addEventListener('click', async (e) => {
            e.preventDefault();

            const categoryMap = {
                'Gamta': 'Gamta',
                'Lauko treniruokliai': 'Lauko_treniruokliai',
                'Slaptos vietos': 'Slaptos_vietos'
            };

            const categoryLT = e.target.textContent.trim();
            const category = categoryMap[categoryLT];

            if (category === 'Slaptos_vietos') {
                const { res, body } = await doFetch('/spots/category/Slaptos_vietos', { credentials: 'include' });
                if (res.ok) renderSpots(body);
                else alert('Please login to view secret spots');
            } else {
                const { res, body } = await doFetch(`/spots/public/category/${category}`);
                if (res.ok) renderSpots(body);
            }
        });
    });
}

    console.log('Setting up DOMContentLoaded listener');
    document.addEventListener('DOMContentLoaded', () => {
    loadSpots();
});

// UI wiring
document.getElementById('btnHealth').addEventListener('click', async () => {
    const { res, body } = await doFetch('/public/health');
    document.getElementById('healthOut').textContent = JSON.stringify({ status: res.status, body }, null, 2);
});

// Public spots
document.getElementById('btnPublicSpots').addEventListener('click', async () => {
    const { res, body } = await doFetch('/spots/');
    document.getElementById('spotsOut').textContent = JSON.stringify({ status: res.status, body }, null, 2);
});

document.getElementById('btnNatureSpots').addEventListener('click', async () => {
    const { res, body } = await doFetch('/spots/public/category/Gamta');
    document.getElementById('spotsOut').textContent = JSON.stringify({ status: res.status, body }, null, 2);
});

document.getElementById('btnTreniruokliaiSpots').addEventListener('click', async () => {
    const { res, body } = await doFetch('/spots/public/category/Lauko_treniruokliai');
    document.getElementById('spotsOut').textContent = JSON.stringify({ status: res.status, body }, null, 2);
});

document.getElementById('btnSecretSpots').addEventListener('click', async () => {
    const { res, body } = await doFetch('/spots/category/Slaptos_vietos', { credentials: 'include' });
    document.getElementById('spotsOut').textContent = JSON.stringify({ status: res.status, body }, null, 2);
});
// Protected spots (browser must send cookie)
document.getElementById('btnProtectedSpots').addEventListener('click', async () => {
    // include credentials so cookie is sent to server
    const { res, body } = await doFetch('/spots/all', { credentials: 'include' });
    document.getElementById('spotsOut').textContent = JSON.stringify({ status: res.status, body }, null, 2);
});

// Register
document.getElementById('btnRegister').addEventListener('click', async () => {
    const email = document.getElementById('regEmail').value;
    const password = document.getElementById('regPassword').value;
    const confirm = document.getElementById('regConfirm').value;
    const { res, body } = await doFetch('/register', { method: 'POST', body: { email, password, confirm_password: confirm }, credentials: 'include' });
    document.getElementById('regMsg').textContent = res.ok ? 'Registered: ' + (body.message||'') : 'Error: ' + (body.error || JSON.stringify(body));
});

// Login
document.getElementById('btnLogin').addEventListener('click', async () => {
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;
    // credentials: 'include' so session cookie set by server is persisted by browser
    const { res, body } = await doFetch('/login', { method: 'POST', body: { email, password }, credentials: 'include' });
    document.getElementById('loginMsg').textContent = res.ok ? 'Logged in' : 'Login failed: ' + (body.error || JSON.stringify(body));
});

// /me
document.getElementById('btnMe').addEventListener('click', async () => {
    const { res, body } = await doFetch('/me', { credentials: 'include' });
    document.getElementById('meOut').textContent = JSON.stringify({ status: res.status, body }, null, 2);
});

// Logout
document.getElementById('btnLogout').addEventListener('click', async () => {
    const { res, body } = await doFetch('/logout', { credentials: 'include' });
    alert(res.ok ? 'Logged out' : 'Logout failed: ' + JSON.stringify(body));
});