// filper Fiber framework helper package
package filper

import "github.com/gofiber/fiber/v2"

func GetBadRequestError(c *fiber.Ctx, message any) error {
	if msg, ok := message.(string); ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": msg})
	}
	return c.Status(fiber.StatusBadRequest).JSON(message)
}

func GetNotFoundError(c *fiber.Ctx, message string) error {

	if len(message) < 1 {
		message = "not found"
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": message,
	})
}

func GetInternalError(c *fiber.Ctx, message string) error {

	if len(message) < 1 {
		message = "somthing unexpected happened"
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": message,
	})
}

func GetSuccessResponse(c *fiber.Ctx, message string) error {

	if len(message) < 1 {
		message = "somthing unexpected happened"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func GetUnAuthError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": message,
	})
}
