GO+MySQL 实现的短链接服务
1. 通过队列解决MySQL的瓶颈问题
2. 存储模型是 短链名 - ID - 原URL
3. ID唯一，作为MySQL主键保证自增；
4. 短链名通过ID得到
