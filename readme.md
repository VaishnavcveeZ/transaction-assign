
# Run

    go build
    ./transaction-assign.exe

# test
     go test -v --cover

## Create cover.out to check function coverage
    go test ./... -coverprofile cover.out

## view function coverage
    go tool cover -func cover.out

# API 
## Create Transaction
> POST: /transaction

- body: json
        
        {
            "amount":788.3,
            "timestamp":"2023-02-10T20:59:00.312Z"
        }

        success code: 201

## Delete Transaction
> DELETE: /transaction

         success code: 204

## Get Statics
> GET: /statics


    - Header: {
                location:bangalore
            }

    - API response: {
                    "sum": 1576.6,
                    "avg": 788.3,
                    "max": 788.3,
                    "min": 788.3,
                    "count": 2
                }
    success code: 200
## Set User City
> POST: /location

     success code: 201

- body: json
        
        {
            "city":"bangalore"
        }

        success code: 201

## Reset City
> PUT: /location

    success code: 205
