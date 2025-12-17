#jtyl_bitable

1.将服务打包成可以在Linux系统的二进制文件
```bash
#在windows命令提示符（CMD）中
set GOOS=linux
set GOARCH=amd64
go build -o myservice main.go
```


#Linux常用命令
```bash
#切换管理员权限
sudo su

#查看nginx服务状态
systemctl status nginx
#systemctl: 系统服务管理命令。
#status: 查看服务运行状态的子命令。
#nginx: 要查看的服务名称。

```
本地文件上传到云服务器
```bash
#在windows命令提示符（CMD）中
scp .\jtyl_bitable root@118.31.238.252:/var/www/
```

```bash
# 1. 创建systemd服务文件
sudo nano /etc/systemd/system/jtyl_bitable.service

# 2. 粘贴上面的配置文件（注意修改路径）
# 3. 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable jtyl_bitable
sudo systemctl start jtyl_bitable

# 4. 检查状态
sudo systemctl status jtyl_bitable
```

systemd配置文件
/etc/systemd/system/jtyl_bitable.service
```ini
[Unit]
Description=JTYL Bitable Go Service
After=network.target

[Service]

Type=simple
User=root
WorkingDirectory=/var/www/jtyl_bitable
ExecStart=/var/www/jtyl_bitable/jtyl_bitable
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

读取日志文件app.log
```bash
tail app.log
```

读取nginx配置文件
```bash
sudo nano /etc/nginx/nginx.conf
```


