# Sandbox for access control following the RBAC model

This repository is a sandbox for the RBAC model. I am using the [Casbin package](https://casbin.org/docs/en/overview) for Golang.

## Structure
RBAC stands for Role-Based Access Control. But in my architecture I have extended the notion of role to the notion of team. Users belong to teams and 
each team is given a set of roles on domains. A domain is a set of resources.

A role indicates which actions can be performed on a type of resources (for example a type of resources can be ```application``` or ```database```).

## Example
The graph below shows a generic example of the structure that I wanted to build.

![rbac_graph](https://user-images.githubusercontent.com/102538155/180420826-0304c288-949e-4286-a19d-3cf37bf285c1.jpg)

You can see that ***user3*** belongs to ***team3***. He is thus ***admin*** on ***domain2*** and ***user***
on ***domain3***. Now, if you want to understand what actions can ***user3*** perform, you must refer to
the table on the top right corner. An admin can perform ```admin_action``` on resources of type 1 and 2.
And a user can perform ```user_action``` on resources of type 1 and 2. If we apply this to our user3,
we conclude that he can do ```admin_action``` on ```data2``` and ```user_action``` on ```data3```.

Another important use-case illustrated by this graph is the inheritance between two domains.
Domain 2 is contained in domain 1. This means that a user having a role on domain 1 will keep his role
on domain 2. Here, ***user1*** is ***admin*** on ***domain1***. As a consequence, he is also ***admin***
on ***domain2***.

On the contrary, ***user2*** is ***user*** on ***domain2*** but doesn't have any role on ***domain1***.

All the requests that can be performed on this example are listed below :

![Capture_requests_3users](https://user-images.githubusercontent.com/102538155/180424674-3cdd71e3-cb7d-471e-9e02-ef08717acddf.PNG)

In this configuration, every user has just the bare necessary permissions, which is what we want for security purposes.

## How to store our RBAC configuration
The simplest way to store the configuration is to seperate it into two files : ```model.conf``` and ```policy.csv```.
This is the easiest way of testing a configuration but this is not really suited for production.

Another way of storing the configuration is by leveraging on a database. You can keep the ```model.conf``` file as it will not scale.
But the ```policy.csv``` file will turn into a table. You can for example create a table ```policy``` in a database named ```rbac```.
For my tests I have used **CockroachDB**. You can create a cluster using this command :
```
cockroach start --insecure --store=node1 --listen-addr=localhost:26257 --http-addr=localhost:8080 --join=localhost:26257,localhost:26258,localhost:26259 --accept-sql-without-tls
```
Then initialize the cluster :
```
cockroach init --insecure
```
Enter the CLI of your cluster :
```
cockroach sql --host=localhost:26257 --insecure
```
Be careful when creating the table, you have to use these commands :
```sql
CREATE DATABASE rbac;
use rbac;
CREATE TABLE public.policy (
  p_type VARCHAR(32) NOT NULL DEFAULT '':::STRING,
  v0 VARCHAR(255) NOT NULL DEFAULT '':::STRING,
  v1 VARCHAR(255) NOT NULL DEFAULT '':::STRING,
  v2 VARCHAR(255) NOT NULL DEFAULT '':::STRING,
  v3 VARCHAR(255) NOT NULL DEFAULT '':::STRING,
  v4 VARCHAR(255) NOT NULL DEFAULT '':::STRING,
  v5 VARCHAR(255) NOT NULL DEFAULT '':::STRING,
  INDEX idx_casbin_rule (p_type ASC, v0 ASC, v1 ASC),
  FAMILY "primary" (p_type, v0, v1, v2, v3, v4, v5, rowid)
);
```
Then, you can insert the lines of the ```policy.csv``` one by one :
```sql
INSERT INTO policy (p_type, v0, v1, v2, v3) VALUES ('p', 'admin', 'domain1', 'type1', 'admin_action');
...
```

For more information on how to store the configuration, please refer to the official documentation :
- [Model storage](https://casbin.org/docs/en/model-storage)
- [Policy storage](https://casbin.org/docs/en/policy-storage)
