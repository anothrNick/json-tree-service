Docker postgres instance to play with JSONB field. Spun up with docker-compose

### Run Locally

```sh
# start services
$ make up

# login with psql
$ psql -h localhost -p 5432 -U testuser -W -d testdb
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

-- Delete nested key from object
SELECT project, data#>'{}' as data FROM trees WHERE project='mydb';
 project |                                       data                                       
---------+----------------------------------------------------------------------------------
 mydb    | {"age": 25, "job": {"title": "clerk"}, "name": "bob", "friends": ["one", "two"]}
(1 row)

SELECT data #- '{friends,0}' as data FROM trees WHERE project='mydb';
                                   data                                    
---------------------------------------------------------------------------
 {"age": 25, "job": {"title": "clerk"}, "name": "bob", "friends": ["two"]}
(1 row)

SELECT project, data#>'{}' as data FROM trees WHERE project='mydb';
 project |                                       data                                       
---------+----------------------------------------------------------------------------------
 mydb    | {"age": 25, "job": {"title": "clerk"}, "name": "bob", "friends": ["one", "two"]}
(1 row)

-- Actually delete the nested object by updating the column
UPDATE trees SET data=(SELECT data #- '{friends,0}' as data FROM trees WHERE project='mydb') where project='mydb';
UPDATE 1
SELECT project, data#>'{}' as data FROM trees WHERE project='mydb';
 project |                                   data                                    
---------+---------------------------------------------------------------------------
 mydb    | {"age": 25, "job": {"title": "clerk"}, "name": "bob", "friends": ["two"]}

-- Perhaps more efficient
UPDATE trees SET data=data #- '{friends,0}' where project='mydb';
UPDATE 1
SELECT project, data#>'{}' as data FROM trees WHERE project='mydb';
 project |                                 data                                 
---------+----------------------------------------------------------------------
 mydb    | {"age": 25, "job": {"title": "clerk"}, "name": "bob", "friends": []}
(1 row)

```
