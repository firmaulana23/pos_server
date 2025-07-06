// Expenses Management JavaScript

let expenses = [];
let isEditingExpense = false;

// Initialize page
document.addEventListener('DOMContentLoaded', async function() {
    if (!isAuthenticated()) {
        window.location.href = '/admin/';
        return;
    }
    
    // Set default date range (current month)
    const now = new Date();
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    
    document.getElementById('startDate').value = formatDateForInput(startOfMonth);
    document.getElementById('endDate').value = formatDateForInput(endOfMonth);
    document.getElementById('expenseDate').value = formatDateForInput(now);
    
    await loadExpenses();
    updateSummary();
});

// Load expenses
async function loadExpenses() {
    try {
        const params = new URLSearchParams();
        
        // Add filters
        const type = document.getElementById('typeFilter').value;
        const startDate = document.getElementById('startDate').value;
        const endDate = document.getElementById('endDate').value;
        const category = document.getElementById('categoryFilter').value;
        
        if (type) params.append('type', type);
        if (startDate) params.append('start_date', startDate);
        if (endDate) params.append('end_date', endDate);
        if (category) params.append('category', category);
        
        const response = await apiCall(`/expenses?${params.toString()}`);
        expenses = response.data || [];
        displayExpenses();
        updateSummary();
    } catch (error) {
        console.error('Failed to load expenses:', error);
        showError('Failed to load expenses: ' + error.message);
    }
}

// Display expenses
function displayExpenses() {
    const tbody = document.getElementById('expensesTable');
    
    if (expenses.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="no-data">No expenses found</td></tr>';
        return;
    }
    
    tbody.innerHTML = expenses.map(expense => `
        <tr>
            <td>${formatDate(expense.date)}</td>
            <td>
                <span class="expense-type ${expense.type}">
                    ${expense.type === 'raw_material' ? 'Raw Materials' : 'Operational'}
                </span>
            </td>
            <td>${expense.category}</td>
            <td>${expense.description}</td>
            <td class="amount">${formatCurrency(expense.amount)}</td>
            <td>${expense.user ? expense.user.username : '-'}</td>
            <td>
                <button class="btn btn-sm btn-primary" onclick="editExpense(${expense.id})">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="btn btn-sm btn-danger" onclick="deleteExpense(${expense.id})">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        </tr>
    `).join('');
}

// Update summary cards
function updateSummary() {
    const rawMaterialsTotal = expenses
        .filter(e => e.type === 'raw_material')
        .reduce((sum, e) => sum + e.amount, 0);
    
    const operationalTotal = expenses
        .filter(e => e.type === 'operational')
        .reduce((sum, e) => sum + e.amount, 0);
    
    const totalExpenses = rawMaterialsTotal + operationalTotal;
    
    document.getElementById('rawMaterialsTotal').textContent = formatCurrency(rawMaterialsTotal);
    document.getElementById('operationalTotal').textContent = formatCurrency(operationalTotal);
    document.getElementById('totalExpenses').textContent = formatCurrency(totalExpenses);
}

// Filter expenses
function filterExpenses() {
    loadExpenses();
}

// Modal Functions
function showCreateExpenseModal() {
    isEditingExpense = false;
    document.getElementById('expenseModalTitle').textContent = 'Add Expense';
    document.getElementById('expenseForm').reset();
    document.getElementById('expenseId').value = '';
    document.getElementById('expenseDate').value = formatDateForInput(new Date());
    document.getElementById('expenseModal').style.display = 'block';
}

function editExpense(expenseId) {
    const expense = expenses.find(e => e.id === expenseId);
    if (!expense) return;
    
    isEditingExpense = true;
    document.getElementById('expenseModalTitle').textContent = 'Edit Expense';
    document.getElementById('expenseId').value = expense.id;
    document.getElementById('expenseType').value = expense.type;
    document.getElementById('expenseCategory').value = expense.category;
    document.getElementById('expenseDescription').value = expense.description;
    document.getElementById('expenseAmount').value = expense.amount;
    document.getElementById('expenseDate').value = formatDateForInput(new Date(expense.date));
    document.getElementById('expenseModal').style.display = 'block';
}

function closeExpenseModal() {
    document.getElementById('expenseModal').style.display = 'none';
}

// Handle form submission
document.getElementById('expenseForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const expenseData = {
        type: document.getElementById('expenseType').value,
        category: document.getElementById('expenseCategory').value,
        description: document.getElementById('expenseDescription').value,
        amount: parseFloat(document.getElementById('expenseAmount').value),
        date: document.getElementById('expenseDate').value + 'T00:00:00Z'
    };
    
    try {
        if (isEditingExpense) {
            const expenseId = document.getElementById('expenseId').value;
            await apiCall(`/expenses/${expenseId}`, {
                method: 'PUT',
                body: JSON.stringify(expenseData)
            });
            showSuccess('Expense updated successfully!');
        } else {
            await apiCall('/expenses', {
                method: 'POST',
                body: JSON.stringify(expenseData)
            });
            showSuccess('Expense created successfully!');
        }
        
        closeExpenseModal();
        await loadExpenses();
    } catch (error) {
        showError('Error saving expense: ' + error.message);
    }
});

// Delete expense
async function deleteExpense(expenseId) {
    if (!confirm('Are you sure you want to delete this expense?')) {
        return;
    }
    
    try {
        await apiCall(`/expenses/${expenseId}`, {
            method: 'DELETE'
        });
        showSuccess('Expense deleted successfully!');
        await loadExpenses();
    } catch (error) {
        showError('Error deleting expense: ' + error.message);
    }
}

// Utility functions
function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('id-ID');
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

function showSuccess(message) {
    alert(message);
}

function showError(message) {
    alert(message);
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('expenseModal');
    if (event.target === modal) {
        closeExpenseModal();
    }
}
