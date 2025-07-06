// POS System Variables
let cart = [];
let categories = [];
let menuItems = [];
let addOns = [];
let currentCategory = null;
let currentItemForAddOns = null;
let isLoading = false; // Add loading state

// Initialize POS
document.addEventListener('DOMContentLoaded', async function() {
    await loadCategories();
    await loadAddOns();
    await loadPaymentMethods();
    updateCartDisplay();
});

// Load categories
async function loadCategories() {
    try {
        const response = await apiCall('/public/menu/categories');
        // Handle both direct array and wrapped response
        categories = Array.isArray(response) ? response : (response.data || []);
        displayCategories();
        if (categories.length > 0) {
            selectCategory(categories[0].id);
        }
    } catch (error) {
        console.error('Failed to load categories:', error);
        showError('Failed to load categories: ' + error.message);
    }
}

// Display categories
function displayCategories() {
    const categoryTabs = document.getElementById('categoryTabs');
    categoryTabs.innerHTML = categories.map(category => 
        `<div class="category-tab" onclick="selectCategory(${category.id})">${category.name}</div>`
    ).join('');
}

// Select category
async function selectCategory(categoryId) {
    if (isLoading || currentCategory === categoryId) return; // Prevent duplicate calls
    
    currentCategory = categoryId;
    
    // Update active tab
    document.querySelectorAll('.category-tab').forEach((tab, index) => {
        tab.classList.toggle('active', categories[index].id === categoryId);
    });
    
    await loadMenuItems(categoryId);
}

// Load menu items
async function loadMenuItems(categoryId) {
    if (isLoading) return; // Prevent concurrent calls
    
    isLoading = true;
    
    // Show loading state
    const menuGrid = document.getElementById('menuGrid');
    menuGrid.innerHTML = '<div class="loading">Loading menu items...</div>';
    
    try {
        const response = await apiCall(`/public/menu/items?category_id=${categoryId}`);
        // Handle wrapped response format
        menuItems = response.data || [];
        displayMenuItems();
    } catch (error) {
        console.error('Failed to load menu items:', error);
        showError('Failed to load menu items: ' + error.message);
        menuItems = []; // Set empty array on error
        displayMenuItems();
    } finally {
        isLoading = false;
    }
}

// Display menu items
function displayMenuItems() {
    const menuGrid = document.getElementById('menuGrid');
    
    if (!menuItems || menuItems.length === 0) {
        menuGrid.innerHTML = '<div class="no-items">No menu items found in this category.</div>';
        return;
    }
    
    menuGrid.innerHTML = menuItems.map(item => {
        const margin = item.price > 0 ? ((item.price - item.cogs) / item.price * 100).toFixed(1) : 0;
        return `
            <div class="menu-item" onclick="addToCart(${item.id})">
                <h4>${item.name}</h4>
                <p class="description">${item.description || ''}</p>
                <p class="price">${formatCurrency(item.price)}</p>
                <p class="margin">Margin: ${margin}%</p>
            </div>
        `;
    }).join('');
}

// Load add-ons
async function loadAddOns() {
    try {
        const response = await apiCall('/public/add-ons?available=true');
        // Handle wrapped response format
        addOns = Array.isArray(response) ? response : (response.data || []);
    } catch (error) {
        console.error('Failed to load add-ons:', error);
        showError('Failed to load add-ons: ' + error.message);
        addOns = []; // Set empty array on error
    }
}

// Add item to cart
function addToCart(menuItemId) {
    const menuItem = menuItems.find(item => item.id === menuItemId);
    if (!menuItem) return;

    if (addOns.length > 0) {
        currentItemForAddOns = menuItem;
        showAddOnModal();
    } else {
        addItemToCart(menuItem, []);
    }
}

// Add item to cart with add-ons
function addItemToCart(menuItem, selectedAddOns) {
    const existingItem = cart.find(item => 
        item.menuItem.id === menuItem.id && 
        JSON.stringify(item.addOns) === JSON.stringify(selectedAddOns)
    );

    if (existingItem) {
        existingItem.quantity++;
    } else {
        cart.push({
            menuItem: menuItem,
            quantity: 1,
            addOns: selectedAddOns.map(addon => ({
                ...addon,
                quantity: addon.quantity || 1
            }))
        });
    }

    updateCartDisplay();
}

