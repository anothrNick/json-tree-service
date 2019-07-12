Docker postgres instance to play with JSONB field.

### Run Locally

```sh
# start services
$ docker-compose up -d

# login to container and create table
$ docker exec -it json-tree-service_mydb_1 /bin/sh
# psql -U postgres
psql (12beta2 (Debian 12~beta2-1.pgdg100+1))
Type "help" for help.

postgres=# CREATE DATABASE testdb;
CREATE DATABASE
postgres=# exit
# exit

# login from local machine
$ psql -h localhost -p 5432 -U postgres -W
Password for user postgres: 
psql (9.5.2, server 12beta2 (Debian 12~beta2-1.pgdg100+1))
Type "help" for help.

postgres=# \l
                          List of databases
   Name    |  Owner   | Encoding |         Access privileges          
-----------+----------+----------+------------------------------------
 postgres  | postgres | UTF8     | 
 template0 | postgres | UTF8     | =c/postgres\npostgres=CTc/postgres
 template1 | postgres | UTF8     | =c/postgres\npostgres=CTc/postgres
 testdb    | postgres | UTF8     | 
(4 rows)

postgres=# \c testdb
```

### Database Schema

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE trees (
    id uuid DEFAULT uuid_generate_v4(), 
    project varchar NOT NULL, 
    data jsonb, 
    PRIMARY KEY (id)
);
```

### Test Data

```sql
INSERT INTO trees (project, data) VALUES ('people', '{"name": "bob", "age": 25, "friends": ["one", "two"], "job": {"title": "clerk"}}');
```

### Test Queries

```sql
-- select all data
SELECT project, data FROM trees WHERE project='people';
 people  | {"age": 25, "job": {"title": "clerk"}, "name": "bob", "friends": ["one", "two"]}

-- select single field
SELECT project, data->'name' FROM trees WHERE project='people';
 people  | "bob"

-- select objects with #>
-- select nests objects with keys/array indices
SELECT project, data#>'{friends,0}' FROM trees WHERE project='people';
 people  | "one"

SELECT project, data#>'{name}' FROM trees WHERE project='people';
 people  | "bob"

SELECT project, data#>'{friends}' FROM trees WHERE project='people';
 people  | ["one", "two"]

SELECT project, data#>'{job}' FROM trees WHERE project='people';
 people  | {"title": "clerk"}

SELECT project, data#>'{job, title}' FROM trees WHERE project='people';
 people  | "clerk"

SELECT project, data#>'{age}' FROM trees WHERE project='people';
 people  | 25

-- update key with jsonb_set, default creates if json path does not exist in {<json field/index path>}
-- update job title key
UPDATE trees set data=jsonb_set(data, '{job, title}', '"engineer"') where project='people';UPDATE 1SELECT project, data#>'{job, title}' FROM trees WHERE project='people'; people  | "engineer"

-- new field (meant to update age but typod)
UPDATE trees set data=jsonb_set(data, '{jage}', '27') where project='people';
UPDATE 1

SELECT project, data#>'{age}' FROM trees WHERE project='people';
 people  | 25

-- query above created new field `jage` (typo)
SELECT project, data#>'{}' FROM trees WHERE project='people';
 people  | {"age": 25, "job": {"title": "engineer"}, "jage": 27, "name": "bob", "friends": ["one", "two"]}

UPDATE trees set data=jsonb_set(data, '{age}', '27') where project='people';
UPDATE 1
UPDATE trees set data=jsonb_set(data, '{age}', '"test"') where project='people';
UPDATE 1

-- update nested array index 2
UPDATE trees set data=jsonb_set(data, '{friends, 2}', '"three"') where project='people';
UPDATE 1
SELECT project, data#>'{}' FROM trees WHERE project='people';
 people  | {"age": "test", "job": {"title": "engineer"}, "jage": 27, "name": "bob", "friends": ["one", "two", "three"]}

UPDATE trees set data=jsonb_set(data, '{friends, 4}', '"five"') where project='people';
UPDATE 1
SELECT project, data#>'{}' FROM trees WHERE project='people';
 people  | {"age": "test", "job": {"title": "engineer"}, "jage": 27, "name": "bob", "friends": ["one", "two", "three", "five"]}

SELECT project, data#>'{friends, one}' FROM trees WHERE project='people';
 people  | 

UPDATE trees set data=jsonb_set(data, '{friends, four}', '"dataaaaa"') where project='people';
ERROR:  path element at position 2 is not an integer: "four"
```

### TODO

Load schema on image startup.
