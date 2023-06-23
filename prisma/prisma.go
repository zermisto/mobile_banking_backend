package prisma

import (
	"boilerplate/prisma/db"
	"context"
	"fmt"
)

var (
	Ctx    = context.Background()
	Client = db.NewClient()
)

func init() {
	err := Client.Prisma.Connect()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database successfully connected")
}
