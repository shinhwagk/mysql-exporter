```sh
mysql 

mysqlslap -hxxxx -P3306 -uzhangxu -p123456 --concurrency=100 --iterations=1 --auto-generate-sql --auto-generate-sql-load-type=mixed --auto-generate-sql-add-autoincrement --engine=innodb --number-of-queries=5000
```