# simple-mysql-mgp-hook
简单的检测指定节点是否在集群内
## 使用说明
在无指定配置文件下,程序会默认读当前执行目录下的config.yml

参考项目内config.yml文件

  ``` yaml 
    node:
       name: test-mysql01 ## 在集群内节点的名称
     mysql:
       user: ## 数据库账号
       password: ## 密码
       host: ## host
       port: ## 端口
       name: performance_schema
     logpath:
       path: ##日志路径
     heartbeat:
       interval: 5 ## 心跳检测 单位为妙
     ## 期望校验值
       check-value: ONLINE
     ## sql 语句
       sql: select field1 from check_table
     ## 查询条件 (字段名:值,可以为多个 需要符合yaml格式规范)
       query-key-value:
        id : 1
       downcommand: ## down的时候执行的相关命令
         - balabla
         - balabla
       upcommand: ## 恢复时候的命令
        - balbal
        - balala
       
  ```

需要创建一个用户并且给予权限访问 **performance_schema** 库的 **select** 权限

节点内名字必须正确,否则程序启动失败
   
##hearbeat解析
hearbeat内**command**为列表配置,当检测到不在集群内会顺序执行配置的命令

hearbeat内**sql** 为查询基本语句,**query-key-value**为查询条件(可以为多个)
例如:
```yaml
    sql: select field1 from check_table
    query-key-value:
      id : 1
```
程序会执行的check语句为
```sql
select field1 from check_table where id = 1
```
比对**field1** 是否等于check-value预定值,如果不为期望值则执行**downcommand**内所有命令

##logpath
logpath内的**path** 默认不指定的情况下会在执行文件的当前目录下生成hook.log文件

## 相关命令

-path your config.yml path

```shell
./simple-mysql-mgp-hook -path /your-path/config.yml
```

