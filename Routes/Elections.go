package Routes

import (
	"example.com/example/Controllers"
	"github.com/gofiber/fiber/v2"
)

func ElectionsRoute(route fiber.Router) {
	route.Post("/new", Controllers.ElectionsNew)
	route.Post("/modify/:id", Controllers.ElectionsModify)
	route.Post("/get_all", Controllers.ElectionsGetAll)
	route.Post("/get_all_populated", Controllers.ElectionsGetAllPopulated)
}
