package main

import (
	"fmt"
	"os"

	"github.com/lumina-tech/gooq/examples/swapi/model"
	"github.com/lumina-tech/gooq/examples/swapi/table"

	"github.com/google/uuid"
	"github.com/lumina-tech/gooq/pkg/database"
	"github.com/lumina-tech/gooq/pkg/gooq"
)

func main() {
	dockerDB := database.NewDockerizedDB(&database.DatabaseConfig{
		Host:         "localhost",
		Port:         5432,
		Username:     "postgres",
		Password:     "password",
		DatabaseName: "swapi",
		SSLMode:      "disable",
	}, "11.4-alpine")
	defer dockerDB.Close()

	database.MigrateDatabase(dockerDB.DB.DB, "migrations")

	stmt := gooq.InsertInto(table.Person).
		Set(table.Person.ID, uuid.New()).
		Set(table.Person.Name, "Frank").
		Set(table.Person.BirthYear, "1998").
		Set(table.Person.Height, "6'1").
		Set(table.Person.HomeWorld, "Runescape").
		Set(table.Person.Gender, model.GenderMale).
		Set(table.Person.EyeColor, model.ColorBrown).
		Set(table.Person.HairColor, model.ColorBlack).
		Set(table.Person.SkinColor, model.ColorOrange).
		Set(table.Person.Mass, "100").
		Returning(table.Person.Asterisk)

	builder := gooq.Builder{}
	stmt.Render(&builder)
	fmt.Println(builder.String())

	p, err := table.Person.ScanRow(dockerDB.DB, stmt)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	fmt.Println(p)
}
