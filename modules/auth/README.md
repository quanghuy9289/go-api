# Auth module

Support authentication for project, include internal authentication (user register account and login to system) and external user (oauth with google, facebook...)
Build up single sign on for multiple services

## Main technical points

- Signed token with jwt: <https://github.com/dgrijalva/jwt-go>
- OAuth: goth - <https://github.com/markbates/goth>

## Solution details

### Internal authentication

- User register a new account with graphql
  
```graphql - register account

    mutation register {
        register(input: {
            email: "user1@gmail.com"
            fullname: "user1"
            nickname: "user1"
            password: "user1"
            avatarBase64: "user1"
            roleID: "user"
        }) {
            id
            email
        }
    }
```

- User login with registered account above

```graphql - login

    mutation login {
        login(input: {
            email: "user1@gmail.com"
            password: "user1"
        }) {
            token
            user {
                id
                email
                fullname
            }
        }
    }
```

After login, system will generate a jwt token for this session and save to AuthenticationToken table. The information for this record contains infor about current logged user, device id using and expired value for this token.

### External authentication

We support this for users that have account in google, facebook... and using our system without create new account (OAuth machanic). Currently we support oauth with google, others third party will be the same.

Steps:

- When user click to login with Google, we handle this by API `GET:{CTB-domain}/auth/provider/:provider` (parameter provider is `google` now). The handler for this can find at `modules/auth/handler/auth.go`
- After `goth` process authenticate, it will redirect user to select their google account to login.
- Then it calls a callback function that we config before in `InitalizeAuthProviders` method at `modules/auth/main.go`
- The callback function will complete authenticate for user, add to our database and return a `access token`.
- Finally, user can use the above `access_token` to do other requests that need authenticate.
