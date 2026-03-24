const API_BASE = '/api/v1';
const TOKEN_KEY = 'shorturl_token';

// 页面加载时检查登录状态
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    loadShortURLs();
    document.getElementById('usernameDisplay').textContent = localStorage.getItem('username') || 'Admin';
    createNotificationContainer();
});

// 创建通知容器
function createNotificationContainer() {
    const container = document.createElement('div');
    container.id = 'notificationContainer';
    container.className = 'fixed top-4 right-4 z-50 space-y-2';
    document.body.appendChild(container);
}

// 显示通知
function showNotification(message, type = 'info') {
    const container = document.getElementById('notificationContainer');
    const notification = document.createElement('div');
    
    const colors = {
        success: 'bg-green-500',
        error: 'bg-red-500',
        info: 'bg-blue-500'
    };
    
    const icons = {
        success: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>',
        error: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>',
        info: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>'
    };
    
    notification.className = `${colors[type]} text-white px-4 py-3 rounded-lg shadow-lg flex items-center space-x-3 transform transition-all duration-300 translate-x-full opacity-0`;
    notification.innerHTML = `
        ${icons[type]}
        <span>${message}</span>
    `;
    
    container.appendChild(notification);
    
    // 动画显示
    setTimeout(() => {
        notification.classList.remove('translate-x-full', 'opacity-0');
    }, 10);
    
    // 3 秒后自动消失
    setTimeout(() => {
        notification.classList.add('translate-x-full', 'opacity-0');
        setTimeout(() => {
            container.removeChild(notification);
        }, 300);
    }, 3000);
}

// 检查认证状态
function checkAuth() {
    const token = localStorage.getItem(TOKEN_KEY);
    const firstLogin = localStorage.getItem('first_login');
    
    if (!token) {
        window.location.href = '/login';
        return;
    }
    
    if (firstLogin === 'true') {
        window.location.href = '/change-password';
        return;
    }
}

// 退出登录
function logout() {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem('first_login');
    localStorage.removeItem('username');
    window.location.href = '/login';
}

// 加载短链列表
async function loadShortURLs() {
    try {
        const token = localStorage.getItem(TOKEN_KEY);
        const response = await fetch(`${API_BASE}/shorturls`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            renderShortURLs(data.data.list);
            updateStats(data.data.list, data.data.total);
        } else {
            throw new Error(data.message);
        }
    } catch (error) {
        console.error('加载短链列表失败:', error);
        showNotification('加载失败：' + error.message, 'error');
    }
}

