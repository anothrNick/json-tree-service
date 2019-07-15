JSON Tree rest service. Access JSON structure with HTTP path parameters as keys/indices to the JSON.

### Run Locally

```sh
# start postgres service
$ make up

# start go web service, TODO: run in docker
$ go run backend/main.go

# make HTTP requests
$ curl -s -X POST -d '{"age": 25, "job": {"title": "clerk"}, "name": "bob", "friends": ["one", "two"]}' localhost:5000/mydb | jq "."
{}

$ curl -s localhost:5000/mydb/ | jq "."
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

$ curl -s localhost:5000/mydb/friends | jq "."
[
  "one",
  "two"
]

$ curl -s localhost:5000/mydb/friends/1 | jq "."
"two"

$ curl -s localhost:5000/mydb/job/title | jq "."
"clerk"
```