// Show add-on modal
function showAddOnModal() {
    const modal = document.getElementById('addOnModal');
    const addOnsList = document.getElementById('addOnsList');
    
    addOnsList.innerHTML = addOns.map(addon => `
        <div class="addon-item">
            <div class="addon-info">
                <h5>${addon.name}</h5>
                <p>${addon.description}</p>
                <span class="addon-price">${formatCurrency(addon.price)}</span>
            </div>
            <div class="addon-controls">
                <button type="button" onclick="decreaseAddonQuantity(${addon.id})" class="quantity-btn">-</button>
                <input type="number" id="addon-${addon.id}" class="addon-quantity" value="0" min="0">
                <button type="button" onclick="increaseAddonQuantity(${addon.id})" class="quantity-btn">+</button>
            </div>
        </div>
    `).join('');
    
    modal.style.display = 'block';
    
    // Close modal handlers
    const closeBtn = modal.querySelector('.close');
    closeBtn.onclick = closeAddOnModal;
    
    window.onclick = (event) => {
        if (event.target === modal) {
            closeAddOnModal();
        }
    };
}

// Close add-on modal
function closeAddOnModal() {
    const modal = document.getElementById('addOnModal');
    modal.style.display = 'none';
    currentItemForAddOns = null;
}

// Increase add-on quantity
function increaseAddonQuantity(addonId) {
    const input = document.getElementById(`addon-${addonId}`);
    input.value = parseInt(input.value) + 1;
}

// Decrease add-on quantity
function decreaseAddonQuantity(addonId) {
    const input = document.getElementById(`addon-${addonId}`);
    const currentValue = parseInt(input.value);
    if (currentValue > 0) {
        input.value = currentValue - 1;
    }
}

// Confirm add-ons
function confirmAddOns() {
    const selectedAddOns = [];
    
    addOns.forEach(addon => {
        const quantity = parseInt(document.getElementById(`addon-${addon.id}`).value);
        if (quantity > 0) {
            selectedAddOns.push({
                ...addon,
                quantity: quantity
            });
        }
    });
    
    addItemToCart(currentItemForAddOns, selectedAddOns);
    closeAddOnModal();
}

// Update cart display
function updateCartDisplay() {
    const cartItems = document.getElementById('cartItems');
    const subtotalEl = document.getElementById('subtotal');
    const taxEl = document.getElementById('tax');
    const discountEl = document.getElementById('discount');
    const totalEl = document.getElementById('total');
    
    // Display cart items
    cartItems.innerHTML = cart.map((item, index) => {
        const addOnsText = item.addOns.length > 0 
            ? item.addOns.map(addon => `${addon.name} (${addon.quantity}x)`).join(', ')
            : '';
        
        const itemTotal = calculateItemTotal(item);
        
        return `
            <div class="cart-item">
                <div class="cart-item-info">
                    <h5>${item.menuItem.name}</h5>
                    ${addOnsText ? `<div class="cart-item-addons">Add-ons: ${addOnsText}</div>` : ''}
                    <div>${formatCurrency(itemTotal)}</div>
                </div>
                <div class="cart-item-controls">
                    <button onclick="decreaseQuantity(${index})" class="quantity-btn">-</button>
                    <span>${item.quantity}</span>
                    <button onclick="increaseQuantity(${index})" class="quantity-btn">+</button>
                    <button onclick="removeFromCart(${index})" class="btn btn-danger btn-sm">×</button>
                </div>
            </div>
        `;
    }).join('');
    
    // Calculate totals
    const subtotal = calculateSubtotal();
    const tax = subtotal * 0.1; // 10% tax
    const discount = 0; // No discount for now
    const total = subtotal + tax - discount;
    
    subtotalEl.textContent = formatCurrency(subtotal);
    taxEl.textContent = formatCurrency(tax);
    discountEl.textContent = formatCurrency(discount);
    totalEl.textContent = formatCurrency(total);
}

