let token = '';

// 获取文件夹列表
async function getFolders() {
    try {
        const response = await fetch('http://14.103.237.30:8081/api/folders', {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        if (response.ok) {
            return await response.json();
        } else {
            console.error('获取文件夹列表失败:', response.status);
            return [];
        }
    } catch (error) {
        console.error('获取文件夹列表失败:', error);
        return [];
    }
}

// 创建文件夹
async function createFolder() {
    const name = document.getElementById('folderName').value.trim();
    const parentId = document.getElementById('parentFolderSelect').value || null;

    if (!name) {
        alert('请输入文件夹名称！');
        return;
    }

    try {
        const response = await fetch('http://14.103.237.30:8081/api/folders', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({
                name: name,
                parent_id: parentId ? parseInt(parentId) : null,
                description: ''
            }),
        });

        const data = await response.json();
        if (response.ok) {
            document.getElementById('folderName').value = '';
            await renderAdminFolders();
            alert('文件夹创建成功！');
        } else {
            alert(data.error || '创建文件夹失败');
        }
    } catch (error) {
        console.error('创建文件夹失败:', error);
        alert('创建文件夹失败，请重试');
    }
}

// 重命名文件夹
async function renameFolder(folderId) {
    const newName = prompt('请输入新的文件夹名称:');
    if (!newName || !newName.trim()) {
        return;
    }

    try {
        const response = await fetch(`http://14.103.237.30:8081/api/folders/${folderId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({
                name: newName.trim(),
                description: ''
            }),
        });

        const data = await response.json();
        if (response.ok) {
            await renderAdminFolders();
            alert('文件夹重命名成功！');
        } else {
            alert(data.error || '重命名失败');
        }
    } catch (error) {
        console.error('重命名失败:', error);
        alert('重命名失败，请重试');
    }
}

// 删除文件夹
async function deleteFolder(folderId) {
    if (!confirm('确定要删除这个文件夹吗？删除前请确保文件夹为空。')) {
        return;
    }

    try {
        const response = await fetch(`http://14.103.237.30:8081/api/folders/${folderId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        const data = await response.json();
        if (response.ok) {
            await renderAdminFolders();
            alert('文件夹删除成功！');
        } else {
            alert(data.error || '删除失败');
        }
    } catch (error) {
        console.error('删除失败:', error);
        alert('删除失败，请重试');
    }
}

// 登录
async function login() {
    const password = document.getElementById('password').value;
    console.log('尝试登录，密码长度:', password.length);
    
    try {
        const response = await fetch('http://14.103.237.30:8081/api/admin/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ password }),
        });

        console.log('登录响应状态:', response.status);
        const data = await response.json();
        console.log('登录响应数据:', data);
        
        if (response.ok) {
            token = data.token;
            document.getElementById('loginForm').style.display = 'none';
            document.getElementById('adminPanel').style.display = 'block';
            await renderAdminFolders();
            renderAdminMaterials();
        } else {
            alert(data.error || '登录失败');
        }
    } catch (error) {
        console.error('登录失败:', error);
        alert('登录失败，请重试');
    }
}

// 获取资料列表
async function getMaterials() {
    try {
        const response = await fetch('http://14.103.237.30:8080/api/materials');
        const materials = await response.json();
        return materials;
    } catch (error) {
        console.error('获取资料列表失败:', error);
        return [];
    }
}

// 渲染管理员文件夹树
async function renderAdminFolders() {
    const folders = await getFolders();

    // 更新父文件夹选择器
    const parentSelect = document.getElementById('parentFolderSelect');
    const targetSelect = document.getElementById('targetFolderSelect');

    const folderOptions = buildFolderOptions(folders);
    parentSelect.innerHTML = '<option value="">根目录</option>' + folderOptions;
    targetSelect.innerHTML = '<option value="">根目录</option>' + folderOptions;

    // 渲染文件夹树
    const folderTree = document.getElementById('adminFolderTree');
    folderTree.innerHTML = renderAdminFolderNodes(folders, 0);
}

// 构建文件夹选项
function buildFolderOptions(folders, level = 0) {
    let options = '';
    const indent = '　'.repeat(level); // 使用全角空格缩进

    folders.forEach(folder => {
        options += `<option value="${folder.id}">${indent}${folder.name}</option>`;
        if (folder.children && folder.children.length > 0) {
            options += buildFolderOptions(folder.children, level + 1);
        }
    });

    return options;
}

// 渲染管理员文件夹节点
function renderAdminFolderNodes(folders, level) {
    return folders.map(folder => {
        const indent = level * 20;
        const hasChildren = folder.children && folder.children.length > 0;

        let html = `
            <div class="admin-folder-item" style="padding-left: ${indent}px;">
                <span class="folder-icon">📁</span>
                <span class="folder-name">${folder.name}</span>
                <div class="folder-actions">
                    <button onclick="renameFolder(${folder.id})" class="btn-small">重命名</button>
                    <button onclick="deleteFolder(${folder.id})" class="btn-small btn-danger">删除</button>
                </div>
            </div>
        `;

        if (hasChildren) {
            html += renderAdminFolderNodes(folder.children, level + 1);
        }

        return html;
    }).join('');
}

// 渲染管理员页面的资料列表
async function renderAdminMaterials() {
    const materialsList = document.getElementById('adminMaterialsList');
    const materials = await getMaterials();

    materialsList.innerHTML = `
        <h3>现有资料</h3>
        ${materials.map(material => `
            <div class="material-card">
                <h3>${material.name}</h3>
                <p class="material-info">
                    ${material.folder ? `文件夹: ${material.folder.name}` : '根目录'}
                </p>
                <button onclick="deleteMaterial(${material.id})">删除</button>
            </div>
        `).join('')}
    `;
}

// 删除资料
async function deleteMaterial(id) {
    if (!confirm('确定要删除这个资料吗？')) {
        return;
    }

    try {
        const response = await fetch(`http://14.103.237.30:8081/api/materials/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            renderAdminMaterials();
        } else {
            const data = await response.json();
            alert(data.error || '删除失败');
        }
    } catch (error) {
        console.error('删除失败:', error);
        alert('删除失败，请重试');
    }
}

// 处理文件上传
document.getElementById('uploadForm').addEventListener('submit', function(e) {
    e.preventDefault();

    const name = document.getElementById('materialName').value;
    const file = document.getElementById('materialFile').files[0];
    const folderId = document.getElementById('targetFolderSelect').value;
    const progressBar = document.getElementById('uploadProgress');
    const progressText = document.getElementById('progressText');

    if (!name || !file) {
        alert('请填写资料名称并选择文件！');
        return;
    }

    const formData = new FormData();
    formData.append('name', name);
    formData.append('file', file);
    if (folderId) {
        formData.append('folder_id', folderId);
    }

    const xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://14.103.237.30:8081/api/materials', true);
    xhr.setRequestHeader('Authorization', `Bearer ${token}`);

    progressBar.style.display = 'block';
    progressText.style.display = 'inline';
    progressBar.value = 0;
    progressText.textContent = '0%';

    xhr.upload.onprogress = function(event) {
        if (event.lengthComputable) {
            const percent = Math.round((event.loaded / event.total) * 100);
            progressBar.value = percent;
            progressText.textContent = percent + '%';
        }
    };

    xhr.onload = function() {
        progressBar.style.display = 'none';
        progressText.style.display = 'none';
        if (xhr.status >= 200 && xhr.status < 300) {
            document.getElementById('materialName').value = '';
            document.getElementById('materialFile').value = '';
            renderAdminMaterials();
            alert('资料上传成功！');
        } else {
            let data = {};
            try { data = JSON.parse(xhr.responseText); } catch {}
            alert(data.error || '上传失败');
        }
    };

    xhr.onerror = function() {
        progressBar.style.display = 'none';
        progressText.style.display = 'none';
        alert('上传失败，请重试');
    };

    xhr.send(formData);
}); 