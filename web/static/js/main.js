// 域名估价系统前端脚本

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 域名输入框自动聚焦
    const domainInput = document.getElementById('domain');
    if (domainInput) {
        domainInput.focus();
    }

    // 表单验证
    const form = document.querySelector('form[action="/estimate"]');
    if (form) {
        form.addEventListener('submit', function(event) {
            const domain = domainInput.value.trim();
            
            // 简单的域名格式验证
            if (!isValidDomain(domain)) {
                event.preventDefault();
                alert('请输入有效的域名格式，例如：example.com');
                domainInput.focus();
            }
        });
    }
});

// 验证域名格式的简单函数
function isValidDomain(domain) {
    // 基本的域名格式验证
    const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/;
    return domainRegex.test(domain);
}
