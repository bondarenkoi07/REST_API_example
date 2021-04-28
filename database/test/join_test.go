package test

import (
	"app/REST_API_example/Models"
	"app/REST_API_example/Service"
	"app/REST_API_example/database"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"testing"
)

var (
	TestUserModel       = &Models.User{Login: "foo14", Password: "bar14"}
	SecondTestUserModel = &Models.User{Login: "foo155", Password: "bar14555"}
	ThirdTestUserModel  = &Models.User{Login: "foo1551337", Password: "bar145551337"}

	TestDeveloperModel = &Models.Developer{
		OrgName: "roga",
		Section: "KoPbITA",
		UserId:  TestUserModel,
	}
	SecondTestDeveloperModel = &Models.Developer{
		OrgName: "techs inc.",
		Section: "tablets",
		UserId:  SecondTestUserModel,
	}
	ThirdTestDeveloperModel = &Models.Developer{
		OrgName: "techs corp.",
		Section: "tablets",
		UserId:  ThirdTestUserModel,
	}
	TestMarketModel = &Models.Market{Name: "roga", MaxProducts: 5}
	market1         = Models.Market{
		Name:        "techs",
		MaxProducts: 10,
	}
	market2 = Models.Market{
		Name:        "techs2",
		MaxProducts: 10,
	}
	product1 = Models.Product{
		Name:      "phone",
		Cost:      10,
		Count:     6,
		Market:    &market1,
		Developer: TestDeveloperModel,
	}
	product2 = Models.Product{
		Name:      "tablet",
		Cost:      10,
		Count:     6,
		Market:    &market1,
		Developer: TestDeveloperModel,
	}

	product3 = Models.Product{
		Name:      "phone",
		Cost:      10,
		Count:     6,
		Market:    &market2,
		Developer: SecondTestDeveloperModel,
	}
	product4 = Models.Product{
		Name:      "tablet",
		Cost:      10,
		Count:     6,
		Market:    &market1,
		Developer: ThirdTestDeveloperModel,
	}
	testProducts = []*Models.Product{
		&product1,
		&product2,
		&product3,
		&product4,
	}
	markets = []*Models.Market{
		&market1,
		&market2,
	}
	developers = []*Models.Developer{
		TestDeveloperModel,
		ThirdTestDeveloperModel,
	}
)

func init() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")
	connection := "db"
	if user == "" || password == "" || databaseName == "" {
		user = "postgres"
		password = "postgres"
		databaseName = "postgres"
	}
	var err error
	dbp, err = database.New(user, password, connection, databaseName)
	if err != nil {
		connection = "localhost"
		dbp, err = database.New(user, password, connection, databaseName)
		if err != nil {
			panic(err)
		}
	}

}

