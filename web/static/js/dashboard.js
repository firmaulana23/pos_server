// Dashboard Variables
let salesChart, expenseChart;
let dashboardRefreshInterval = null;

// Auto-refresh dashboard every 30 seconds if enabled
function enableAutoRefresh(intervalSeconds = 30) {
    if (dashboardRefreshInterval) {
        clearInterval(dashboardRefreshInterval);
    }
    
    dashboardRefreshInterval = setInterval(() => {
        console.log('Auto-refreshing dashboard...');
        updateDashboard();
    }, intervalSeconds * 1000);
    
    console.log(`Dashboard auto-refresh enabled (${intervalSeconds}s)`);
}

function disableAutoRefresh() {
    if (dashboardRefreshInterval) {
        clearInterval(dashboardRefreshInterval);
        dashboardRefreshInterval = null;
        console.log('Dashboard auto-refresh disabled');
    }
}

// Add manual refresh button functionality
function refreshDashboard() {
    console.log('Manual dashboard refresh triggered');
    updateDashboard();
}

// Window focus refresh - refresh when user returns to dashboard tab
let lastFocusTime = Date.now();
window.addEventListener('focus', () => {
    // Only refresh if it's been more than 5 seconds since last focus
    if (Date.now() - lastFocusTime > 5000) {
        console.log('Window focus - refreshing dashboard');
        updateDashboard();
    }
    lastFocusTime = Date.now();
});

// Listen for storage events (when other tabs make changes)
window.addEventListener('storage', (e) => {
    if (e.key === 'dashboard_refresh_needed') {
        console.log('Detected changes from other tab - refreshing dashboard');
        updateDashboard();
        localStorage.removeItem('dashboard_refresh_needed');
    }
});

// Initialize Dashboard
document.addEventListener('DOMContentLoaded', function() {
    console.log('Dashboard script loaded');
    
    // Leave date inputs empty to show all data initially
    // Users can set specific dates if they want to filter
    
    console.log('About to update dashboard with all data');
    updateDashboard();
    
    // Enable auto-refresh every 30 seconds
    enableAutoRefresh(30);
    
    // Add refresh button functionality if it exists
    const refreshBtn = document.getElementById('refreshDashboard');
    if (refreshBtn) {
        refreshBtn.addEventListener('click', refreshDashboard);
    }
});

// Update dashboard
async function updateDashboard() {
    console.log('updateDashboard called');
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;
    
    console.log('Date range:', startDate, 'to', endDate);
    
    try {
        await Promise.all([
            loadDashboardStats(startDate, endDate),
            loadSalesChart(startDate, endDate),
            loadExpenseChart(startDate, endDate)
        ]);
        console.log('Dashboard updated successfully');
    } catch (error) {
        console.error('Dashboard update error:', error);
        showError('Failed to load dashboard data: ' + error.message);
    }
}

