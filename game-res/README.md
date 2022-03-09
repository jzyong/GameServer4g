## 一、目录说明

路径                         |描述
--------------------------- |------------------------------              
docker                      |docker运行，镜像制作脚本

## 二、常见问题

1. Github push `Empty reply from server`


      git fetch origin --prune

2. Github push `OpenSSL SSL_read: Connection was reset, errno 10054`  
 

     git config --global http.sslVerify "false"