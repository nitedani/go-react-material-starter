package users

import (
	"net/http"
	"server/generated/db"

	"github.com/labstack/echo/v5"
)

func Setup(g *echo.Group) {
	user_g := g.Group("/users")

	user_g.POST("", createUser)
	user_g.GET("", getUsers)
	user_g.GET("/:id", getUser)
}

func createUser(c echo.Context) error {
	newUser := new(db.UserModel)
	c.Bind(newUser)
	createdUsed, _ := db.Client.User.CreateOne(
		db.User.Provider.Set(db.ProviderLOCAL),
		db.User.Username.Set(newUser.Username),
	).Exec(db.Ctx)
	return c.JSON(http.StatusCreated, createdUsed)
}

func getUsers(c echo.Context) error {
	users, _ := db.Client.User.FindMany().Exec(db.Ctx)
	return c.JSON(http.StatusOK, users)
}

func getUser(c echo.Context) error {
	user, _ := db.Client.User.FindUnique(
		db.User.ID.Equals((c.PathParam("id"))),
	).Exec(db.Ctx)

	return c.JSON(http.StatusOK, user)
}
