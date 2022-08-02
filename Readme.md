# Project structure

1. app/ - layer of app-building logic
2. app/tools/geoimport - main package for importer of geodata.
3. app/services/geoapi - main package for web api
4. app/services/geoapi/handlers - input layer
5. business/ - layer of business logic
6. business/core/{city,country,location} - layer of access to business entities
7. business/core/{city,country,location}/db - layer of access to database entities (city, country, location)
8. business/data - helpers to manage data (migrations, seeds) and setup tests ()
9. foundation/ - all non-business related logic
10. foundation/database/ - common database related helpers
11. foundation/docker/ - common docker related helpers
12. foundation/web/ - common web related helpers
13. infra/ - configuration for infrastructure

# Layers responsibilities (resps.)
```
                                Request  /  Response
---------------------------------------------------------------------------------------
                    Input/Output layer (app/services/geoapi/handlers)
                    
Resps.:   Decode request 
          Encode response
          Manage user-faced errors
---------------------------------------------------------------------------------------
                    Core layer (business/core/{city,country,location, importer})
                    
Resps.:   Validate input data 
          Prepare data to database layer
          Building result from database layer
          Handling errors from database layer 
---------------------------------------------------------------------------------------
                    Database layer (business/core/{city,country,location}/db)
                    
Resps.:   Communicate with mysql by queries
          Return entities from mysql to core layer
---------------------------------------------------------------------------------------
```


# Import csv-data logic

As file is quite big I have 3 ideas how to implement this:
1. Read all data in a loop, put data in 1 huge transaction and commit it. This way requires a lot of memory
   (cause of huge transaction) and slow
2. Read data in a loop and pass them to workers (I used 8 workers as I have 8 logical cores laptop). Workers
insert data without transactions. So no additional memory. However, requires additional time as indexes are being updated
on every database modification. 1M modifications - 10M indexes updates (as I used 10 indexes)
3. Read data in a loop and pass them to workers. Start as a previous one but with important change: every worker works
with defined set of country. Or to say another, all records for specific country will be handled by one specific worker.
This allows to use multiple transactions. One transaction per worker. However, during implementing this I experienced
an issues. I supposed it's a mysql driver bug. For some reasons, it holds transactions for 50s from time to time. So
I end up with the solution #2. In real life, I'd deep into this issue and debug mysql driver's code. Or try to eliminate
of this issue by  changes in logic. For instance, I can firstly import countries and cities in one transaction. And import
locations after this. This might save ~20min of execution time. Taking into account that there is no target time and this
data will be imported one a day I think result time ~1.5h is ok. We could import it at 4am without any issues. Of course,
in case of any additional requirements to time, pprof analyse of CPU usage will be done and code will be optimized for best
performance.


# Run
1. make all
2. make import
3. make api
