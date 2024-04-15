package handlers

import (
	"errors"

	"github.com/adrianowsh/bookfy-api/db"
	"github.com/adrianowsh/bookfy-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandlerInterface struct {
	userStore db.UserStoreInterface
}

func NewUserHandler(userStore db.UserStoreInterface) *UserHandlerInterface {
	return &UserHandlerInterface{
		userStore: userStore,
	}
}

func (h *UserHandlerInterface) HandlehPostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}

func (h *UserHandlerInterface) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params(("id"))
	)

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userID})
}

func (h *UserHandlerInterface) HandleDeleteUser(c *fiber.Ctx) error {
	var userID = c.Params("id")

	if err := h.userStore.RemoveUser(c.Context(), userID); err != nil {
		return err
	}

	return c.JSON(map[string]string{"deleted": userID})
}

func (h *UserHandlerInterface) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsersPaginated(c.Context())

	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *UserHandlerInterface) HandleGetUser(c *fiber.Ctx) error {
	var userID = c.Params("id")

	user, err := h.userStore.GetUserById(c.Context(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}

	return c.JSON(user)
}