// Calculate item total including add-ons
function calculateItemTotal(item) {
    let total = item.menuItem.price * item.quantity;
    
    item.addOns.forEach(addon => {
        total += addon.price * addon.quantity * item.quantity;
    });
    
    return total;
}

// Calculate subtotal
function calculateSubtotal() {
    return cart.reduce((total, item) => total + calculateItemTotal(item), 0);
}

// Increase quantity
function increaseQuantity(index) {
    cart[index].quantity++;
    updateCartDisplay();
}

// Decrease quantity
function decreaseQuantity(index) {
    if (cart[index].quantity > 1) {
        cart[index].quantity--;
    } else {
        cart.splice(index, 1);
    }
    updateCartDisplay();
}

// Remove from cart
function removeFromCart(index) {
    cart.splice(index, 1);
    updateCartDisplay();
}

// Clear cart
function clearCart() {
    if (confirm('Are you sure you want to clear the cart?')) {
        cart = [];
        updateCartDisplay();
    }
}

// Save transaction
async function saveTransaction() {
    if (cart.length === 0) {
        showError('Cart is empty');
        return;
    }
    
    try {
        const transactionData = prepareTransactionData();
        const response = await apiCall('/transactions', {
            method: 'POST',
            body: JSON.stringify(transactionData)
        });
        
        alert(`Transaction saved successfully! Transaction ID: ${response.transaction_no}`);
        clearCart();
    } catch (error) {
        showError('Failed to save transaction: ' + error.message);
    }
}

// Pay transaction
function payTransaction() {
    if (cart.length === 0) {
        showError('Cart is empty');
        return;
    }
    
    showPaymentModal();
}

// Show payment modal
function showPaymentModal() {
    const modal = document.getElementById('paymentModal');
    const paymentTotal = document.getElementById('paymentTotal');
    
    const subtotal = calculateSubtotal();
    const tax = subtotal * 0.1;
    const total = subtotal + tax;
    
    paymentTotal.textContent = formatCurrency(total);
    modal.style.display = 'block';
    
    // Close modal handlers
    const closeBtn = modal.querySelector('.close');
    closeBtn.onclick = closePaymentModal;
    
    window.onclick = (event) => {
        if (event.target === modal) {
            closePaymentModal();
        }
    };
}

// Close payment modal
function closePaymentModal() {
    const modal = document.getElementById('paymentModal');
    modal.style.display = 'none';
}

// Confirm payment
async function confirmPayment() {
    const paymentMethod = document.getElementById('paymentMethod').value;
    
    try {
        // First save the transaction
        const transactionData = prepareTransactionData();
        const transaction = await apiCall('/transactions', {
            method: 'POST',
            body: JSON.stringify(transactionData)
        });
        
        // Then process payment
        await apiCall(`/transactions/${transaction.id}/pay`, {
            method: 'PUT',
            body: JSON.stringify({ payment_method: paymentMethod })
        });
        
        alert(`Payment processed successfully! Transaction ID: ${transaction.transaction_no}`);
        clearCart();
        closePaymentModal();
    } catch (error) {
        showError('Failed to process payment: ' + error.message);
    }
}

// Prepare transaction data
function prepareTransactionData() {
    const subtotal = calculateSubtotal();
    const tax = subtotal * 0.1;
    
    const items = cart.map(item => ({
        menu_item_id: item.menuItem.id,
        quantity: item.quantity,
        add_ons: item.addOns.map(addon => ({
            add_on_id: addon.id,
            quantity: addon.quantity
        }))
    }));
    
    return {
        items: items,
        tax: tax,
        discount: 0
    };
}

// Load payment methods
async function loadPaymentMethods() {
    try {
        const paymentMethods = await apiCall('/public/payment-methods');
        const select = document.getElementById('paymentMethod');
        
        select.innerHTML = paymentMethods.map(method => 
            `<option value="${method.code}">${method.name}</option>`
        ).join('');
    } catch (error) {
        console.error('Failed to load payment methods:', error);
    }
}