// Load dashboard stats
async function loadDashboardStats(startDate, endDate) {
    console.log('loadDashboardStats called with:', startDate, endDate);
    try {
        const params = new URLSearchParams();
        
        // Only add date parameters if they are provided
        if (startDate && startDate.trim() !== '') {
            params.append('start_date', startDate);
        }
        if (endDate && endDate.trim() !== '') {
            params.append('end_date', endDate);
        }
        
        const queryString = params.toString();
        const url = queryString ? `/dashboard/stats?${queryString}` : '/dashboard/stats';
        
        console.log('Making API call to:', url);
        const stats = await apiCall(url);
        console.log('Stats received:', stats);
        
        // Debug: Check if DOM elements exist
        console.log('Checking DOM elements...');
        const totalSalesEl = document.getElementById('totalSales');
        const grossProfitEl = document.getElementById('grossProfit');
        const totalCogsEl = document.getElementById('totalCogs');
        
        console.log('totalSales element:', totalSalesEl);
        console.log('grossProfit element:', grossProfitEl);
        console.log('totalCogs element:', totalCogsEl);
        
        // Debug: Check formatCurrency function
        console.log('formatCurrency function exists:', typeof formatCurrency);
        
        if (typeof formatCurrency === 'function') {
            console.log('formatCurrency test with 57700:', formatCurrency(57700));
        } else {
            console.error('formatCurrency function not found!');
        }
        
        // Debug: Check the actual values before formatting
        console.log('Raw values:');
        console.log('- total_sales:', stats.total_sales, typeof stats.total_sales);
        console.log('- gross_profit:', stats.gross_profit, typeof stats.gross_profit);
        console.log('- total_cogs:', stats.total_cogs, typeof stats.total_cogs);
        
        // Update stat cards with debugging
        if (totalSalesEl) {
            const formattedSales = formatCurrency(stats.total_sales);
            console.log('Setting totalSales to:', formattedSales);
            totalSalesEl.textContent = formattedSales;
            console.log('totalSales after update:', totalSalesEl.textContent);
        } else {
            console.error('totalSales element not found!');
        }
        
        if (grossProfitEl) {
            const formattedProfit = formatCurrency(stats.gross_profit);
            console.log('Setting grossProfit to:', formattedProfit);
            grossProfitEl.textContent = formattedProfit;
            console.log('grossProfit after update:', grossProfitEl.textContent);
        } else {
            console.error('grossProfit element not found!');
        }
        
        if (totalCogsEl) {
            const formattedCogs = formatCurrency(stats.total_cogs);
            console.log('Setting totalCogs to:', formattedCogs);
            totalCogsEl.textContent = formattedCogs;
            console.log('totalCogs after update:', totalCogsEl.textContent);
        } else {
            console.error('totalCogs element not found!');
        }
        
        document.getElementById('totalOrders').textContent = stats.total_orders;
        document.getElementById('netProfit').textContent = formatCurrency(stats.net_profit);
        document.getElementById('pendingOrders').textContent = stats.pending_orders;
        
        // Update gross margin with color coding
        const grossMarginElement = document.getElementById('grossMargin');
        const marginPercent = stats.gross_margin_percent || 0;
        grossMarginElement.textContent = `${marginPercent.toFixed(1)}%`;
        
        // Apply color coding based on margin percentage
        grossMarginElement.className = 'margin-text';
        if (marginPercent >= 25) {
            grossMarginElement.classList.add('positive');
        } else if (marginPercent >= 15) {
            grossMarginElement.classList.add('neutral');
        } else {
            grossMarginElement.classList.add('negative');
        }
        
        // Update top menu items table
        updateTopMenuTable(stats.top_menu_items || []);
        
        // Update top add-ons table
        updateTopAddOnsTable(stats.top_add_ons || []);
        
        console.log('Stats updated in UI');
        return stats;
    } catch (error) {
        console.error('loadDashboardStats error:', error);
        throw error;
    }
}

// Update top menu items table
function updateTopMenuTable(topItems) {
    const tbody = document.querySelector('#topMenuTable tbody');
    tbody.innerHTML = topItems.map(item => `
        <tr>
            <td>${item.name}</td>
            <td>${item.total_sold}</td>
            <td>${formatCurrency(item.total_revenue)}</td>
        </tr>
    `).join('');
}

// Update top add-ons table
function updateTopAddOnsTable(topAddOns) {
    const tbody = document.querySelector('#topAddOnsTable tbody');
    tbody.innerHTML = topAddOns.map(addon => `
        <tr>
            <td>${addon.name}</td>
            <td>${addon.total_sold}</td>
            <td>${formatCurrency(addon.total_revenue)}</td>
        </tr>
    `).join('');
}

