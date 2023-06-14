# Casbin Hybrid ABAC/RBAC with Domains Model

Welcome to the repository for our hybrid access control model, implemented using the Casbin library in Golang. This model combines the principles of Role-Based Access Control (RBAC) with domains (or tenants) and Attribute-Based Access Control (ABAC) to manage access to data resources.

## Resources definition

Data resources are defined as a struct with four string properties:

```plaintext
type Data struct {
	Name    string
	Owner   string
	Domain  string
	Creator string
}
```

This structure allows us to check the "Owner" or "Creator" attributes of the subject requesting access within our matchers.

## Model Overview

Our model uses the following definitions:

```plaintext
[request_definition]
r = dom, sub, obj, act

[policy_definition]
p = dom, sub, obj, act

[role_definition]
g = _, _
g2= _, _
g3= _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.dom== p.dom && (r.obj.Name== p.obj || g(r.sub, p.sub)) && r.act == p.act && (g2(r.obj.Name, p.obj) || g3(r.sub, p.sub, r.dom)) || r.sub == r.obj.Owner && r.dom == r.obj.Domain && r.act == 'read' || r.sub == r.obj.Creator && r.dom == r.obj.Domain
```

The goal of our model is to provide access to users based on their roles and attributes. 

### ABAC model:

In the ABAC component of our model, we define the attributes "OWNER" and "CREATOR". The owner of a data resource is granted read-only access to their data, while the creator of a data resource has full control over their data. This is just an example, and the matchers can be modified to align with your specific policy rules.

### RBAC with domains model:

The RBAC with domains component of our model allows us to define different roles for the same user, depending on the resource. Each role is defined within a specific perimeter, or "Domain", to provide granular access control for different resources.

### Hierarchy in the model

Our model supports the establishment of hierarchical relationships, both among users and domains.

#### User Hierarchy:

We can group users and define hierarchical relationships among them using the "g3" matcher function. This allows us to create complex user structures, where certain users or groups of users have authority or permissions that others do not.

#### Domain Hierarchy:

The "g2" matcher function enables us to define a hierarchy among domains. This is particularly useful when you have a complex system with multiple domains, each with its own set of rules and permissions. Furthermore, our model supports nested domains, meaning that a domain can contain other domains as well as resources. This provides an additional layer of complexity and control, allowing you to fine-tune access and permissions based on the specific needs of your system.
