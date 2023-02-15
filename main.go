package main

import (
	"fmt"
	"net/http"
	"os"
	"vandesar/handler/api"
	"vandesar/middleware"
	"vandesar/repository"
	"vandesar/service"
	"vandesar/utils"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler *api.UserAPI
}

func main() {
	os.Setenv("DATABASE_URL", "postgres://root:secret@localhost:5432/pos")

	mux := http.NewServeMux()

	err := utils.ConnectDB()
	if err != nil {
		panic(err)
	}

	db := utils.GetDBConnection()
	mux = RunServer(db, mux)

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func RunServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)

	userAPIHandler := api.NewUserAPI(userService)

	apiHandler := APIHandler{UserAPIHandler: userAPIHandler}

	MuxRoute(mux, "POST", "/api/v1/users/admin/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.AdminLogin)))
	MuxRoute(mux, "POST", "/api/v1/users/admin/register", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.AdminRegister)))

	MuxRoute(mux, "POST", "/api/v1/users/cashier/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.CashierLogin)))
	MuxRoute(mux, "POST", "/api/v1/users/cashier/register", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.CashierRegister)))

	return mux
}

func MuxRoute(mux *http.ServeMux, method string, path string, handler http.Handler, opt ...string) {
	if len(opt) > 0 {
		fmt.Printf("[%s]: %s %v \n", method, path, opt)
	} else {
		fmt.Printf("[%s]: %s \n", method, path)
	}

	mux.Handle(path, handler)
}