// 渲染短链列表
function renderShortURLs(shortUrls) {
    const tbody = document.getElementById('shortUrlTableBody');
    
    if (shortUrls.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="7" class="px-6 py-12 text-center">
                    <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"></path>
                    </svg>
                    <p class="mt-2 text-sm text-gray-500">暂无数据</p>
                </td>
            </tr>
        `;
        return;
    }
    
    tbody.innerHTML = shortUrls.map(url => `
        <tr class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${url.id}</td>
            <td class="px-6 py-4 whitespace-nowrap">
                <a href="${url.short_url}" target="_blank" class="text-indigo-600 hover:text-indigo-900 font-medium">
                    ${url.short_code}
                </a>
            </td>
            <td class="px-6 py-4 text-sm text-gray-500 max-w-xs truncate" title="${url.original_url}">
                ${url.original_url}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">${url.clicks}</td>
            <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${url.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                    ${url.is_active ? '激活' : '禁用'}
                </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                ${new Date(url.created_at).toLocaleString('zh-CN')}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                <button onclick="copyShortUrl('${url.short_url}')" class="text-indigo-600 hover:text-indigo-900" title="复制链接">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                    </svg>
                </button>
                <button onclick="editShortUrl(${url.id})" class="text-blue-600 hover:text-blue-900" title="编辑">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                    </svg>
                </button>
                <button onclick="toggleShortUrlStatus(${url.id}, ${url.is_active})" class="${url.is_active ? 'text-yellow-600 hover:text-yellow-900' : 'text-green-600 hover:text-green-900'}" title="${url.is_active ? '禁用' : '启用'}">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${url.is_active ? 'M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636' : 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'}"></path>
                    </svg>
                </button>
                <button onclick="deleteShortUrl(${url.id})" class="text-red-600 hover:text-red-900" title="删除">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                </button>
            </td>
        </tr>
    `).join('');
}

// 更新统计信息
function updateStats(shortUrls, total) {
    document.getElementById('totalCount').textContent = total;
    document.getElementById('activeCount').textContent = shortUrls.filter(url => url.is_active).length;
    document.getElementById('totalClicks').textContent = shortUrls.reduce((sum, url) => sum + url.clicks, 0);
}

// 显示创建模态框
function showCreateModal() {
    document.getElementById('modalTitle').textContent = '创建短链';
    document.getElementById('shortUrlForm').reset();
    document.getElementById('editId').value = '';
    document.getElementById('modalErrorMessage').classList.add('hidden');
    document.getElementById('modal').classList.remove('hidden');
}

// 关闭模态框
function closeModal() {
    document.getElementById('modal').classList.add('hidden');
}

// 编辑短链
async function editShortUrl(id) {
    try {
        const token = localStorage.getItem(TOKEN_KEY);
        const response = await fetch(`${API_BASE}/shorturls/${id}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            const url = data.data;
            document.getElementById('modalTitle').textContent = '编辑短链';
            document.getElementById('editId').value = url.id;
            document.getElementById('originalUrl').value = url.original_url;
            document.getElementById('customCode').value = url.short_code;
            document.getElementById('modalErrorMessage').classList.add('hidden');
            document.getElementById('modal').classList.remove('hidden');
        } else {
            throw new Error(data.message);
        }
    } catch (error) {
        console.error('获取短链详情失败:', error);
        showNotification('获取失败：' + error.message, 'error');
    }
}

// 复制短链
function copyShortUrl(shortUrl) {
    // 使用兼容的方式复制文本
    const textArea = document.createElement('textarea');
    textArea.value = shortUrl;
    textArea.style.position = 'fixed';
    textArea.style.left = '-999999px';
    textArea.style.top = '-999999px';
    document.body.appendChild(textArea);
    textArea.select();
    
    try {
        const successful = document.execCommand('copy');
        if (successful) {
            showNotification('短链已复制到剪贴板', 'success');
        } else {
            showNotification('复制失败', 'error');
        }
    } catch (err) {
        console.error('复制失败:', err);
        // 尝试使用现代 API
        navigator.clipboard.writeText(shortUrl).then(() => {
            showNotification('短链已复制到剪贴板', 'success');
        }).catch(() => {
            showNotification('复制失败', 'error');
        });
    }
    
    document.body.removeChild(textArea);
}

// 切换短链状态
async function toggleShortUrlStatus(id, isActive) {
    if (!confirm(`确定要${isActive ? '禁用' : '启用'}该短链吗？`)) {
        return;
    }
    
    try {
        const token = localStorage.getItem(TOKEN_KEY);
        const response = await fetch(`${API_BASE}/shorturls/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ is_active: !isActive })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            loadShortURLs();
        } else {
            throw new Error(data.message);
        }
    } catch (error) {
        console.error('更新状态失败:', error);
        showNotification('更新失败：' + error.message, 'error');
    }
}

// 删除短链
async function deleteShortUrl(id) {
    if (!confirm('确定要删除该短链吗？此操作不可恢复。')) {
        return;
    }
    
    try {
        const token = localStorage.getItem(TOKEN_KEY);
        const response = await fetch(`${API_BASE}/shorturls/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            loadShortURLs();
        } else {
            throw new Error(data.message);
        }
    } catch (error) {
        console.error('删除失败:', error);
        showNotification('删除失败：' + error.message, 'error');
    }
}

// 表单提交
document.getElementById('shortUrlForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const editId = document.getElementById('editId').value;
    const originalUrl = document.getElementById('originalUrl').value;
    const customCode = document.getElementById('customCode').value;
    const modalSubmitBtn = document.getElementById('modalSubmitBtn');
    const modalErrorMessage = document.getElementById('modalErrorMessage');
    const modalErrorText = document.getElementById('modalErrorText');
    
    modalSubmitBtn.disabled = true;
    modalErrorMessage.classList.add('hidden');
    
    try {
        const token = localStorage.getItem(TOKEN_KEY);
        const url = editId 
            ? `${API_BASE}/shorturls/${editId}` 
            : `${API_BASE}/shorturls`;
        
        const method = editId ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                original_url: originalUrl,
                custom_code: customCode
            })
        });
        
        const data = await response.json();
        
        if (data.code === 200) {
            closeModal();
            loadShortURLs();
        } else {
            throw new Error(data.message);
        }
    } catch (error) {
        modalErrorText.textContent = error.message;
        modalErrorMessage.classList.remove('hidden');
        showNotification('操作失败：' + error.message, 'error');
    } finally {
        modalSubmitBtn.disabled = false;
    }
});
