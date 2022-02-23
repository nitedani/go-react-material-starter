package db

import "context"

var Client = NewClient()
var Ctx = context.Background()

func Connect() {
	if err := Client.Prisma.Connect(); err != nil {
		panic(err)
	}
}

func Disconnect() {
	if err := Client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}
