rsync -a -v --exclude .git ./ /mnt/c/Users/shinh/Desktop/gk@mysqld_exporter/

# kubectl set image deployment/mysql-exporter1 mysql-exporter1=shinhwagk/mysqld_exporter:0.0.2 -n monitoring