# Testbed for OpenIDConnect

The two OIDC providers chosen for this testbed are Google and Gitlab. Thus you can choose to connect using Google or Gitlab authentication.

The frontend is written in Quasar and the backend is written in Golang.

For the testbed to work, you have to declare an OIDC application on both providers, and insert the ```CLIENT_ID``` and ```CLIENT_SECRET``` in the ***backend/conf.env*** file. When creating the application, you will be asked for a ```callback_URI```. You have to provide the URI of the frontend page which will receive the callback : http://localhost:8080/callback

You can launch the backend server by running the command ```go run .``` inside the ***backend*** folder.
You can launch the frontend application by running the command ```quasar dev``` inside the ***frontend*** folder.

## What is stored in the localStorage of the browser

When a user logs in using the frontend app, the backend fetches the OIDC tokens (```id_token```, ```refresh_token``` and ```access_token```) from the provider and stores them, they are never sent to the frontend as it is considered to be an unsafe environment. For each user session, the backend generates a ```session_code``` which is sent to the frontend localStorage. Each time the frontend wants to access the backend API, it has to send the ```session_code``` as an **Authorization: Bearer**. The backend can then checks the validity of the OIDC tokens corresponding to this ```session_code```.

## How to add a new provider

You can see the list of OAuth and OIDC providers here : [List of OAuth providers](https://en.wikipedia.org/wiki/List_of_OAuth_providers)

### In the backend

First you need to declare an application on the provider services.
Then, in the ```conf.env``` file, you need to add three environment variables :
```
PROVIDERNAME_CLIENT_ID= ...
PROVIDERNAME_CLIENT_SECRET= ...
PROVIDERNAME_ISSUER= ...
```
In the ```provider.go``` file, change the ***GetProviders()*** function to get the ```providerNameConfig, providerNameVerifier, providerNameProvider``` of the new provider. Then, change the ***CreateProvider(name string)*** to cover the instantiation of the new provider.
You now need to create a file ```providerName.go```. This file must follow the pattern of ```google.go``` and ```gitlab.go```. All the functions listed in the **Provider** interface must be implemented.

### In the frontend

For the frontend you need to change the ```LoginPage.vue``` file. You just need to add a button that will call the ```LoginClick(prov)``` function with the name of the new provider :
```javascript
<q-btn @click="() => LoginClick('providerName')" color="primary" icon="login" label="Login with providerName" />
```

Now you are good, your new provider is configured !
