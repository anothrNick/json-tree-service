JSON Tree rest service. Access JSON structure with HTTP path parameters as keys/indices to the JSON.

### Run Locally

```sh
# start services
$ make up
```

#### Make HTTP requests

```sh
# create tree for new project
$ curl -s -X POST -d '{"age": 25, "job": {"title": "clerk"}, "name": "bob", "friends": ["one", "two"]}' localhost:5001/api/mydb | jq "."
{}

# retrieve full tree
$ curl -s localhost:5001/api/mydb/ | jq "."
{
  "age": 25,
  "friends": [
    "one",
    "two"
  ],
  "job": {
    "title": "clerk"
  },
  "name": "bob"
}

# retrieve individual keys
$ curl -s localhost:5001/api/mydb/friends | jq "."
[
  "one",
  "two"
]

# or by index of array
$ curl -s localhost:5001/api/mydb/friends/1 | jq "."
"two"

$ curl -s localhost:5001/api/mydb/job/title | jq "."
"clerk"
```

#### Add key
```sh
$ curl -s -X POST -d '4' localhost:5001/api/mydb/job/years | jq "."
{}
$ curl -s localhost:5001/api/mydb/ | jq "."
{
  "age": 25,
  "friends": [
    "one",
    "two"
  ],
  "job": {
    "title": "clerk",
    "years": 4
  },
  "name": "bob"
}
```

#### Update key
```sh
$ curl -s -X PUT -d '{"title": "Engineer", "years": 1}' localhost:5001/api/mydb/job | jq "."
{}
$ curl -s localhost:5001/api/mydb/ | jq "."
{
  "age": 25,
  "friends": [
    "one",
    "two"
  ],
  "job": {
    "title": "Engineer",
    "years": 1
  },
  "name": "bob"
}
```

#### Delete key
```sh
$ curl -s -X DELETE localhost:5001/api/mydb/job | jq "."
{}
$ curl -s localhost:5001/api/mydb/ | jq "."
{
  "age": 25,
  "friends": [
    "one",
    "two"
  ],
  "name": "bob"
}
```

### Credit

* Websocket work basically taken from [gorilla/websocket/chat example](https://github.com/gorilla/websocket/blob/master/examples/chat)

### License

MIT Copyright (c) 2019 Nick Sjostrom