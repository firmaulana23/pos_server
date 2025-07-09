// Add-ons Management JavaScript

let addOns = [];
let menuItems = [];
let isEditingAddOn = false;

// Initialize page
document.addEventListener('DOMContentLoaded', async function() {
    if (!isAuthenticated()) {
        window.location.href = '/admin/';
        return;
    }
    
    await loadMenuItems();
    await loadAddOns();
});

// Load menu items for the dropdown
async function loadMenuItems() {
    try {
        const response = await apiCall('/public/menu/items');
        menuItems = response.data || [];
        populateMenuItemFilter();
        populateMenuItemDropdown();
    } catch (error) {
        console.error('Failed to load menu items:', error);
        showError('Failed to load menu items: ' + error.message);
    }
}

// Populate menu item filter
function populateMenuItemFilter() {
    const select = document.getElementById('menuItemFilter');
    const currentOptions = Array.from(select.options).slice(2); // Keep first 2 options
    currentOptions.forEach(option => option.remove());
    
    menuItems.forEach(item => {
        const option = document.createElement('option');
        option.value = item.id;
        option.textContent = item.name;
        select.appendChild(option);
    });
}

// Populate menu item dropdown in form
function populateMenuItemDropdown() {
    const select = document.getElementById('addOnMenuItemId');
    const currentOptions = Array.from(select.options).slice(1); // Keep first option
    currentOptions.forEach(option => option.remove());
    
    menuItems.forEach(item => {
        const option = document.createElement('option');
        option.value = item.id;
        option.textContent = item.name;
        select.appendChild(option);
    });
}

// Load add-ons
async function loadAddOns() {
    try {
        const response = await apiCall('/add-ons');
        addOns = Array.isArray(response) ? response : (response.data || []);
        displayAddOns();
    } catch (error) {
        console.error('Failed to load add-ons:', error);
        showError('Failed to load add-ons: ' + error.message);
    }
}

// Display add-ons
function displayAddOns() {
    const tbody = document.getElementById('addOnsTable');
    
    if (addOns.length === 0) {
        tbody.innerHTML = '<tr><td colspan="9" class="no-data">No add-ons found</td></tr>';
        return;
    }
    
    tbody.innerHTML = addOns.map(addOn => {
        const margin = addOn.price > 0 ? ((addOn.price - addOn.cogs) / addOn.price * 100).toFixed(1) : 0;
        const menuItemName = addOn.menu_item_id ? 
            (menuItems.find(item => item.id === addOn.menu_item_id)?.name || 'Unknown') : 
            'Global';
        return `
            <tr>
                <td>${addOn.id}</td>
                <td>${addOn.name}</td>
                <td>${addOn.description || '-'}</td>
                <td>
                    <span class="menu-item-tag ${addOn.menu_item_id ? 'specific' : 'global'}">
                        ${menuItemName}
                    </span>
                </td>
                <td>${formatCurrency(addOn.price)}</td>
                <td>${formatCurrency(addOn.cogs)}</td>
                <td>${margin}%</td>
                <td>
                    <span class="status ${addOn.is_available ? 'available' : 'unavailable'}">
                        ${addOn.is_available ? 'Available' : 'Unavailable'}
                    </span>
                </td>
                <td>
                    <button class="btn btn-sm btn-primary" onclick="editAddOn(${addOn.id})">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="deleteAddOn(${addOn.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            </tr>
        `;
    }).join('');
}

