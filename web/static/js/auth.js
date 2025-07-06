// API Base URL
const API_BASE = '/api/v1';

// Auth utilities
function getToken() {
    return localStorage.getItem('token');
}

function setToken(token) {
    localStorage.setItem('token', token);
}

function removeToken() {
    localStorage.removeItem('token');
}

function getHeaders() {
    const token = getToken();
    return {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` })
    };
}

// API wrapper
async function apiCall(url, options = {}) {
    console.log('apiCall called with URL:', url, 'Options:', options);
    const response = await fetch(`${API_BASE}${url}`, {
        ...options,
        headers: {
            ...getHeaders(),
            ...options.headers
        }
    });

    console.log('Response status:', response.status);
    console.log('Response headers:', response.headers);

    if (!response.ok) {
        if (response.status === 401) {
            console.log('Unauthorized, removing token and redirecting');
            removeToken();
            window.location.href = '/admin/';
            return;
        }
        
        const error = await response.json();
        console.error('API error:', error);
        throw new Error(error.error || 'An error occurred');
    }

    const data = await response.json();
    console.log('API response data:', data);
    return data;
}

// Login function
async function login(username, password) {
    try {
        const response = await apiCall('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });

        setToken(response.token);
        return response;
    } catch (error) {
        throw error;
    }
}

// Logout function
function logout() {
    removeToken();
    window.location.href = '/admin/';
}

// Check if user is authenticated
function isAuthenticated() {
    return !!getToken();
}

// Redirect to dashboard if already authenticated
function redirectIfAuthenticated() {
    if (isAuthenticated()) {
        window.location.href = '/admin/dashboard';
    }
}

// Require authentication
function requireAuth() {
    if (!isAuthenticated()) {
        window.location.href = '/admin/';
    }
}

// Login form handler
document.addEventListener('DOMContentLoaded', function() {
    console.log('auth.js loaded');
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        console.log('Login form found, checking if already authenticated');
        redirectIfAuthenticated();
        
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            
            try {
                await login(username, password);
                window.location.href = '/admin/dashboard';
            } catch (error) {
                showError(error.message);
            }
        });
    } else {
        console.log('No login form found, requiring auth');
        requireAuth();
    }
});

// Error modal
function showError(message) {
    const modal = document.getElementById('errorModal');
    const errorMessage = document.getElementById('errorMessage');
    
    if (modal && errorMessage) {
        errorMessage.textContent = message;
        modal.style.display = 'block';
        
        // Close modal handlers
        const closeBtn = modal.querySelector('.close');
        closeBtn.onclick = () => modal.style.display = 'none';
        
        window.onclick = (event) => {
            if (event.target === modal) {
                modal.style.display = 'none';
            }
        };
    } else {
        alert(message);
    }
}

// Format currency
function formatCurrency(amount) {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        minimumFractionDigits: 0
    }).format(amount);
}

// Format date
function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('id-ID');
}

// Format datetime
function formatDateTime(dateString) {
    return new Date(dateString).toLocaleString('id-ID');
}
