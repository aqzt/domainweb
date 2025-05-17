# 域名估价系统安装指南

本文档提供了域名估价系统的详细安装和配置说明。

## 系统要求

- Go 1.16 或更高版本
- MySQL 5.7 或更高版本
- 支持现代浏览器（Chrome、Firefox、Safari、Edge等）

## 安装步骤

### 1. 获取源代码

```bash
git clone https://github.com/yourusername/domainweb.git
cd domainweb
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

#### 3.1 创建数据库和表

使用提供的SQL脚本初始化数据库：

```bash
mysql -u root -p < scripts/init_db.sql
```

#### 3.2 配置数据库连接

编辑`config/config.json`文件，修改数据库连接信息：

```json
{
  "database": {
    "driver": "mysql",
    "host": "127.0.0.1",
    "port": 3306,
    "username": "your_username",
    "password": "your_password",
    "dbname": "domainweb",
    "charset": "utf8mb4",
    "parseTime": true,
    "maxOpenConns": 10,
    "maxIdleConns": 5,
    "connMaxLifetime": 3600
  }
}
```

### 4. 构建和运行

```bash
go build -o domainweb main.go
./domainweb
```

或者直接运行：

```bash
go run main.go
```

### 5. 访问系统

打开浏览器，访问：http://localhost:8080

## 配置说明

### 服务器配置

编辑`config/config.json`文件中的`server`部分：

```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0",
    "readTimeout": 10,
    "writeTimeout": 10,
    "maxHeaderBytes": 1048576
  }
}
```

| 参数 | 描述 | 默认值 |
|------|------|--------|
| port | 服务器监听端口 | 8080 |
| host | 服务器监听地址 | 0.0.0.0 |
| readTimeout | 读取超时时间（秒） | 10 |
| writeTimeout | 写入超时时间（秒） | 10 |
| maxHeaderBytes | 最大请求头大小 | 1048576 |

### 估价配置

编辑`config/config.json`文件中的`estimation`部分：

```json
{
  "estimation": {
    "basePrice": 25.0,
    "baseGrade": -0.5,
    "defaultHistoryLimit": 50
  }
}
```

| 参数 | 描述 | 默认值 |
|------|------|--------|
| basePrice | 估价基数（元） | 25.0 |
| baseGrade | 等级基数 | -0.5 |
| defaultHistoryLimit | 默认历史记录限制 | 50 |

## 动态属性API配置（可选）

要使用真实的动态属性数据，需要配置相应的API密钥。编辑`config/config.json`文件，添加以下部分：

```json
{
  "apis": {
    "alexa": {
      "apiKey": "your_alexa_api_key",
      "endpoint": "https://awis.api.alexa.com/api"
    },
    "whois": {
      "apiKey": "your_whois_api_key",
      "endpoint": "https://whoisapi.whoisxmlapi.com/api/v1"
    },
    "search": {
      "apiKey": "your_search_api_key",
      "endpoint": "https://api.example.com/search"
    }
  }
}
```

## 生产环境部署

### 使用Systemd服务

创建systemd服务文件：

```bash
sudo nano /etc/systemd/system/domainweb.service
```

添加以下内容：

```
[Unit]
Description=Domain Valuation System
After=network.target mysql.service

[Service]
User=www-data
WorkingDirectory=/path/to/domainweb
ExecStart=/path/to/domainweb/domainweb
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
sudo systemctl enable domainweb
sudo systemctl start domainweb
```

### 使用Nginx反向代理

安装Nginx：

```bash
sudo apt install nginx
```

创建Nginx配置文件：

```bash
sudo nano /etc/nginx/sites-available/domainweb
```

添加以下内容：

```
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用站点并重启Nginx：

```bash
sudo ln -s /etc/nginx/sites-available/domainweb /etc/nginx/sites-enabled/
sudo systemctl restart nginx
```

## 故障排除

### 数据库连接问题

如果遇到数据库连接问题，请检查：

1. MySQL服务是否运行
2. 数据库用户名和密码是否正确
3. 数据库名称是否正确
4. MySQL用户是否有足够的权限

### 服务无法启动

如果服务无法启动，请检查：

1. 日志文件中的错误信息
2. 端口是否被占用
3. 配置文件是否正确

### 获取帮助

如有其他问题，请提交GitHub Issue或联系项目维护者。