// Filter add-ons
function filterAddOns() {
    const availabilityFilter = document.getElementById('availabilityFilter').value;
    const menuItemFilter = document.getElementById('menuItemFilter').value;
    
    let filteredAddOns = addOns;
    
    if (availabilityFilter !== '') {
        const isAvailable = availabilityFilter === 'true';
        filteredAddOns = filteredAddOns.filter(addOn => addOn.is_available === isAvailable);
    }
    
    if (menuItemFilter !== '') {
        if (menuItemFilter === 'global') {
            filteredAddOns = filteredAddOns.filter(addOn => !addOn.menu_item_id);
        } else {
            filteredAddOns = filteredAddOns.filter(addOn => addOn.menu_item_id === parseInt(menuItemFilter));
        }
    }
    
    const tbody = document.getElementById('addOnsTable');
    
    if (filteredAddOns.length === 0) {
        tbody.innerHTML = '<tr><td colspan="9" class="no-data">No add-ons match the selected filter</td></tr>';
        return;
    }
    
    tbody.innerHTML = filteredAddOns.map(addOn => {
        const margin = addOn.price > 0 ? ((addOn.price - addOn.cogs) / addOn.price * 100).toFixed(1) : 0;
        const menuItemName = addOn.menu_item_id ? 
            (menuItems.find(item => item.id === addOn.menu_item_id)?.name || 'Unknown') : 
            'Global';
        return `
            <tr>
                <td>${addOn.id}</td>
                <td>${addOn.name}</td>
                <td>${addOn.description || '-'}</td>
                <td>
                    <span class="menu-item-tag ${addOn.menu_item_id ? 'specific' : 'global'}">
                        ${menuItemName}
                    </span>
                </td>
                <td>${formatCurrency(addOn.price)}</td>
                <td>${formatCurrency(addOn.cogs)}</td>
                <td>${margin}%</td>
                <td>
                    <span class="status ${addOn.is_available ? 'available' : 'unavailable'}">
                        ${addOn.is_available ? 'Available' : 'Unavailable'}
                    </span>
                </td>
                <td>
                    <button class="btn btn-sm btn-primary" onclick="editAddOn(${addOn.id})">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="deleteAddOn(${addOn.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            </tr>
        `;
    }).join('');
}

// Modal Functions
function showCreateAddOnModal() {
    isEditingAddOn = false;
    document.getElementById('addOnModalTitle').textContent = 'Add New Add-on';
    document.getElementById('addOnForm').reset();
    document.getElementById('addOnId').value = '';
    document.getElementById('addOnAvailable').checked = true;
    document.getElementById('addOnModal').style.display = 'block';
}

function editAddOn(addOnId) {
    const addOn = addOns.find(a => a.id === addOnId);
    if (!addOn) return;
    
    isEditingAddOn = true;
    document.getElementById('addOnModalTitle').textContent = 'Edit Add-on';
    document.getElementById('addOnId').value = addOn.id;
    document.getElementById('addOnName').value = addOn.name;
    document.getElementById('addOnDescription').value = addOn.description || '';
    document.getElementById('addOnMenuItemId').value = addOn.menu_item_id || '';
    document.getElementById('addOnPrice').value = addOn.price;
    document.getElementById('addOnCOGS').value = addOn.cogs;
    document.getElementById('addOnAvailable').checked = addOn.is_available;
    document.getElementById('addOnModal').style.display = 'block';
}

function closeAddOnModal() {
    document.getElementById('addOnModal').style.display = 'none';
}

// Handle form submission
document.getElementById('addOnForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const menuItemId = document.getElementById('addOnMenuItemId').value;
    const addOnData = {
        name: document.getElementById('addOnName').value,
        description: document.getElementById('addOnDescription').value,
        price: parseFloat(document.getElementById('addOnPrice').value),
        cogs: parseFloat(document.getElementById('addOnCOGS').value),
        is_available: document.getElementById('addOnAvailable').checked
    };
    
    // Only include menu_item_id if a menu item is selected
    if (menuItemId) {
        addOnData.menu_item_id = parseInt(menuItemId);
    }
    
    try {
        if (isEditingAddOn) {
            const addOnId = document.getElementById('addOnId').value;
            await apiCall(`/add-ons/${addOnId}`, {
                method: 'PUT',
                body: JSON.stringify(addOnData)
            });
            showSuccess('Add-on updated successfully!');
        } else {
            await apiCall('/add-ons', {
                method: 'POST',
                body: JSON.stringify(addOnData)
            });
            showSuccess('Add-on created successfully!');
        }
        
        closeAddOnModal();
        await loadAddOns();
    } catch (error) {
        showError('Error saving add-on: ' + error.message);
    }
});

// Delete add-on
async function deleteAddOn(addOnId) {
    if (!confirm('Are you sure you want to delete this add-on?')) {
        return;
    }
    
    try {
        await apiCall(`/add-ons/${addOnId}`, {
            method: 'DELETE'
        });
        showSuccess('Add-on deleted successfully!');
        await loadAddOns();
    } catch (error) {
        showError('Error deleting add-on: ' + error.message);
    }
}

// Utility functions
function formatCurrency(amount) {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        minimumFractionDigits: 0
    }).format(amount);
}

function showSuccess(message) {
    alert(message);
}

function showError(message) {
    alert(message);
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('addOnModal');
    if (event.target === modal) {
        closeAddOnModal();
    }
}
