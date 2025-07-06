// Transactions Management JavaScript

let transactions = [];
let currentPage = 1;
let totalPages = 1;
const itemsPerPage = 20;
let currentTransactionId = null;

// Initialize page
document.addEventListener('DOMContentLoaded', async function() {
    if (!isAuthenticated()) {
        window.location.href = '/admin/';
        return;
    }
    
    // Set default date range (last 30 days)
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - 30);
    
    document.getElementById('startDate').value = formatDateForInput(startDate);
    document.getElementById('endDate').value = formatDateForInput(endDate);
    
    await loadTransactions();
});

// Load transactions
async function loadTransactions(page = 1) {
    try {
        const params = new URLSearchParams({
            page: page,
            limit: itemsPerPage
        });
        
        // Add filters
        const status = document.getElementById('statusFilter').value;
        const paymentMethod = document.getElementById('paymentMethodFilter').value;
        const startDate = document.getElementById('startDate').value;
        const endDate = document.getElementById('endDate').value;
        
        if (status) params.append('status', status);
        if (paymentMethod) params.append('payment_method', paymentMethod);
        if (startDate) params.append('start_date', startDate);
        if (endDate) params.append('end_date', endDate);
        
        const response = await apiCall(`/transactions?${params.toString()}`);
        transactions = response.data || [];
        currentPage = response.page || 1;
        totalPages = Math.ceil((response.total || 0) / itemsPerPage);
        
        displayTransactions();
        displayPagination();
    } catch (error) {
        console.error('Failed to load transactions:', error);
        showError('Failed to load transactions: ' + error.message);
    }
}

// Display transactions
function displayTransactions() {
    const tbody = document.getElementById('transactionsTable');
    
    if (transactions.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" class="no-data">No transactions found</td></tr>';
        return;
    }
    
    tbody.innerHTML = transactions.map(transaction => {
        const itemCount = transaction.items ? transaction.items.length : 0;
        return `
            <tr>
                <td>${transaction.transaction_no}</td>
                <td>${formatDateTime(transaction.created_at)}</td>
                <td>${transaction.user ? transaction.user.username : '-'}</td>
                <td>${itemCount} item(s)</td>
                <td>${formatCurrency(transaction.total)}</td>
                <td>${formatPaymentMethod(transaction.payment_method)}</td>
                <td>
                    <span class="status ${transaction.status}">
                        ${transaction.status.charAt(0).toUpperCase() + transaction.status.slice(1)}
                    </span>
                </td>
                <td>
                    <button class="btn btn-sm btn-primary" onclick="viewTransaction(${transaction.id})">
                        <i class="fas fa-eye"></i>
                    </button>
                    ${transaction.status === 'pending' ? `
                        <button class="btn btn-sm btn-success" onclick="markAsPaid(${transaction.id})">
                            <i class="fas fa-check"></i>
                        </button>
                    ` : ''}
                </td>
            </tr>
        `;
    }).join('');
}

// Display pagination
function displayPagination() {
    const pagination = document.getElementById('pagination');
    
    if (totalPages <= 1) {
        pagination.innerHTML = '';
        return;
    }
    
    let paginationHTML = '';
    
    // Previous button
    if (currentPage > 1) {
        paginationHTML += `<button class="btn btn-sm" onclick="loadTransactions(${currentPage - 1})">Previous</button>`;
    }
    
    // Page numbers
    for (let i = Math.max(1, currentPage - 2); i <= Math.min(totalPages, currentPage + 2); i++) {
        paginationHTML += `<button class="btn btn-sm ${i === currentPage ? 'btn-primary' : ''}" onclick="loadTransactions(${i})">${i}</button>`;
    }
    
    // Next button
    if (currentPage < totalPages) {
        paginationHTML += `<button class="btn btn-sm" onclick="loadTransactions(${currentPage + 1})">Next</button>`;
    }
    
    pagination.innerHTML = paginationHTML;
}

// Filter transactions
function filterTransactions() {
    currentPage = 1;
    loadTransactions();
}