// Load sales chart
async function loadSalesChart(startDate, endDate) {
    try {
        const params = new URLSearchParams();
        
        // Only add date parameters if they are provided
        if (startDate && startDate.trim() !== '') {
            params.append('start_date', startDate);
        }
        if (endDate && endDate.trim() !== '') {
            params.append('end_date', endDate);
        }
        
        const queryString = params.toString();
        const url = queryString ? `/dashboard/stats?${queryString}` : '/dashboard/stats';
        
        const stats = await apiCall(url);
        
        const ctx = document.getElementById('salesChart').getContext('2d');
        
        // Destroy existing chart if it exists
        if (salesChart) {
            salesChart.destroy();
        }
        
        const salesData = stats.sales_chart || [];
        const labels = salesData.map(item => formatDate(item.date));
        const amounts = salesData.map(item => item.amount);
        const orders = salesData.map(item => item.orders);
        
        salesChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Sales Amount',
                        data: amounts,
                        borderColor: '#3498db',
                        backgroundColor: 'rgba(52, 152, 219, 0.1)',
                        tension: 0.4,
                        yAxisID: 'y'
                    },
                    {
                        label: 'Orders Count',
                        data: orders,
                        borderColor: '#e74c3c',
                        backgroundColor: 'rgba(231, 76, 60, 0.1)',
                        tension: 0.4,
                        yAxisID: 'y1'
                    }
                ]
            },
            options: {
                responsive: true,
                interaction: {
                    mode: 'index',
                    intersect: false,
                },
                scales: {
                    x: {
                        display: true,
                        title: {
                            display: true,
                            text: 'Date'
                        }
                    },
                    y: {
                        type: 'linear',
                        display: true,
                        position: 'left',
                        title: {
                            display: true,
                            text: 'Sales Amount'
                        },
                        ticks: {
                            callback: function(value) {
                                return formatCurrency(value);
                            }
                        }
                    },
                    y1: {
                        type: 'linear',
                        display: true,
                        position: 'right',
                        title: {
                            display: true,
                            text: 'Orders Count'
                        },
                        grid: {
                            drawOnChartArea: false,
                        },
                    }
                },
                plugins: {
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                if (context.dataset.label === 'Sales Amount') {
                                    return `Sales: ${formatCurrency(context.parsed.y)}`;
                                } else {
                                    return `Orders: ${context.parsed.y}`;
                                }
                            }
                        }
                    }
                }
            }
        });
    } catch (error) {
        throw error;
    }
}

// Load expense chart
async function loadExpenseChart(startDate, endDate) {
    try {
        const params = new URLSearchParams();
        
        // Only add date parameters if they are provided
        if (startDate && startDate.trim() !== '') {
            params.append('start_date', startDate);
        }
        if (endDate && endDate.trim() !== '') {
            params.append('end_date', endDate);
        }
        
        const queryString = params.toString();
        const url = queryString ? `/dashboard/stats?${queryString}` : '/dashboard/stats';
        
        const stats = await apiCall(url);
        
        const ctx = document.getElementById('expenseChart').getContext('2d');
        
        // Destroy existing chart if it exists
        if (expenseChart) {
            expenseChart.destroy();
        }
        
        const expenseData = stats.expense_chart || [];
        
        // Group expenses by type
        const rawMaterialData = [];
        const operationalData = [];
        const dates = [...new Set(expenseData.map(item => item.date))].sort();
        
        dates.forEach(date => {
            const rawMaterial = expenseData.find(item => item.date === date && item.type === 'raw_material');
            const operational = expenseData.find(item => item.date === date && item.type === 'operational');
            
            rawMaterialData.push(rawMaterial ? rawMaterial.amount : 0);
            operationalData.push(operational ? operational.amount : 0);
        });
        
        expenseChart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: dates.map(date => formatDate(date)),
                datasets: [
                    {
                        label: 'Raw Materials',
                        data: rawMaterialData,
                        backgroundColor: 'rgba(231, 76, 60, 0.7)',
                        borderColor: '#e74c3c',
                        borderWidth: 1
                    },
                    {
                        label: 'Operational',
                        data: operationalData,
                        backgroundColor: 'rgba(52, 152, 219, 0.7)',
                        borderColor: '#3498db',
                        borderWidth: 1
                    }
                ]
            },
            options: {
                responsive: true,
                scales: {
                    x: {
                        stacked: true,
                        title: {
                            display: true,
                            text: 'Date'
                        }
                    },
                    y: {
                        stacked: true,
                        title: {
                            display: true,
                            text: 'Expense Amount'
                        },
                        ticks: {
                            callback: function(value) {
                                return formatCurrency(value);
                            }
                        }
                    }
                },
                plugins: {
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return `${context.dataset.label}: ${formatCurrency(context.parsed.y)}`;
                            }
                        }
                    }
                }
            }
        });
    } catch (error) {
        throw error;
    }
}
