
```
nohup  ./cpsql-webdav -server-name 127.0.0.1 -project-name 本地项目 -webdav-username wangmoxiang  -webdav-password 594Wmx1991 -sql-database test  -sql-user root  -sql-password 123456 > output.log 2>&1 &
```

- server-name 服务名称 可以使用服务器IP
- project-name 项目名称 区别同一个服务中不同项目
- webdav-username webdav 用户名
- webdav-password  webdav 密码
- sql-database 数据库
- sql-user 数据库用户名
- sql-password 数据库密码


更多配置请看 config.yaml
