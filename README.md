# jtyl_bitable

## 1.将服务打包成可以在Linux系统的二进制文件
```powershell
$env:GOOS='linux'; $env:GOARCH='amd64'; $env:CGO_ENABLED='0'; go build -ldflags="-s -w" -o myService
```

## 2.Linux常用命令
```bash
#切换到root用户
sudo su -
#查看所有运行中的服务
systemctl list-units --type=service --state=running
#停止运行服务
systemctl stop jtyl_bitable
#删除文件
rm /var/www/jtyl_bitable/myService
```

## 3.上传服务文件到云服务器
```bash
scp .\myService root@118.31.238.252:/var/www/jtyl_bitable/
#password: Ye230516@
```

## 4.启动服务
```bash
sudo systemctl daemon-reload
sudo systemctl enable jtyl_bitable
sudo systemctl start jtyl_bitable
sudo systemctl status jtyl_bitable
```