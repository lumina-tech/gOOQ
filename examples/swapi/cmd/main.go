package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/lumina-tech/gooq/examples/swapi/model"
	"github.com/lumina-tech/gooq/examples/swapi/table"
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

	speciesStmt := gooq.InsertInto(table.Species).
		Set(table.Species.ID, uuid.New()).
		Set(table.Species.Name, "Human").
		Set(table.Species.Classification, "Mammal").
		Set(table.Species.AverageHeight, 160.5).
		Set(table.Species.AverageLifespan, 70).
		Set(table.Species.HairColor, model.ColorBlack).
		Set(table.Species.SkinColor, model.ColorOrange).
		Set(table.Species.EyeColor, model.ColorBrown).
		Set(table.Species.HomeWorld, "Earth").
		Set(table.Species.Language, "English").
		Returning(table.Species.Asterisk)
	species, err := table.Species.ScanRow(dockerDB.DB, speciesStmt)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	personStmt := gooq.InsertInto(table.Person).
		Set(table.Person.ID, uuid.New()).
		Set(table.Person.Name, "Frank").
		Set(table.Person.Height, 170.3).
		Set(table.Person.Mass, 150.5).
		Set(table.Person.BirthYear, 1998).
		Set(table.Person.HomeWorld, "Runescape").
		Set(table.Person.Gender, model.GenderMale).
		Set(table.Person.EyeColor, model.ColorBrown).
		Set(table.Person.HairColor, model.ColorBlack).
		Set(table.Person.SkinColor, model.ColorOrange).
		Set(table.Person.SpeciesID, species.ID).
		Returning(table.Person.Asterisk)
	_, err = table.Person.ScanRow(dockerDB.DB, personStmt)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	type PersonWithSpecies struct {
		model.Person
		Species *model.Species `db:"species"`
	}
	stmt := gooq.Select(
		table.Person.Asterisk,
		table.Species.Name.As("species.name"),
		table.Species.Classification.As("species.name"),
		table.Species.AverageHeight.As("species.average_height"),
		table.Species.AverageLifespan.As("species.average_lifespan"),
		table.Species.HairColor.As("species.hair_color"),
		table.Species.SkinColor.As("species.skin_color"),
		table.Species.EyeColor.As("species.eye_color"),
		table.Species.HomeWorld.As("species.home_world"),
		table.Species.Language.As("species.language"),
	).From(table.Person).
		Join(table.Species).
		On(table.Person.SpeciesID.Eq(table.Species.ID))

	builder := &gooq.Builder{}
	stmt.Render(builder)
	fmt.Println(builder.String())

	var results []PersonWithSpecies
	if err := gooq.ScanRows(dockerDB.DB, stmt, &results); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	bytes, _ := json.Marshal(results)
	fmt.Println(string(bytes))
}
