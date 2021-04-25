# REST_API_example
  This app was made to demonstrate how golang interactive web application may be developed to сorresond REST architecture.
Server is made with Go package `net/http`. REST API was implemented using package `gorilla.mux`, which is able to create flexing routing.
For the sake of interest I reject using `GORM` and decide to use specific POSTGRESQL-compatible package `jackc/pgx` to create queries manually.  
---  
## Schema
![Schema's diagram](https://github.com/bondarenkoi07/REST_API_example/blob/master/schema.png)
---  
## API
```go
  	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.ReadAll).Methods(http.MethodGet)
	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.DeleteAll).Methods(http.MethodDelete)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.ReadById).Methods(http.MethodGet)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.DeleteById).Methods(http.MethodDelete)
	r.HandleFunc("/api/{table:[A-Za-z]+}/filter/{f_key:[A-Za-z]+}/{id:[0-9]+}", controller.FilterProducts).Methods(http.MethodGet)
```
