# monGO
A cli showing basic usage of the offical go driver for mongo db.

## Usage
### uri
Use env variable `MONGODB_URI` to set the mongo uri, if unset `mongodb://localhost:27017` is used by default.

### commands
    $ ./main -cmd=ping
    connecting to: mongodb://localhost:27017
    Successfully pinged!

The available commands are
##### ping
test that the connection is working
##### count
count the number of database entries
##### insert
generate a random entry
##### insertMany
generate `-num=<int>` random entires
##### view
retrieve object with object id: `-id`
##### list
list objects with default sort, use `-limit=<int>` and `-offset=<int>` for pagination
##### delete
delete object with object id: `-id`
##### filter
Search for objects; use `-filter=<field name>` to set the search field and `-term=<value>` to define the search term. Use `-sort=<field name>` to define the field to use for sorting and `-dir=<ASC|DESC>` to specificy ascending or descending sort. And use `-limit=<int>` and `-offset=<int>` for pagination.

### Example
##### generate entries
    $ ./main -cmd=insertMany -num 15
    connecting to: mongodb://localhost:27017
    Inserted document with _id: ObjectID("61f857e33301e08580fba374")
    Inserted document with _id: ObjectID("61f857e33301e08580fba375")
    Inserted document with _id: ObjectID("61f857e33301e08580fba376")
    Inserted document with _id: ObjectID("61f857e33301e08580fba377")
    Inserted document with _id: ObjectID("61f857e33301e08580fba378")
    Inserted document with _id: ObjectID("61f857e33301e08580fba379")
    Inserted document with _id: ObjectID("61f857e33301e08580fba37a")
    Inserted document with _id: ObjectID("61f857e33301e08580fba37b")
    Inserted document with _id: ObjectID("61f857e33301e08580fba37c")
    Inserted document with _id: ObjectID("61f857e33301e08580fba37d")
    Inserted document with _id: ObjectID("61f857e33301e08580fba37e")
    Inserted document with _id: ObjectID("61f857e33301e08580fba37f")
    Inserted document with _id: ObjectID("61f857e33301e08580fba380")
    Inserted document with _id: ObjectID("61f857e33301e08580fba381")
    Inserted document with _id: ObjectID("61f857e33301e08580fba382")
    Documents inserted: 15

##### view an entry
    $ ./main -cmd=view -id=61f857e33301e08580fba382
    connecting to: mongodb://localhost:27017
    61f857e33301e08580fba382       Brad      Smith    Gresham   WA

##### list entries
    $ ./main -cmd=list -offset=0 -limit=5          
    connecting to: mongodb://localhost:27017
    61f5a9d9ed57f46a4a7f3e63       Brad    Simpson    Gresham   WA
    61f5a9ea30f8fc48591e223b     Ashley   Peterson    Gresham   WA
    61f5a9ea30f8fc48591e223c       John    Johnson    Gresham   WA
    61f5a9ea30f8fc48591e223d       Brad    Simpson    Gresham   WA
    61f5a9ea30f8fc48591e223e     Heater    Simpson   Portland   WA

##### filter
    $ ./main -cmd=filter -filter=address.state -term=OR -sort=address.city -dir ASC
    connecting to: mongodb://localhost:27017
    61f8598a84389043a523553d       Brad    Johnson    Gresham   OR
    61f8598a84389043a5235548      Laura    Johnson    Gresham   OR
    61f8598a84389043a5235545       Brad      Smith    Gresham   OR
    61f8598a84389043a5235541    Heather    Simpson    Gresham   OR
    61f8598a84389043a5235543      Laura   Peterson    Gresham   OR