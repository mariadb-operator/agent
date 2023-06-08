#!/bin/bash

set -e

if [ -z "$ENTRYPOINT" ]; then
    ENTRYPOINT="/usr/local/bin/docker-entrypoint.sh"
fi

cat <<EOF | tee /etc/mysql/mariadb.conf.d/0-galera.cnf
[mysqld]
bind-address=0.0.0.0
default_storage_engine=InnoDB
binlog_format=row
innodb_autoinc_lock_mode=2

[galera]
wsrep_on=ON
wsrep_provider=/usr/lib/galera/libgalera_smm.so
wsrep_cluster_address="gcomm://mariadb"
wsrep_cluster_name="mariadb-galera-cluster"
wsrep_sst_method=rsync
wsrep_node_address="${HOSTNAME}"
wsrep_node_name="${HOSTNAME}"
EOF

if [ ! -n "$(ls -A /var/lib/mysql)" ]; then 
    bash -c "$ENTRYPOINT mariadbd --wsrep-new-cluster"
else
    bash -c "$ENTRYPOINT mariadbd"
fi