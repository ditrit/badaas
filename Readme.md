# Access Control Web App

This is a simple web application that demonstrates user authentication and access control using Casbin and Go. The application allows users to log in, select a domain, specify an object, and an action, and then enforces access control policies to determine whether the requested action is allowed.

## Features

- User Authentication and Session Management: Securely manage user sessions and allow basic authentication.
- Domain Selection and Input of Object and Action: Users can select domains and specify objects and actions.
- Enforcement of Access Control Policies using Casbin: Utilizes Casbin to enforce access control policies and returns results.

## Configuration

- model.conf: Defines the access control model and matchers for a hybrid ABAC/RBAC with domains.
- policy.csv: Defines the access control policies.    

## Testing

To test, please enter a struct like object so it can be converted and avoid errors:

```go
{
  "Name": "file_name",
  "Owner": "owner_name",
  "Domain": "domain_name",
  "Creator": "creator_name"
}
```