// View transaction details
async function viewTransaction(transactionId) {
    try {
        const transaction = await apiCall(`/transactions/${transactionId}`);
        displayTransactionDetails(transaction);
        document.getElementById('transactionModal').style.display = 'block';
    } catch (error) {
        showError('Failed to load transaction details: ' + error.message);
    }
}

// Display transaction details in modal
function displayTransactionDetails(transaction) {
    document.getElementById('transactionNo').textContent = transaction.transaction_no;
    document.getElementById('transactionDate').textContent = formatDateTime(transaction.created_at);
    document.getElementById('transactionCashier').textContent = transaction.user ? transaction.user.username : '-';
    document.getElementById('transactionPaymentMethod').textContent = formatPaymentMethod(transaction.payment_method);
    
    const statusElement = document.getElementById('transactionStatus');
    statusElement.textContent = transaction.status.charAt(0).toUpperCase() + transaction.status.slice(1);
    statusElement.className = `status ${transaction.status}`;
    
    // Display items
    const itemsTableBody = document.getElementById('transactionItems');
    itemsTableBody.innerHTML = transaction.items.map(item => {
        const addOnsText = item.add_ons && item.add_ons.length > 0 
            ? item.add_ons.map(addon => `${addon.add_on.name} (${addon.quantity}x)`).join(', ')
            : '-';
        
        return `
            <tr>
                <td>${item.menu_item.name}</td>
                <td>${item.quantity}</td>
                <td>${formatCurrency(item.unit_price)}</td>
                <td>${addOnsText}</td>
                <td>${formatCurrency(item.total_price)}</td>
            </tr>
        `;
    }).join('');
    
    // Display summary
    document.getElementById('transactionSubtotal').textContent = formatCurrency(transaction.sub_total);
    document.getElementById('transactionTax').textContent = formatCurrency(transaction.tax || 0);
    document.getElementById('transactionDiscount').textContent = formatCurrency(transaction.discount || 0);
    document.getElementById('transactionTotal').textContent = formatCurrency(transaction.total);
}

// Mark transaction as paid
async function markAsPaid(transactionId) {
    // Set the current transaction ID and show the payment modal
    currentTransactionId = transactionId;
    document.getElementById('paymentModal').style.display = 'block';
}

// Print receipt
function printReceipt() {
    const transactionNo = document.getElementById('transactionNo').textContent;
    const printWindow = window.open('', '_blank');
    
    const receiptContent = generateReceiptContent();
    
    printWindow.document.write(`
        <html>
            <head>
                <title>Receipt - ${transactionNo}</title>
                <style>
                    body { font-family: monospace; font-size: 12px; margin: 20px; }
                    .receipt { max-width: 300px; margin: 0 auto; }
                    .header { text-align: center; margin-bottom: 20px; }
                    .items { margin: 20px 0; }
                    .summary { border-top: 1px dashed #000; padding-top: 10px; }
                    .total { font-weight: bold; border-top: 1px solid #000; padding-top: 5px; }
                </style>
            </head>
            <body>
                ${receiptContent}
            </body>
        </html>
    `);
    
    printWindow.document.close();
    printWindow.print();
}

