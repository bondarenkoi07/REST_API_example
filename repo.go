package REST_API_example

import "app/REST_API_example/Service"

func test(service Service.UserService) {
	print(service.ReadAll())
}

func init() {
	var service Service.UserService

	test(service)
}
