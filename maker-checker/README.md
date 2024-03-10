# maker-checker

1. Maker makes a transaction on the Admin Panel Ul. A transaction refers to any action
that is configured in the Maker-checker system, e.g. points adjustment.
2. Checker receives notification via email and approves the transaction. You can
implement the approval step in Admin Panel Ul or directly via the email.
3. Transaction is processed on the backend.


## transaction schema
```
{
    "action": {
        "resource": "points",
        "operation": {
            "type": "increment",
            "value": 10
            "where": {
                "user_id": ""
            }
        }
    },
    "maker": "",
    "description": "",
    "checker": "",
    "created_at": "",
    "updated_at": "",
}
```

## process

### create a transaction
makers --> lambda (mc) --> aurora (txn) --triggers?--> lambda (find checkers based on makers) --> ses --> checkers

### validate transaction
checkers --> lambda (mc) --> aurora (txn)
                         --> aurora (resource)

