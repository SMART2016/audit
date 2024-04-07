## API documentation and Usage
- Run the service by going one folder up and firing docker compose up --build
- Ones all the services are up and running, you can click on the below link , that will show the swagger UI , from where you can directly open postman and start using the API's
- https://documenter.getpostman.com/view/5673453/2sA35MyJQV

## ES queries

QUERY1 = `{
"query": {
"range": {
"time": {
"gte": "now-10m",
"lte": "now"
}
}
}
}`

QUERY2 = `{
"query": {
"bool": {
"must": [
{
"match": {
"action": "HEALTH"
}
},
{
"match": {
"system": "user-service"
}
}
]
}
},
"_source": true
}

`
QUERY3 = `{"query":{"bool":{"must":[{"match":{"action":"HEALTH"}},{"match":{"system":"user-service"}},{"range":{"time":{"gte":"2024-04-06T00:00:00+05:30","lte":"2024-04-06T23:59:59+05:30"}}}],"should":[{"match":{"currentUser":"Dipanjan"}}],"minimum_should_match":1}}}`

QUERY4 = `{
"query": {
"bool": {
"must": [
{
"match": {
"action": "HEALTH"
}
},
{
"match": {
"system": "user-service"
}
},
{
"range": {
"time": {
"gte": "now-24h", // 5 minutes ago
"lte": "now"
}
}
}
],
"should": [
{
"match": {
"currentUser": "Dipanjan"
}
}
],
"minimum_should_match": 1
}
}
}`

QUERY5 = `
                {
                  "query": {
                    "range": {
                      "time": {
                        "gte": "now-24h",
                        "lte": "now"
                      }
                    }
                  }
                }`

QUERY6 = `{
"query": {
"bool": {
"must": [
{
"match": {
"currentUser": "dipanjanmazumder"
}
},
{
"term": {
"system.keyword": "monitoring-service"
}
},
{
"range": {
"time": {
"gte": "now-200h",
"lte": "now"
}
}
}
]
}
}
}
`

