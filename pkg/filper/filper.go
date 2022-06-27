//filper Fiber framework helper package
package filper

import "github.com/gofiber/fiber/v2"

func GetBadRequestError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
		Message: message,
	})
}
