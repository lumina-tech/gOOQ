package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	model "github.com/lumina-tech/gooq/examples/swapi/table"
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

	stmt := gooq.InsertInto(model.Person).
		Set(model.Person.ID, uuid.New()).
		Set(model.Person.Name, "Frank").
		Set(model.Person.EyeColor, "blue").
		Set(model.Person.BirthYear, "1998").
		Set(model.Person.Gender, "male").
		Set(model.Person.HairColor, "#000000").
		Set(model.Person.Height, "6'1").
		Set(model.Person.HomeWorld, "Runescape").
		Set(model.Person.Mass, "100").
		Set(model.Person.SkinColor, "#f0d7b9").
		Returning(model.Person.Asterisk)

	builder := gooq.Builder{}
	stmt.Render(&builder)
	fmt.Println(builder.String())

	p, err := model.Person.ScanRow(dockerDB.DB, stmt)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	fmt.Println(p)
}
