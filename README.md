# DrawSQL
Tool for drawing relationship diagrams from SQL databases.  It will connect ot a db, pull the table schemas and build 
a nice relationship chart similar to:

![Graph](https://github.com/RadhiFadlillah/sqldiagram/blob/master/example/basic/erd.png)

Only supports CockroachDB at the moment.  Feel free to add interfaces for other databases. ( See: pkg/db/interface.go )

I wrote this to scratch my own itch for something unrelated I'm working on.

# schema.json

The tool will attempt to load a config file called schema.json that allows you to define relationships between tables
by hand for cases where foreign keys are not being used, and lets you link the tables to each other.


```aiignore

Format:
{
   <table>: {
      <column name>: <related table>
   },
   ...
}

Example:
{
  "*": {
    "customer_id": "customers",
    "device_id":   "devices",
  },
  "users": {
    "customer": "customers"
  }
}
```


`*`: For all tables if a column is named customer_id, it's pointing to the customer table

For the `users` tbale the field named `customer` points to the customers table.


# Attribution:
I borrowed (lifted, copied, whatever you want to call it) a bunch of code from: https://github.com/RadhiFadlillah/sqldiagram 
(specifically the stuff in pkg/renderer) Give that project a look and a star.

# Licenses 

DrawSQL is distributed under MIT license, which means you can pretty much do whatever you want with it.  
If you do have changes, fixes, want to support more database, make a pull request. Or clone the repo and do whatever,
it's cool. 
