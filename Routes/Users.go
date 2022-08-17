package Routes

import (
	"example.com/example/Controllers"
	"github.com/gofiber/fiber/v2"
)

func UsersRoute(route fiber.Router) {
	route.Post("/validate_username/:username", Controllers.UsersValidateUsersUsername)
	route.Post("/validate_email/:email", Controllers.UsersValidateUsersEmail)
	route.Post("/new", Controllers.UsersSignUp)
	route.Post("/register", Controllers.UsersRegister)
	route.Post("/login", Controllers.UsersLogin)
	route.Post("/get_user_by_id/:id", Controllers.UsersGetByID)
	route.Post("/get_all", Controllers.UsersGetAll)
}
