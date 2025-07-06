// Menu Management JavaScript

let categories = [];
let menuItems = [];
let isEditingCategory = false;
let isEditingMenuItem = false;

// Initialize page
document.addEventListener('DOMContentLoaded', async function() {
    if (!isAuthenticated()) {
        window.location.href = '/admin/';
        return;
    }
    
    await loadCategories();
    await loadMenuItems();
    populateCategorySelects();
});

// Load categories
async function loadCategories() {
    try {
        const response = await apiCall('/menu/categories');
        categories = response.data || [];
        displayCategories();
    } catch (error) {
        console.error('Failed to load categories:', error);
        showError('Failed to load categories: ' + error.message);
    }
}

// Display categories
function displayCategories() {
    const tbody = document.getElementById('categoriesTable');
    
    if (categories.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="no-data">No categories found</td></tr>';
        return;
    }
    
    tbody.innerHTML = categories.map(category => `
        <tr>
            <td>${category.id}</td>
            <td>${category.name}</td>
            <td>${category.description || '-'}</td>
            <td>${category.menu_items ? category.menu_items.length : 0}</td>
            <td>
                <button class="btn btn-sm btn-primary" onclick="editCategory(${category.id})">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="btn btn-sm btn-danger" onclick="deleteCategory(${category.id})">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        </tr>
    `).join('');
}

// Load menu items
async function loadMenuItems() {
    try {
        const response = await apiCall('/menu/items');
        menuItems = response.data || [];
        displayMenuItems();
    } catch (error) {
        console.error('Failed to load menu items:', error);
        showError('Failed to load menu items: ' + error.message);
    }
}

// Display menu items
function displayMenuItems() {
    const tbody = document.getElementById('menuItemsTable');
    
    if (menuItems.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" class="no-data">No menu items found</td></tr>';
        return;
    }
    
    tbody.innerHTML = menuItems.map(item => {
        const margin = item.price > 0 ? ((item.price - item.cogs) / item.price * 100).toFixed(1) : 0;
        return `
            <tr>
                <td>${item.id}</td>
                <td>${item.name}</td>
                <td>${item.category ? item.category.name : '-'}</td>
                <td>${formatCurrency(item.price)}</td>
                <td>${formatCurrency(item.cogs)}</td>
                <td>${margin}%</td>
                <td>
                    <span class="status ${item.is_available ? 'available' : 'unavailable'}">
                        ${item.is_available ? 'Available' : 'Unavailable'}
                    </span>
                </td>
                <td>
                    <button class="btn btn-sm btn-primary" onclick="editMenuItem(${item.id})">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="deleteMenuItem(${item.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            </tr>
        `;
    }).join('');
}

// Populate category selects
function populateCategorySelects() {
    const categoryFilter = document.getElementById('categoryFilter');
    const menuItemCategory = document.getElementById('menuItemCategory');
    
    const categoryOptions = categories.map(cat => 
        `<option value="${cat.id}">${cat.name}</option>`
    ).join('');
    
    categoryFilter.innerHTML = '<option value="">All Categories</option>' + categoryOptions;
    menuItemCategory.innerHTML = '<option value="">Select Category</option>' + categoryOptions;
}

// Filter menu items
function filterMenuItems() {
    const categoryFilter = document.getElementById('categoryFilter').value;
    const availabilityFilter = document.getElementById('availabilityFilter').value;
    
    let filteredItems = menuItems;
    
    if (categoryFilter) {
        filteredItems = filteredItems.filter(item => item.category_id == categoryFilter);
    }
    
    if (availabilityFilter !== '') {
        const isAvailable = availabilityFilter === 'true';
        filteredItems = filteredItems.filter(item => item.is_available === isAvailable);
    }
    
    const tbody = document.getElementById('menuItemsTable');
    
    if (filteredItems.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" class="no-data">No items match the selected filters</td></tr>';
        return;
    }
    
    tbody.innerHTML = filteredItems.map(item => {
        const margin = item.price > 0 ? ((item.price - item.cogs) / item.price * 100).toFixed(1) : 0;
        return `
            <tr>
                <td>${item.id}</td>
                <td>${item.name}</td>
                <td>${item.category ? item.category.name : '-'}</td>
                <td>${formatCurrency(item.price)}</td>
                <td>${formatCurrency(item.cogs)}</td>
                <td>${margin}%</td>
                <td>
                    <span class="status ${item.is_available ? 'available' : 'unavailable'}">
                        ${item.is_available ? 'Available' : 'Unavailable'}
                    </span>
                </td>
                <td>
                    <button class="btn btn-sm btn-primary" onclick="editMenuItem(${item.id})">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="deleteMenuItem(${item.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                </td>
            </tr>
        `;
    }).join('');
}

// Category Modal Functions
function showCreateCategoryModal() {
    isEditingCategory = false;
    document.getElementById('categoryModalTitle').textContent = 'Add Category';
    document.getElementById('categoryForm').reset();
    document.getElementById('categoryId').value = '';
    document.getElementById('categoryModal').style.display = 'block';
}

