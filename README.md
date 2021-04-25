# REST_API_example

  This app was made to demonstrate how golang interactive web application may be developed to сorresond *REST architecture*.
Server is made with Go package `net/http`. REST API was implemented using package `gorilla.mux`, which is able to create flexing routing.
For the sake of interest I reject using `GORM` and decide to use specific POSTGRESQL-compatible package [`jackc/pgx`](https://github.com/jackc/pgx) to create queries manually.  

---  

## Schema
![Schema's diagram](https://github.com/bondarenkoi07/REST_API_example/blob/master/schema.png)

  According to above schema, our app includes 4 models: `User`, `Developer`, `Market` and `Product`.  
**Business logic**: Model `User` incudes informaton required for authentication, other information is contained in dependent model `Developer`.
Model `Market` contains info place of  product's sale including info about maximum storage capacity.
Model `Product` depends on both `Market` and `Developer` models by id, but Market's also imposes constraint on count of dependent products.

---  
## API

Below is the implementation of the REST API:
Implemented CRUD functionality to route requests to handlers. It will complete different manipulations with data.
Routes have a similiar structure: path `/api`, name of used table, and, if required, id, filter, etc.  
All handlers are also mapped with semantically appropriated methods (*POST*, *GET*, *PUT*, *DELETE*)

```go
  	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.ReadAll).Methods(http.MethodGet)
	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.DeleteAll).Methods(http.MethodDelete)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.ReadById).Methods(http.MethodGet)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.DeleteById).Methods(http.MethodDelete)
	r.HandleFunc("/api/{table:[A-Za-z]+}/filter/{f_key:[A-Za-z]+}/{id:[0-9]+}", controller.FilterProducts).Methods(http.MethodGet)  
```

## Project Structure

### package Models

  This package includes implementations of DB tables. All them implements interface 
  [IterableModel](https://github.com/bondarenkoi07/REST_API_example/blob/c73f9fd1c8fda436400d040060d50593e12758ff/database/database.go#L12).
  This interface was made to [Create](https://github.com/bondarenkoi07/REST_API_example/blob/c73f9fd1c8fda436400d040060d50593e12758ff/database/database.go#L41) 
  and [Update](https://github.com/bondarenkoi07/REST_API_example/blob/c73f9fd1c8fda436400d040060d50593e12758ff/database/database.go#L128)
  methods.
  
### package Database

```go
type Database struct {
	//пул подклчений
	pool *pgxpool.Pool
	ctx  context.Context
}

func New(user, password, connection, databaseName string) (*Database, error) {
	var ctx context.Context
	var pool *pgxpool.Pool

	ctx = context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, connection, databaseName)
	var err error
	pool, err = pgxpool.Connect(ctx, dsn)

	if err != nil {
		return nil, err
	} else {
		conn := &Database{pool: pool, ctx: ctx}
		return conn, nil
	}
}
```

  This package provides access to postgres DB. All it's methods are a wrappers for pgx funcs.  
  
### package Service

```go
type DeveloperService struct {
	dbp *database.Database
}

func NewDeveloperService(dbp *database.Database) *DeveloperService {
	return &DeveloperService{dbp: dbp}
}

func (ds DeveloperService) Create(model Models.Developer) error {
	err := ds.dbp.Create("developers", model)
	return err
}
```

  This package is responsive for data distribution. Methods `Deserialization` ([developer's](https://github.com/bondarenkoi07/REST_API_example/blob/c73f9fd1c8fda436400d040060d50593e12758ff/Service/DeveloperService.go#L140) in example)
is necessary to successful deserialization of recieved JSON. Every single model's service is contained in [MultiService](https://github.com/bondarenkoi07/REST_API_example/blob/master/Service/MultiService.go)

### package Controller

  This package is responsible for requests routing, REST API implementation is nested there. All methods implements basic `net/http` handler's interface.  
  ```go
  type Controller struct {
	service Service.Service
	Router  *mux.Router
}

func NewController() Controller {
	var controller Controller
	r := mux.NewRouter()
	controller.service = Service.NewService()
  
  //there is REST API implemetation
  //...
  
	controller.Router = r

	return controller
}
  ```

### package Utils

This package consist small func responible for final request and response handling (encode/decode JSON, writing HTTP body, etc.)
