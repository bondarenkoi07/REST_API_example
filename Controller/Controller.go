package Controller

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/Service"
	"app/REST_API_example/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Controller struct {
	service Service.Service
	Router  *mux.Router
}

func NewController() Controller {
	var controller Controller
	r := mux.NewRouter()
	controller.service = Service.NewService()

	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.ReadAll).Methods(http.MethodGet)
	r.HandleFunc("/api/{table:[A-Za-z]+}", controller.DeleteAll).Methods(http.MethodDelete)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.ReadById).Methods(http.MethodGet)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/{table:[A-Za-z]+}/{id:[0-9]+}", controller.DeleteById).Methods(http.MethodDelete)
	r.HandleFunc("/api/{table:[A-Za-z]+}/filter/{f_key:[A-Za-z]+}/{id:[0-9]+}", controller.FilterProducts).Methods(http.MethodGet)

	controller.Router = r

	return controller
}

func (c Controller) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		pathVars := mux.Vars(r)
		table, isSet := pathVars["table"]
		if !isSet || table == "" {
			log.Println("maybe this?")
			http.Error(w, "wrong path", http.StatusBadRequest)
			return
		}
		var err error
		switch table {
		case "product":
			var data = make(map[string]string)
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "wrong body", http.StatusBadRequest)
				return
			}

			var user Models.Product

			err, user = c.service.ProductService.Deserialize(data, *c.service.DeveloperService, *c.service.MarketService)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = c.service.ProductService.Create(user)
		case "users":
			var data = make(map[string]string)
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "wrong body", http.StatusBadRequest)
				return
			}

			var user Models.User

			err, user = c.service.UserService.Deserialize(data)

			log.Println(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = c.service.UserService.Create(user)
		case "developers":
			var data = make(map[string]string)
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "wrong body", http.StatusBadRequest)
				return
			}

			var user Models.Developer

			err, user = c.service.DeveloperService.Deserialize(data, *c.service.UserService)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = c.service.DeveloperService.Create(user)
		case "markets":
			var data = make(map[string]string)
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "wrong body", http.StatusBadRequest)
				return
			}

			var user Models.Market

			user, err = c.service.MarketService.Deserialize(data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = c.service.MarketService.Create(user)
		default:
			http.Error(w, "table does not exist", http.StatusBadRequest)
		}

		utils.RespondJSON(w, map[string]string{"status": "created"}, err)
	}
}

func (c Controller) ReadAll(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	var err error
	var data interface{}
	switch pathVars["table"] {
	case "markets":
		data, err = c.service.MarketService.ReadAll()
	case "product":
		data, err = c.service.ProductService.ReadAll()
	case "developers":
		data, err = c.service.DeveloperService.ReadAll()
	case "users":
		data, err = c.service.UserService.ReadAll()
	default:
		http.Error(w, "table does not exist", http.StatusBadRequest)
	}

	utils.RespondJSON(w, data, err)
}

func (c Controller) ReadById(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	var data interface{}
	id, err := strconv.Atoi(pathVars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch pathVars["table"] {
	case "markets":
		data, err = c.service.MarketService.ReadOne(int64(id))
	case "product":
		data, err = c.service.ProductService.ReadOne(int64(id))
	case "developers":
		data, err = c.service.DeveloperService.ReadOne(int64(id))
	case "users":
		data, err = c.service.UserService.ReadOne(int64(id))
	default:
		http.Error(w, "table does not exist", http.StatusBadRequest)
	}

	utils.RespondJSON(w, data, err)
}

func (c Controller) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		pathVars := mux.Vars(r)
		table, isSetTable := pathVars["table"]

		id, isSetId := pathVars["id"]

		if !isSetTable || table == "" || !isSetId || id == "" {
			http.Error(w, "wrong path", http.StatusBadRequest)
			return
		}

		Id, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "wrong id", http.StatusBadRequest)
			return
		}

		switch table {
		case "users":
			var data = make(map[string]string)
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "wrong body", http.StatusBadRequest)
				return
			}

			var user Models.User

			err, user = c.service.UserService.Deserialize(data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = c.service.UserService.Update(user, int64(Id))
		case "developers":
			var data = make(map[string]string)
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "wrong body", http.StatusBadRequest)
				return
			}

			var user Models.Developer

			err, user = c.service.DeveloperService.Deserialize(data, *c.service.UserService)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = c.service.DeveloperService.Update(user, int64(Id))
		case "product":
			var data = make(map[string]string)
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "wrong body", http.StatusBadRequest)
				return
			}

			var user Models.Product

			err, user = c.service.ProductService.Deserialize(data, *c.service.DeveloperService, *c.service.MarketService)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = c.service.ProductService.Update(user, int64(Id))
		case "markets":
			var data = make(map[string]string)
			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "wrong body", http.StatusBadRequest)
				return
			}

			var user Models.Market

			user, err = c.service.MarketService.Deserialize(data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = c.service.MarketService.Update(user, int64(Id))

		default:
			http.Error(w, "table does not exist", http.StatusBadRequest)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.RespondJSON(w, map[string]string{"status": "updated"}, err)
	}
}

func (c Controller) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		pathVars := mux.Vars(r)
		var err error
		switch pathVars["table"] {
		case "markets":
			err = c.service.MarketService.DeleteAll()
		case "product":
			err = c.service.ProductService.DeleteAll()
		case "developers":
			err = c.service.DeveloperService.DeleteAll()
		case "users":
			err = c.service.UserService.DeleteAll()
		default:
			http.Error(w, "table does not exist", http.StatusBadRequest)
		}
		utils.RespondJSON(w, map[string]string{"status": "deleted"}, err)
	}
}

func (c Controller) DeleteById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		pathVars := mux.Vars(r)
		id, err := strconv.Atoi(pathVars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch pathVars["table"] {
		case "markets":
			err = c.service.MarketService.DeleteOne(int64(id))
		case "product":
			err = c.service.ProductService.DeleteOne(int64(id))
		case "developers":
			err = c.service.DeveloperService.DeleteOne(int64(id))
		case "users":
			err = c.service.UserService.DeleteOne(int64(id))
		default:
			http.Error(w, "table does not exist", http.StatusBadRequest)
		}
		utils.RespondJSON(w, map[string]string{"status": "deleted"}, err)
	}
}

func (c Controller) FilterProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		pathVars := mux.Vars(r)
		id, err := strconv.Atoi(pathVars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var data interface{}
		switch pathVars["table"] {
		case "product":
			switch pathVars["f_key"] {
			case "markets":
				data, err = c.service.ProductService.FilterProductsByMarket(int64(id))
			default:
				http.Error(w, "where is any constraints based on this table or table does not exist", http.StatusBadRequest)
			}
		default:
			http.Error(w, "where is any constraints based on this table or table does not exist", http.StatusBadRequest)
		}

		log.Println(data)

		utils.RespondJSON(w, data, err)
	}
}
