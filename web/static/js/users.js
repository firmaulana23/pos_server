// User Management JavaScript

let currentEditingUserId = null;
let currentRoleUserId = null;

// Initialize Users Page
document.addEventListener('DOMContentLoaded', function() {
    console.log('Users script loaded');
    
    // Check authentication
    requireAuth();
    
    // Load users
    loadUsers();
    
    // Setup form handlers
    setupFormHandlers();
});

// Setup form event handlers
function setupFormHandlers() {
    const userForm = document.getElementById('userForm');
    const roleForm = document.getElementById('roleForm');
    
    if (userForm) {
        userForm.addEventListener('submit', handleUserSubmit);
    }
    
    if (roleForm) {
        roleForm.addEventListener('submit', handleRoleChange);
    }
}

// Load all users
async function loadUsers() {
    try {
        console.log('Loading users...');
        const users = await apiCall('/users');
        displayUsers(users);
    } catch (error) {
        console.error('Failed to load users:', error);
        showError('Failed to load users: ' + error.message);
    }
}

// Display users in table
function displayUsers(users) {
    const tbody = document.querySelector('#usersTable tbody');
    tbody.innerHTML = '';
    
    if (!users || users.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="text-center">No users found</td></tr>';
        return;
    }
    
    users.forEach(user => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${user.id}</td>
            <td>${user.username}</td>
            <td>${user.email || '-'}</td>
            <td>${user.full_name}</td>
            <td>
                <span class="role-badge role-${user.role}">${user.role}</span>
            </td>
            <td>${formatDateTime(user.created_at)}</td>
            <td>
                <button onclick="openRoleModal(${user.id}, '${user.username}', '${user.role}')" 
                        class="btn btn-sm btn-secondary" title="Change Role">
                    <i class="fas fa-user-cog"></i>
                </button>
                <button onclick="editUser(${user.id})" 
                        class="btn btn-sm btn-primary" title="Edit User">
                    <i class="fas fa-edit"></i>
                </button>
                <button onclick="deleteUser(${user.id}, '${user.username}')" 
                        class="btn btn-sm btn-danger" title="Delete User">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Open user modal for adding new user
function openUserModal() {
    currentEditingUserId = null;
    document.getElementById('modalTitle').textContent = 'Add User';
    document.getElementById('userForm').reset();
    document.getElementById('password').required = true;
    document.getElementById('userModal').style.display = 'block';
}

// Close user modal
function closeUserModal() {
    document.getElementById('userModal').style.display = 'none';
    currentEditingUserId = null;
}

// Open role change modal
function openRoleModal(userId, username, currentRole) {
    currentRoleUserId = userId;
    document.getElementById('currentUser').value = username;
    document.getElementById('newRole').value = currentRole;
    document.getElementById('roleModal').style.display = 'block';
}

// Close role modal
function closeRoleModal() {
    document.getElementById('roleModal').style.display = 'none';
    currentRoleUserId = null;
}

// Handle user form submission
async function handleUserSubmit(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const userData = Object.fromEntries(formData.entries());
    
    try {
        if (currentEditingUserId) {
            // Update existing user
            await apiCall(`/users/${currentEditingUserId}`, {
                method: 'PUT',
                body: JSON.stringify(userData)
            });
            showSuccess('User updated successfully');
        } else {
            // Create new user
            await apiCall('/auth/register', {
                method: 'POST',
                body: JSON.stringify(userData)
            });
            showSuccess('User created successfully');
        }
        
        closeUserModal();
        loadUsers();
    } catch (error) {
        console.error('Failed to save user:', error);
        showError('Failed to save user: ' + error.message);
    }
}

// Handle role change
async function handleRoleChange(e) {
    e.preventDefault();
    
    const newRole = document.getElementById('newRole').value;
    
    if (!currentRoleUserId || !newRole) {
        showError('Invalid role change request');
        return;
    }
    
    try {
        await apiCall(`/users/${currentRoleUserId}/role`, {
            method: 'PUT',
            body: JSON.stringify({ role: newRole })
        });
        
        showSuccess('User role updated successfully');
        closeRoleModal();
        loadUsers();
    } catch (error) {
        console.error('Failed to update user role:', error);
        showError('Failed to update user role: ' + error.message);
    }
}

// Edit user
async function editUser(userId) {
    try {
        const user = await apiCall(`/users/${userId}`);
        
        currentEditingUserId = userId;
        document.getElementById('modalTitle').textContent = 'Edit User';
        document.getElementById('username').value = user.username;
        document.getElementById('email').value = user.email || '';
        document.getElementById('fullName').value = user.full_name;
        document.getElementById('role').value = user.role;
        document.getElementById('password').value = '';
        document.getElementById('password').required = false;
        
        document.getElementById('userModal').style.display = 'block';
    } catch (error) {
        console.error('Failed to load user for editing:', error);
        showError('Failed to load user details: ' + error.message);
    }
}

// Delete user
async function deleteUser(userId, username) {
    if (!confirm(`Are you sure you want to delete user "${username}"? This action cannot be undone.`)) {
        return;
    }
    
    try {
        await apiCall(`/users/${userId}`, {
            method: 'DELETE'
        });
        
        showSuccess('User deleted successfully');
        loadUsers();
    } catch (error) {
        console.error('Failed to delete user:', error);
        showError('Failed to delete user: ' + error.message);
    }
}

// Show success message
function showSuccess(message) {
    const successDiv = document.getElementById('successMessage');
    successDiv.textContent = message;
    successDiv.style.display = 'block';
    
    setTimeout(() => {
        successDiv.style.display = 'none';
    }, 5000);
    
    // Hide error message if visible
    document.getElementById('errorMessage').style.display = 'none';
}

// Show error message
function showError(message) {
    const errorDiv = document.getElementById('errorMessage');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    
    setTimeout(() => {
        errorDiv.style.display = 'none';
    }, 5000);
    
    // Hide success message if visible
    document.getElementById('successMessage').style.display = 'none';
}

// Close modals when clicking outside
window.onclick = function(event) {
    const userModal = document.getElementById('userModal');
    const roleModal = document.getElementById('roleModal');
    
    if (event.target === userModal) {
        closeUserModal();
    }
    if (event.target === roleModal) {
        closeRoleModal();
    }
};
