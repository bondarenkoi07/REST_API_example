package Test

import (
	"app/REST_API_example/Controller"
	"app/REST_API_example/Models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var controller Controller.Controller

var id int64

var ts = []string{
	`{
		  "login": "login",
		  "password": "password"
		}`,
	`{
			"login": "login1",
			"password": "password2"
		}`,
}

func TestCreateHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	controller = Controller.NewController()

	rr := httptest.NewRecorder()

	controller.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error(rr.Code, rr.Body)
	}

	for _, value := range ts {
		req, err = http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(value))
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()

		controller.Router.ServeHTTP(rr, req)

		var test map[string]string

		err = json.NewDecoder(rr.Body).Decode(&test)

		if rr.Code != http.StatusOK || test["status"] != "created" {
			t.Error(rr.Code, rr.Body)
		}

	}
}

func TestReadHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	controller.Router.ServeHTTP(rr, req)

	var data = make([]Models.User, 0)

	if rr.Body == nil {
		t.Error("body is nil", rr.Body.String())
	}

	err = json.Unmarshal([]byte(rr.Body.String()), &data)

	if err != nil {
		t.Error(err, rr.Body.String())
	}

	if rr.Code != http.StatusOK {
		t.Error("Shieeeeeeeet", rr.Code)
	}

	id = data[0].Id

}

func TestUpdateHandler(t *testing.T) {
	Id := strconv.Itoa(int(id))
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/users/%s", Id), strings.NewReader(`{"login":"loginany","password": "any"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	controller.Router.ServeHTTP(rr, req)

	var test map[string]string

	err = json.NewDecoder(rr.Body).Decode(&test)

	if rr.Code != http.StatusOK || test["status"] != "updated" {
		t.Error(rr.Code, rr.Body)
	}

}

func TestDeleteHandler(t *testing.T) {
	Id := strconv.Itoa(int(id))
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/users/%s", Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	controller.Router.ServeHTTP(rr, req)

	var test map[string]string

	err = json.NewDecoder(rr.Body).Decode(&test)

	if rr.Code != http.StatusOK || test["status"] != "deleted" {
		t.Error(rr.Code, rr.Body)
	}

	req, err = http.NewRequest(http.MethodDelete, "/api/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	controller.Router.ServeHTTP(rr, req)

	err = json.NewDecoder(rr.Body).Decode(&test)

	if rr.Code != http.StatusOK || test["status"] != "deleted" {
		t.Error(rr.Code, test["status"])
	}

	req, err = http.NewRequest(http.MethodDelete, "/api/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	controller.Router.ServeHTTP(rr, req)

	err = json.NewDecoder(rr.Body).Decode(&test)

	if rr.Code != http.StatusBadRequest {
		t.Error(rr.Code, rr.Body.String())
	}
}
