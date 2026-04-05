# Theory of Operation

Purpose: Describes the data flow between the store and handlers components within the userservice package.

Tags: architecture, design

## Data Flow

An incoming HTTP request is routed by the server mux in `main.go` to the appropriate handler in the `handlers` package. Each handler parses the request, calls the corresponding method on the `store`, and writes a JSON response.

```
HTTP Client
    |
    v
+---------------------------+
|  main.go (net/http mux)   |
|  :8080                    |
+---------------------------+
    |
    | routes request to
    v
+---------------------------+
|  handlers package         |
|  CreateUser               |
|  GetUser                  |
|  UpdateUser               |
|  DeleteUser               |
+---------------------------+
    |
    | calls
    v
+---------------------------+
|  store package            |
|  in-memory map[string]User|
|  protected by sync.RWMutex|
+---------------------------+
```

## Request Lifecycle

```
Request arrives
    |
    +-- POST /users ---------> handlers.CreateUser
    |                              |-> store.Create(user) -> assigns UUID -> returns User
    |                              |<- 201 JSON User
    |
    +-- GET /users/{id} -----> handlers.GetUser
    |                              |-> store.Get(id) -> returns (User, bool)
    |                              |<- 200 JSON User   (found)
    |                              |<- 404             (not found)
    |
    +-- PUT /users/{id} -----> handlers.UpdateUser
    |                              |-> store.Update(id, user) -> returns (User, bool)
    |                              |<- 200 JSON User   (found)
    |                              |<- 404             (not found)
    |
    +-- DELETE /users/{id} --> handlers.DeleteUser
                                   |-> store.Delete(id) -> returns bool
                                   |<- 204             (found)
                                   |<- 404             (not found)
```
