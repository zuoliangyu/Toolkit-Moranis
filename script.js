// 从后端API获取资料列表
async function getMaterials() {
    try {
        console.log('准备请求 /api/materials');
        const response = await fetch('/api/materials');
        if (!response.ok) {
            console.error('获取资料列表失败，状态码:', response.status);
            return [];
        }
        const data = await response.json();
        console.log('获取到资料列表:', data);
        return data;
    } catch (e) {
        console.error('请求 /api/materials 出错:', e);
        return [];
    }
}

// 渲染资料列表
async function renderMaterials() {
    const materialsList = document.getElementById('materialsList');
    if (!materialsList) {
        console.error('未找到 #materialsList 元素');
        return;
    }
    const materials = await getMaterials();
    if (!Array.isArray(materials)) {
        console.error('后端返回的资料不是数组:', materials);
        return;
    }
    materialsList.innerHTML = materials.map(material => `
        <div class="material-card">
            <h3>${material.name}</h3>
            <a href="/api/materials/${material.id}/download" download="${material.name}">下载</a>
        </div>
    `).join('');
    console.log('资料渲染完成');
}

// 页面加载时渲染资料列表
document.addEventListener('DOMContentLoaded', renderMaterials); 