// Generate receipt content
function generateReceiptContent() {
    const transactionNo = document.getElementById('transactionNo').textContent;
    const date = document.getElementById('transactionDate').textContent;
    const cashier = document.getElementById('transactionCashier').textContent;
    const paymentMethod = document.getElementById('transactionPaymentMethod').textContent;
    
    let receiptHTML = `
        <div class="receipt">
            <div class="header">
                <h2>Coffee Shop POS</h2>
                <p>Receipt</p>
                <p>Transaction: ${transactionNo}</p>
                <p>Date: ${date}</p>
                <p>Cashier: ${cashier}</p>
            </div>
            <div class="items">
    `;
    
    // Add items
    const itemRows = document.querySelectorAll('#transactionItems tr');
    itemRows.forEach(row => {
        const cells = row.querySelectorAll('td');
        receiptHTML += `
            <div style="display: flex; justify-content: space-between; margin: 5px 0;">
                <span>${cells[0].textContent} (${cells[1].textContent})</span>
                <span>${cells[4].textContent}</span>
            </div>
        `;
        if (cells[3].textContent !== '-') {
            receiptHTML += `<div style="font-size: 10px; color: #666; margin-left: 10px;">+ ${cells[3].textContent}</div>`;
        }
    });
    
    receiptHTML += `
            </div>
            <div class="summary">
                <div style="display: flex; justify-content: space-between;">
                    <span>Subtotal:</span>
                    <span>${document.getElementById('transactionSubtotal').textContent}</span>
                </div>
                <div style="display: flex; justify-content: space-between;">
                    <span>Tax:</span>
                    <span>${document.getElementById('transactionTax').textContent}</span>
                </div>
                <div style="display: flex; justify-content: space-between;">
                    <span>Discount:</span>
                    <span>${document.getElementById('transactionDiscount').textContent}</span>
                </div>
                <div class="total" style="display: flex; justify-content: space-between;">
                    <span>Total:</span>
                    <span>${document.getElementById('transactionTotal').textContent}</span>
                </div>
            </div>
            <div style="text-align: center; margin-top: 20px; font-size: 10px;">
                <p>Payment Method: ${paymentMethod}</p>
                <p>Thank you for your visit!</p>
            </div>
        </div>
    `;
    
    return receiptHTML;
}

// Close modal
function closeTransactionModal() {
    document.getElementById('transactionModal').style.display = 'none';
}

// Show payment method selection modal
function showPaymentMethodModal() {
    return new Promise((resolve) => {
        // Create modal HTML
        const modalHTML = `
            <div id="paymentModal" class="modal" style="display: block;">
                <div class="modal-content">
                    <div class="modal-header">
                        <h3>Select Payment Method</h3>
                        <span class="close" onclick="closePaymentModal()">&times;</span>
                    </div>
                    <div class="modal-body">
                        <div class="payment-methods">
                            <button class="payment-method-btn" onclick="selectPaymentMethod('cash')">
                                <i class="fas fa-money-bill"></i> Cash
                            </button>
                            <button class="payment-method-btn" onclick="selectPaymentMethod('card')">
                                <i class="fas fa-credit-card"></i> Credit Card
                            </button>
                            <button class="payment-method-btn" onclick="selectPaymentMethod('digital_wallet')">
                                <i class="fas fa-wallet"></i> Digital Wallet
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        // Add modal to body
        document.body.insertAdjacentHTML('beforeend', modalHTML);
        
        // Store resolve function globally so buttons can access it
        window.paymentModalResolve = resolve;
    });
}

// Close payment modal
function closePaymentModal() {
    document.getElementById('paymentModal').style.display = 'none';
    currentTransactionId = null;
}

// Close modal when clicking outside of it
window.onclick = function(event) {
    const modal = document.getElementById('paymentModal');
    if (event.target === modal) {
        closePaymentModal();
    }
}

// Process payment with selected method
async function processPayment(paymentMethod) {
    if (!currentTransactionId) {
        showError('No transaction selected');
        return;
    }

    try {
        await apiCall(`/transactions/${currentTransactionId}/pay`, {
            method: 'PUT',
            body: JSON.stringify({
                payment_method: paymentMethod
            })
        });

        // Close the modal
        closePaymentModal();
        
        // Refresh the transactions list
        await loadTransactions();
        
        showSuccess('Transaction marked as paid successfully!');
    } catch (error) {
        console.error('Error processing payment:', error);
        showError('Error processing payment: ' + error.message);
    }
}

// Utility functions
function formatDateTime(dateString) {
    return new Date(dateString).toLocaleString('id-ID');
}

function formatDateForInput(date) {
    return date.toISOString().split('T')[0];
}

function formatCurrency(amount) {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        minimumFractionDigits: 0
    }).format(amount);
}

function formatPaymentMethod(method) {
    const methods = {
        'cash': 'Cash',
        'card': 'Credit Card',
        'digital_wallet': 'Digital Wallet'
    };
    return methods[method] || method;
}

function showSuccess(message) {
    alert(message);
}

function showError(message) {
    alert(message);
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('transactionModal');
    if (event.target === modal) {
        closeTransactionModal();
    }
}
