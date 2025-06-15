// 从后端API获取资料列表，支持按文件夹筛选
async function getMaterials(folderId = null) {
    try {
        const url = folderId ? `/api/materials?folder_id=${folderId}` : '/api/materials';
        console.log('准备请求:', url);
        const response = await fetch(url);
        if (!response.ok) {
            console.error('获取资料列表失败，状态码:', response.status);
            return [];
        }
        const data = await response.json();
        console.log('获取到资料列表:', data);
        return data;
    } catch (e) {
        console.error('请求资料列表出错:', e);
        return [];
    }
}

// 从后端API获取文件夹树
async function getFolders() {
    try {
        console.log('准备请求 /api/folders');
        const response = await fetch('/api/folders');
        if (!response.ok) {
            console.error('获取文件夹列表失败，状态码:', response.status);
            return [];
        }
        const data = await response.json();
        console.log('获取到文件夹列表:', data);
        return data;
    } catch (e) {
        console.error('请求 /api/folders 出错:', e);
        return [];
    }
}

// 全局变量存储当前选中的文件夹
let currentFolderId = null;

// 渲染文件夹树
async function renderFolderTree() {
    const folderTree = document.getElementById('folderTree');
    if (!folderTree) {
        console.error('未找到 #folderTree 元素');
        return;
    }

    const folders = await getFolders();
    if (!Array.isArray(folders)) {
        console.error('后端返回的文件夹不是数组:', folders);
        return;
    }

    // 添加"所有资料"选项
    let treeHTML = `
        <div class="folder-item ${currentFolderId === null ? 'active' : ''}" data-folder-id="">
            <span class="folder-icon">📁</span>
            <span class="folder-name">所有资料</span>
        </div>
    `;

    // 递归渲染文件夹树
    treeHTML += renderFolderNodes(folders, 0);

    folderTree.innerHTML = treeHTML;

    // 添加点击事件监听
    folderTree.addEventListener('click', handleFolderClick);

    console.log('文件夹树渲染完成');
}

// 递归渲染文件夹节点
function renderFolderNodes(folders, level) {
    return folders.map(folder => {
        const indent = level * 20; // 每层缩进20px
        const hasChildren = folder.children && folder.children.length > 0;
        const isExpanded = false; // 默认折叠

        let html = `
            <div class="folder-item ${currentFolderId === folder.id ? 'active' : ''}"
                 data-folder-id="${folder.id}"
                 data-level="${level}"
                 style="padding-left: ${indent}px;">
                <span class="folder-toggle ${hasChildren ? 'has-children' : ''} ${isExpanded ? 'expanded' : ''}"
                      data-folder-id="${folder.id}">
                    ${hasChildren ? (isExpanded ? '▼' : '▶') : ''}
                </span>
                <span class="folder-icon">📁</span>
                <span class="folder-name">${folder.name}</span>
            </div>
        `;

        // 如果有子文件夹且展开，递归渲染子文件夹
        if (hasChildren && isExpanded) {
            html += `<div class="folder-children" data-parent-id="${folder.id}">`;
            html += renderFolderNodes(folder.children, level + 1);
            html += '</div>';
        } else if (hasChildren) {
            // 如果有子文件夹但未展开，创建隐藏的子文件夹容器
            html += `<div class="folder-children hidden" data-parent-id="${folder.id}">`;
            html += renderFolderNodes(folder.children, level + 1);
            html += '</div>';
        }

        return html;
    }).join('');
}

// 渲染资料列表
async function renderMaterials(folderId = null) {
    const materialsList = document.getElementById('materialsList');
    const currentFolderName = document.getElementById('currentFolderName');

    if (!materialsList) {
        console.error('未找到 #materialsList 元素');
        return;
    }

    const materials = await getMaterials(folderId);
    if (!Array.isArray(materials)) {
        console.error('后端返回的资料不是数组:', materials);
        return;
    }

    // 更新当前文件夹名称显示
    if (currentFolderName) {
        if (folderId === null) {
            currentFolderName.textContent = '所有资料';
        } else {
            // 从文件夹树中找到对应的文件夹名称
            const folderItem = document.querySelector(`[data-folder-id="${folderId}"] .folder-name`);
            currentFolderName.textContent = folderItem ? folderItem.textContent : '未知文件夹';
        }
    }

    materialsList.innerHTML = materials.map(material => `
        <div class="material-card">
            <h3>${material.name}</h3>
            <p class="material-info">
                ${material.folder ? `文件夹: ${material.folder.name}` : ''}
            </p>
            <a href="/api/materials/${material.id}/download" download="${material.name}">下载</a>
        </div>
    `).join('');

    console.log('资料渲染完成，文件夹ID:', folderId);
}

// 处理文件夹点击事件
function handleFolderClick(event) {
    const target = event.target;

    // 处理文件夹展开/折叠
    if (target.classList.contains('folder-toggle') && target.classList.contains('has-children')) {
        event.stopPropagation();
        const folderId = target.dataset.folderId;
        toggleFolder(folderId);
        return;
    }

    // 处理文件夹选择
    const folderItem = target.closest('.folder-item');
    if (folderItem) {
        const folderId = folderItem.dataset.folderId;
        selectFolder(folderId);
    }
}

// 切换文件夹展开/折叠状态
function toggleFolder(folderId) {
    const toggle = document.querySelector(`[data-folder-id="${folderId}"].folder-toggle`);
    const children = document.querySelector(`[data-parent-id="${folderId}"].folder-children`);

    if (!toggle || !children) return;

    const isExpanded = toggle.classList.contains('expanded');

    if (isExpanded) {
        // 折叠
        toggle.classList.remove('expanded');
        toggle.textContent = '▶';
        children.classList.add('hidden');
    } else {
        // 展开
        toggle.classList.add('expanded');
        toggle.textContent = '▼';
        children.classList.remove('hidden');
    }
}

// 选择文件夹
function selectFolder(folderId) {
    // 移除之前的选中状态
    document.querySelectorAll('.folder-item.active').forEach(item => {
        item.classList.remove('active');
    });

    // 添加新的选中状态
    const selectedItem = document.querySelector(`[data-folder-id="${folderId || ''}"]`);
    if (selectedItem) {
        selectedItem.classList.add('active');
    }

    // 更新全局变量
    currentFolderId = folderId || null;

    // 重新渲染资料列表
    renderMaterials(currentFolderId);
}

// 初始化页面
async function initializePage() {
    await renderFolderTree();
    await renderMaterials();
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', initializePage);