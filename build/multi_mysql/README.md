# 1. 构建主从镜像
docker build -f ./master/Dockerfile -t mysql_master:5.7.17  
docker build -f ./slave/Dockerfile -t mysql_slave:5.7.17  

# 2. 启动主从
docker run --name mysql_master -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql_master:5.7.17  
docker run --name mysql_slave_1 -p 3307:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql_slave:5.7.17  
docker run --name mysql_slave_2 -p 3308:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql_slave:5.7.17  

# 3. 配置主从关系

## 获取主从ip

master
$ docker inspect --format="{{.NetworkSettings.IPAddress}}" containerID
$ 172.17.0.2

slave
$ docker inspect --format="{{.NetworkSettings.IPAddress}}" containerID
$ 172.17.0.3

slave
$ docker inspect --format="{{.NetworkSettings.IPAddress}}" containerID
$ 172.17.0.4

## 配置master  
创建数据同步用户   

CREATE USER 'slave'@'%' IDENTIFIED BY '123456';  
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'slave'@'%';     
## 查询master状态
```
mysql> SHOW MASTER STATUS\G
*************************** 1. row ***************************
File: replicas-mysql-bin.000003
Position: 1307
Binlog_Do_DB: test
Binlog_Ignore_DB: manual, mysql
Executed_Gtid_Set: 3E11FA47-71CA-11E1-9E33-C80AA9429562:1-5
1 row in set (0.00 sec)
```

## 配置slave 

连接slave, 执行:
```
change master to master_host='172.17.0.2', master_user='slave', master_password='123456', master_port=3306, master_log_file='replicas-mysql-bin.000003', master_log_pos=1307, master_connect_retry=30;  
```
master_log_file与master状态中的FIle保持一致  
master_log_pos与master状态中的Position保持一致  
参数解释:
```
master_host: Master 的IP地址
master_user: 在 Master 中授权的用于数据同步的用户
master_password: 同步数据的用户的密码
master_port: Master 的数据库的端口号
master_log_file: 指定 Slave 从哪个日志文件开始复制数据，即上文中提到的 File 字段的值
master_log_pos: 从哪个 Position 开始读，即上文中提到的 Position 字段的值
master_connect_retry: 当重新建立主从连接时，如果连接失败，重试的时间间隔，单位是秒，默认是60秒。
```
开始同步:
```
start slave;
```

## 查询slave状态
 
```
root@localhost (none)>show slave status\G
*************************** 1. row ***************************
               Slave_IO_State: Waiting for master to send event
                  Master_Host: 192.168.1.100
                  Master_User: mysync
                  Master_Port: 3306
                Connect_Retry: 60
              Master_Log_File: mysql-bin.001822
          Read_Master_Log_Pos: 290072815
               Relay_Log_File: mysqld-relay-bin.005201
                Relay_Log_Pos: 256529594
        Relay_Master_Log_File: mysql-bin.001821
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
              Replicate_Do_DB: 
          Replicate_Ignore_DB: 
           Replicate_Do_Table: 
       Replicate_Ignore_Table: 
      Replicate_Wild_Do_Table: 
  Replicate_Wild_Ignore_Table: 
                   Last_Errno: 0
                   Last_Error: 
                 Skip_Counter: 0
          Exec_Master_Log_Pos: 256529431
              Relay_Log_Space: 709504534
              Until_Condition: None
               Until_Log_File: 
                Until_Log_Pos: 0
           Master_SSL_Allowed: No
           Master_SSL_CA_File: 
           Master_SSL_CA_Path: 
              Master_SSL_Cert: 
            Master_SSL_Cipher: 
               Master_SSL_Key: 
        Seconds_Behind_Master: 2923
Master_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error: 
               Last_SQL_Errno: 0
               Last_SQL_Error: 
  Replicate_Ignore_Server_Ids: 
             Master_Server_Id: 1
                  Master_UUID: 13ee75bb-99e2-11e6-be4d-b499baa80e6e
             Master_Info_File: /home/data/mysql/master.info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Slave_SQL_Running_State: Reading event from the relay log
           Master_Retry_Count: 86400
                  Master_Bind: 
      Last_IO_Error_Timestamp: 
     Last_SQL_Error_Timestamp: 
               Master_SSL_Crl: 
           Master_SSL_Crlpath: 
           Retrieved_Gtid_Set: 
            Executed_Gtid_Set: 
                Auto_Position: 0
1 row in set (0.02 sec)

```