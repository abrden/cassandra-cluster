# Cassandra cluster
Followed [this](http://tonyzampogna.com/2017/02/10/how-to-setup-an-apache-cassandra-cluster-2/) to create the Vagrant box

To bring the 3 Vagrant VMs up.
```
$ vagrant up –no-provision
$ vagrant provision
```

To get in node 1:
```
$ vagrant ssh cassandra1
```

To check for cluster status:
```
$ ~/apache-cassandra-3.11.4/bin/nodetool status
Datacenter: datacenter1
=======================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
--  Address        Load       Tokens       Owns (effective)  Host ID                               Rack
DN  192.168.50.41  164.84 KiB  256          100.0%            1480e5cc-4002-4426-bb12-038565b9ab53  rack1
UN  192.168.50.42  235.12 KiB  256          100.0%            ee651aaf-7e59-4659-b961-67b6acf9334e  rack1
UN  192.168.50.43  91.85 KiB  256          100.0%            2bf5d77c-9afa-40a0-9c58-724e92791fbc  rack1
```

To perform queries from node3: (Use .41 or .42 to do them from nodes 1 and 2)
```
$ ~/apache-cassandra-3.11.4/bin/cqlsh 192.168.50.43
Connected to My Cluster at 192.168.50.43:9042.
[cqlsh 5.0.1 | Cassandra 3.11.4 | CQL spec 3.4.4 | Native protocol v4]
Use HELP for help.
cqlsh> USE usertest;
cqlsh:usertest> SELECT * FROM users;

 user_id | email                | fname | lname
---------+----------------------+-------+---------
    1001 | john.doe@example.com |  john |     doe
    1002 |      bob@example.com |   bob | johnson
    1000 |  smith_j@example.com |  john |   smith

(3 rows)
cqlsh:usertest>  INSERT INTO users (user_id, fname, lname, email) VALUES (1003, 'bobby', 'johanson', 'bob@example.com');
cqlsh:usertest> SELECT * FROM users;

 user_id | email                | fname | lname
---------+----------------------+-------+----------
    1001 | john.doe@example.com |  john |      doe
    1003 |      bob@example.com | bobby | johanson
    1002 |      bob@example.com |   bob |  johnson
    1000 |  smith_j@example.com |  john |    smith

(4 rows)
cqlsh:usertest> exit
```

To create tables and stuff:
```
$ ~/apache-cassandra-3.11.4/bin/cqlsh 192.168.50.43
Connected to My Cluster at 192.168.50.43:9042.
[cqlsh 5.0.1 | Cassandra 3.11.4 | CQL spec 3.4.4 | Native protocol v4]
Use HELP for help.
cqlsh> SELECT cluster_name, listen_address FROM system.local;

 cluster_name | listen_address
--------------+----------------
   My Cluster |  192.168.50.43

(1 rows)
cqlsh> CONSISTENCY ONE
Consistency level set to ONE.
cqlsh> SHOW HOST
Connected to My Cluster at 192.168.50.43:9042.
cqlsh> CONSISTENCY
Current consistency level is ONE.
cqlsh>  CREATE KEYSPACE usertest WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };
cqlsh> USE usertest;
cqlsh:usertest> CREATE TABLE users (
            ...    user_id int,
            ...    fname text,
            ...    lname text,
            ...    email text,
            ...    primary key (user_id)
            ...  );
cqlsh:usertest> INSERT INTO users (user_id, fname, lname, email) VALUES (1000, 'john', 'smith', 'smith_j@example.com');
cqlsh:usertest>  INSERT INTO users (user_id, fname, lname, email) VALUES (1001, 'john', 'doe', 'john.doe@example.com');
cqlsh:usertest>  INSERT INTO users (user_id, fname, lname, email) VALUES (1002, 'bob', 'johnson', 'bob@example.com');
cqlsh:usertest> SELECT * FROM users;

 user_id | email                | fname | lname
---------+----------------------+-------+---------
    1001 | john.doe@example.com |  john |     doe
    1002 |      bob@example.com |   bob | johnson
    1000 |  smith_j@example.com |  john |   smith

(3 rows)
cqlsh:usertest> exit
```