function editCategory(categoryId) {
    const category = categories.find(c => c.id === categoryId);
    if (!category) return;
    
    isEditingCategory = true;
    document.getElementById('categoryModalTitle').textContent = 'Edit Category';
    document.getElementById('categoryId').value = category.id;
    document.getElementById('categoryName').value = category.name;
    document.getElementById('categoryDescription').value = category.description || '';
    document.getElementById('categoryModal').style.display = 'block';
}

function closeCategoryModal() {
    document.getElementById('categoryModal').style.display = 'none';
}

// Handle category form submission
document.getElementById('categoryForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const categoryData = {
        name: document.getElementById('categoryName').value,
        description: document.getElementById('categoryDescription').value
    };
    
    try {
        if (isEditingCategory) {
            const categoryId = document.getElementById('categoryId').value;
            await apiCall(`/menu/categories/${categoryId}`, {
                method: 'PUT',
                body: JSON.stringify(categoryData)
            });
            showSuccess('Category updated successfully!');
        } else {
            await apiCall('/menu/categories', {
                method: 'POST',
                body: JSON.stringify(categoryData)
            });
            showSuccess('Category created successfully!');
        }
        
        closeCategoryModal();
        await loadCategories();
        populateCategorySelects();
    } catch (error) {
        showError('Error saving category: ' + error.message);
    }
});

// Delete category
async function deleteCategory(categoryId) {
    if (!confirm('Are you sure you want to delete this category? This will also delete all menu items in this category.')) {
        return;
    }
    
    try {
        await apiCall(`/menu/categories/${categoryId}`, {
            method: 'DELETE'
        });
        showSuccess('Category deleted successfully!');
        await loadCategories();
        await loadMenuItems();
        populateCategorySelects();
    } catch (error) {
        showError('Error deleting category: ' + error.message);
    }
}

// Menu Item Modal Functions
function showCreateMenuItemModal() {
    isEditingMenuItem = false;
    document.getElementById('menuItemModalTitle').textContent = 'Add Menu Item';
    document.getElementById('menuItemForm').reset();
    document.getElementById('menuItemId').value = '';
    document.getElementById('menuItemAvailable').checked = true;
    document.getElementById('menuItemModal').style.display = 'block';
}

function editMenuItem(menuItemId) {
    const item = menuItems.find(i => i.id === menuItemId);
    if (!item) return;
    
    isEditingMenuItem = true;
    document.getElementById('menuItemModalTitle').textContent = 'Edit Menu Item';
    document.getElementById('menuItemId').value = item.id;
    document.getElementById('menuItemName').value = item.name;
    document.getElementById('menuItemCategory').value = item.category_id;
    document.getElementById('menuItemDescription').value = item.description || '';
    document.getElementById('menuItemPrice').value = item.price;
    document.getElementById('menuItemCOGS').value = item.cogs;
    document.getElementById('menuItemAvailable').checked = item.is_available;
    document.getElementById('menuItemModal').style.display = 'block';
}

function closeMenuItemModal() {
    document.getElementById('menuItemModal').style.display = 'none';
}

// Handle menu item form submission
document.getElementById('menuItemForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const menuItemData = {
        name: document.getElementById('menuItemName').value,
        category_id: parseInt(document.getElementById('menuItemCategory').value),
        description: document.getElementById('menuItemDescription').value,
        price: parseFloat(document.getElementById('menuItemPrice').value),
        cogs: parseFloat(document.getElementById('menuItemCOGS').value),
        is_available: document.getElementById('menuItemAvailable').checked
    };
    
    try {
        if (isEditingMenuItem) {
            const menuItemId = document.getElementById('menuItemId').value;
            await apiCall(`/menu/items/${menuItemId}`, {
                method: 'PUT',
                body: JSON.stringify(menuItemData)
            });
            showSuccess('Menu item updated successfully!');
        } else {
            await apiCall('/menu/items', {
                method: 'POST',
                body: JSON.stringify(menuItemData)
            });
            showSuccess('Menu item created successfully!');
        }
        
        closeMenuItemModal();
        await loadMenuItems();
    } catch (error) {
        showError('Error saving menu item: ' + error.message);
    }
});

// Delete menu item
async function deleteMenuItem(menuItemId) {
    if (!confirm('Are you sure you want to delete this menu item?')) {
        return;
    }
    
    try {
        await apiCall(`/menu/items/${menuItemId}`, {
            method: 'DELETE'
        });
        showSuccess('Menu item deleted successfully!');
        await loadMenuItems();
    } catch (error) {
        showError('Error deleting menu item: ' + error.message);
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
    // You can implement a toast notification system here
    alert(message);
}

function showError(message) {
    // You can implement a toast notification system here
    alert(message);
}

// Close modals when clicking outside
window.onclick = function(event) {
    const categoryModal = document.getElementById('categoryModal');
    const menuItemModal = document.getElementById('menuItemModal');
    
    if (event.target === categoryModal) {
        closeCategoryModal();
    }
    if (event.target === menuItemModal) {
        closeMenuItemModal();
    }
}
