Using sqlx to map nested query

```
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
```
The following sql statement is generated
```
SELECT
  person.*, 
  species.name AS "species.name", 
  species.classification AS "species.name", 
  species.average_height AS "species.average_height", 
  species.average_lifespan AS "species.average_lifespan", 
  species.hair_color AS "species.hair_color", 
  species.skin_color AS "species.skin_color", 
  species.eye_color AS "species.eye_color", 
  species.home_world AS "species.home_world", 
  species.language AS "species.language" 
FROM public.person 
JOIN public.species ON person.species_id = species.id
```
When results is mashalled as a JSON
```
[
  {
    "id": "534df3d9-8239-4cad-a6bf-7cd9fe3c82be",
    "name": "Frank",
    "height": 170.3,
    "mass": 150.5,
    "hair_color": "black",
    "skin_color": "orange",
    "eye_color": "brown",
    "birth_year": 1998,
    "gender": "male",
    "home_world": "Runescape",
    "species_id": "d0269559-7772-470b-a5df-c67ad59c68dc",
    "species": {
      "id": "00000000-0000-0000-0000-000000000000",
      "name": "Mammal",
      "classification": "",
      "average_height": 160.5,
      "average_lifespan": 70,
      "hair_color": "black",
      "skin_color": "orange",
      "eye_color": "brown",
      "home_world": "Earth",
      "language": "English"
    }
  }
]
```