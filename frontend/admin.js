let token = '';

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

// 渲染管理员页面的资料列表
async function renderAdminMaterials() {
    const materialsList = document.getElementById('adminMaterialsList');
    const materials = await getMaterials();
    
    materialsList.innerHTML = `
        <h3>现有资料</h3>
        ${materials.map(material => `
            <div class="material-card">
                <h3>${material.name}</h3>
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
    const progressBar = document.getElementById('uploadProgress');
    const progressText = document.getElementById('progressText');
    
    if (!name || !file) {
        alert('请填写资料名称并选择文件！');
        return;
    }
    
    const formData = new FormData();
    formData.append('name', name);
    formData.append('file', file);

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