//filper Fiber framework helper package
package filper

import "github.com/gofiber/fiber/v2"

func GetBadRequestError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": message,
	})
}

func GetInternalError(c *fiber.Ctx, message string) error {

	if len(message) < 1 {
		message = "somthing unexpected happened"
	}

	return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
		"message": message,
	})
}

func GetInvalidCredentialsError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
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