func TestManyToOne(t *testing.T) {

	err := dbp.DeleteAll("users")

	if err != nil {
		t.Error("couldn't vanish users: ", err)
	}

	err = dbp.DeleteAll("markets")

	if err != nil {
		t.Error("couldn't vanish markets: ", err)
	}

	defer func(dbp *database.Database, tableName string) {
		err := dbp.DeleteAll(tableName)
		if err != nil {
			t.Error(err)
		}
	}(dbp, "users")
	defer func(dbp *database.Database, tableName string) {
		err := dbp.DeleteAll(tableName)
		if err != nil {
			t.Error(err)
		}
	}(dbp, "markets")

	err = dbp.Create("users", TestUserModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	var rows *pgx.Rows

	rows, err = dbp.ReadAll("users")

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	defer (*rows).Close()

	for (*rows).Next() {
		err = (*rows).Scan(&TestUserModel.Id, &TestUserModel.Login, &TestUserModel.Password)
		if err != nil {
			t.Error("Cannot serialize row: ", TestUserModel)
		}
	}

	err = dbp.Create("developers", TestDeveloperModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	err = dbp.Create("markets", TestMarketModel)

	if err != nil {
		t.Error("Could not insert model: ", err)
	}

	rows, err = dbp.ReadAll("markets")

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	defer (*rows).Close()

	for (*rows).Next() {
		err = (*rows).Scan(&TestMarketModel.Id, &TestMarketModel.Name, &TestMarketModel.MaxProducts)
		if err != nil {
			t.Error("Cannot serialize row: ", TestUserModel)
		}
	}

	for i := 0; i < 5; i++ {
		var TestProductModel = &Models.Product{Name: fmt.Sprintf("name%d", i), Cost: int8(i * 10), Count: int8(1), Developer: TestDeveloperModel, Market: TestMarketModel}

		err = dbp.Create("product", TestProductModel)
	}

	row, err := dbp.GetProductsCount(TestMarketModel.Id)

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	id := 0
	cnt := 0
	err = (*row).Scan(&id, &cnt)
	if err != nil {
		t.Error("Cannot serialize row: ", TestUserModel, err)
	}

	if cnt != 5 {
		t.Error("Counted not properly enough", TestUserModel)
	}

	var list = make([]interface{}, 0)

	rows, err = dbp.GetMarketProducts(TestMarketModel.Id)

	if err != nil {
		t.Error("Could not read rows: ", err)
	} else if rows == nil {
		t.Error("nil pointer: ")
	}

	for (*rows).Next() {
		data, err := (*rows).Values()
		if err != nil {
			t.Error(err)
		}
		list = append(list, data)
	}

	if len(list) != 5 {
		t.Error("query fetched rows not properly")
	}
}

func Vanish(dbp *database.Database, t *testing.T) {
	defer func(dbp *database.Database, tableName string) {
		err := dbp.DeleteAll(tableName)
		if err != nil {
			t.Error(err)
		}
	}(dbp, "users")
	defer func(dbp *database.Database, tableName string) {
		err := dbp.DeleteAll(tableName)
		if err != nil {
			t.Error(err)
		}
	}(dbp, "markets")
}

func TestJoint(t *testing.T) {

	Vanish(dbp, t)
	var err error

	TestUserModel.Id, err = dbp.Save("users", TestUserModel)
	if err != nil {
		t.Error(err)
	}
	ThirdTestUserModel.Id, err = dbp.Save("users", ThirdTestUserModel)
	if err != nil {
		t.Error(err)
	}
	SecondTestUserModel.Id, err = dbp.Save("users", SecondTestUserModel)
	if err != nil {
		t.Error(err)
	}

	err = dbp.Create("developers", TestDeveloperModel)
	if err != nil {
		t.Error(err)
	}
	err = dbp.Create("developers", SecondTestDeveloperModel)
	if err != nil {
		t.Error(err)
	}
	err = dbp.Create("developers", ThirdTestDeveloperModel)
	if err != nil {
		t.Error(err)
	}

	for _, market := range markets {
		market.Id, err = dbp.Save("markets", market)
		if err != nil {
			t.Error(err)
		}
	}
	for _, product := range testProducts {
		err = dbp.Create("product", product)
	}
	rows, err := dbp.GetMarketDevelopers(market1.Id)
	if err != nil {
		t.Error(err)
	}
	var productCounts = make([]int64, 0)
	var prodSum = make([]int64, 0)
	for iter := 0; (*rows).Next(); iter++ {
		var iterModel Models.Developer
		var id int64
		var count int64
		var sum int64
		err = (*rows).Scan(&id, &iterModel.OrgName, &iterModel.Section, &count, &sum)
		if err != nil {
			t.Error(err)
		}

		userService := Service.NewUserService(dbp)

		user, err := userService.ReadOne(id)
		if err != nil {
			t.Error(err)
		}

		iterModel.UserId = user
		if user.Id != developers[iter].UserId.Id {
			t.Errorf("assertion error:%d, %d", user.Id, developers[iter].UserId.Id)
		}
		prodSum = append(prodSum, sum)
		productCounts = append(productCounts, count)
	}

	if len(productCounts) != 2 {
		t.Error("there odds developer(s) out :", len(productCounts))
		return
	}

	if productCounts[0] != 2 || productCounts[1] != 1 {
		t.Errorf("calculated counts is wrong: %d %d and %d %d", productCounts[0], 2, productCounts[1], 1)
	}
	if prodSum[0] != 12 || prodSum[1] != 6 {
		t.Errorf("calculated counts is wrong: %d %d and %d %d", prodSum[0], 12, prodSum[1], 6)
	}
}
