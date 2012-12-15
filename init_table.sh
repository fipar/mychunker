#!/bin/bash
mysql -e 'drop table if exists test.test; create table test.test (id int unsigned not null auto_increment primary key, c int default null) engine=Innodb'
i=0
while [ $i -lt 1000000 ]; do
	echo "insert into test.test values (null,$RANDOM);"
	i=$((i+1))
done|mysql
