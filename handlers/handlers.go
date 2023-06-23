package handlers

import (
	"boilerplate/database"
	"boilerplate/models"
	prisma "boilerplate/prisma"
	"boilerplate/prisma/db"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/exp/slices"
)

// UserList returns a list of users
func UserList(c *fiber.Ctx) error {
	users := database.Get()

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

// UserCreate registers a user
func UserCreate(c *fiber.Ctx) error {
	user := &models.User{
		// Note: when writing to external database,
		// we can simply use - Name: c.FormValue("user")
		Name: utils.CopyString(c.FormValue("user")),
	}
	database.Insert(user)

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}

type Parent struct {
	Email string
}

func CreateParentHandler(c *fiber.Ctx) error {
	var parent Parent
	c.BodyParser(&parent)

	created, err := parent.Create()
	if err != nil {
		return err
	}
	return c.JSON(created)
}

func (Parent) Find(c *fiber.Ctx) error {
	parentId := c.Params("id")
	result, err := prisma.Client.Parent.FindUnique(db.Parent.ID.Equals(parentId)).With(db.Parent.Student.Fetch()).Exec(prisma.Ctx)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (Parent) FindAll(c *fiber.Ctx) error {
	result, err := prisma.Client.Parent.FindMany().With(db.Parent.Student.Fetch()).Exec(prisma.Ctx)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (Student) FindAll(c *fiber.Ctx) error {
	result, err := prisma.Client.Student.FindMany().Exec(prisma.Ctx)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (Payment) Find(c *fiber.Ctx) error {

	studentId := c.Query("studentId")
	semester := c.QueryInt("semester")
	year := c.Query("year")

	result, err := prisma.Client.Payment.FindUnique(
		db.Payment.StudentIDSemesterYear(
			db.Payment.StudentID.Equals(studentId),
			db.Payment.Semester.Equals(semester),
			db.Payment.Year.Equals(year),
		),
	).Exec(prisma.Ctx)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (p Parent) Create() (*db.ParentModel, error) {
	created, err := prisma.Client.Parent.CreateOne(
		db.Parent.Email.Set(p.Email),
	).Exec(prisma.Ctx)
	if err != nil {
		return nil, err
	}
	return created, nil
}

type Student struct{}

// CreateStudentHandler creates a new student associated with a parent
func CreateStudentHandler(c *fiber.Ctx) error {
	var student db.InnerStudent
	err := c.BodyParser(&student)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Create student
	created, err := prisma.Client.Student.CreateOne(
		db.Student.ID.Set(student.ID),
		db.Student.FirstName.Set(student.FirstName),
		db.Student.LastName.Set(student.LastName),
		db.Student.Parent.Link(db.Parent.ID.Equals(student.ParentID)),
		db.Student.GPA.SetIfPresent(student.GPA),
		db.Student.Credits.SetIfPresent(student.Credits),
		db.Student.Class.SetIfPresent(student.Class),
		db.Student.Email.SetIfPresent(student.Email),
		db.Student.Phone.SetIfPresent(student.Phone),
		db.Student.Address.SetIfPresent(student.Address),
	).Exec(prisma.Ctx)
	if err != nil {
		return err
	}

	return c.JSON(created)
}

type Payment struct{}

// CreatePaymentHandler creates a new Payment associated with a parent
func CreatePaymentHandler(c *fiber.Ctx) error {
	var payment db.InnerPayment
	err := c.BodyParser(&payment)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}

	// Create Payment
	created, err := prisma.Client.Payment.CreateOne(
		db.Payment.Student.Link(db.Student.ID.Equals(payment.StudentID)),
		db.Payment.Amount.Set(payment.Amount),
		db.Payment.Year.Set(payment.Year),
		db.Payment.Semester.Set(payment.Semester),
	).Exec(prisma.Ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"student": created,
	})
}

type PaymentMethod struct {
	Method string
}

func PaidHandler(c *fiber.Ctx) error {
	paymentId := c.Params("id")
	allowedMethod := []string{"qr", "creditcard", "installment"}

	var payload PaymentMethod
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	m, err := json.Marshal(allowedMethod)
	if err != nil {
		return err
	}

	if !slices.Contains(allowedMethod, payload.Method) {
		return c.Status(400).SendString("Payment method not allowed, allowed: " + string(m))
	}

	updated, err := prisma.Client.Payment.FindUnique(
		db.Payment.ID.Equals(paymentId),
	).Update(
		db.Payment.Paid.Set(true),
		db.Payment.PaymentDate.Set(time.Now()),
		db.Payment.PaymentMethod.Set(payload.Method),
	).Exec(prisma.Ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}

	return c.JSON(updated)
}
