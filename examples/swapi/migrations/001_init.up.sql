CREATE TABLE color_reference_table(
  value text primary key NOT NULL
);

INSERT INTO color_reference_table (value) VALUES
  ('black'),
  ('brown'),
  ('red'),
  ('orange'),
  ('yellow'),
  ('green'),
  ('blue'),
  ('purple');

CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE person(
  id uuid primary key NOT NULL,
  name text NOT NULL,
  height decimal NOT NULL,
  mass decimal NOT NULL,
  hair_color text NOT NULL,
  skin_color text NOT NULL,
  eye_color text NOT NULL,
  birth_year int NOT NULL,
  gender gender NOT NULL,
  home_world text NOT NULL,
  FOREIGN KEY (hair_color) REFERENCES color_reference_table(value),
  FOREIGN KEY (skin_color) REFERENCES color_reference_table(value),
  FOREIGN KEY (eye_color) REFERENCES color_reference_table(value)
);

-- type Person struct {
-- 	Name      string   `json:"name"`
-- 	Height    string   `json:"height"`
-- 	Mass      string   `json:"mass"`
-- 	HairColor string   `json:"hair_color"`
-- 	SkinColor string   `json:"skin_color"`
-- 	EyeColor  string   `json:"eye_color"`
-- 	BirthYear string   `json:"birth_year"`
-- 	Gender    string   `json:"gender"`
-- 	Homeworld string   `json:"homeworld"`
-- 	Films     []string `json:"films"`
-- 	Species   []string `json:"species"`
-- 	Vehicles  []string `json:"vehicles"`
-- 	Starships []string `json:"starships"`
-- 	Created   string   `json:"created"`
-- 	Edited    string   `json:"edited"`
-- 	URL       string   `json:"url"`
-- }
--
-- type Film struct {
-- 	Title        string   `json:"title"`
-- 	EpisodeID    int64    `json:"episode_id"`
-- 	OpeningCrawl string   `json:"opening_crawl"`
-- 	Director     string   `json:"director"`
-- 	Producer     string   `json:"producer"`
-- 	Characters   []string `json:"characters"`
-- 	Planets      []string `json:"planets"`
-- 	Starships    []string `json:"starships"`
-- 	Vehicles     []string `json:"vehicles"`
-- 	Species      []string `json:"species"`
-- 	Created      string   `json:"created"`
-- 	Edited       string   `json:"edited"`
-- 	URL          string   `json:"url"`
-- }
--
-- type Planet struct {
-- 	Name           string   `json:"name"`
-- 	RotationPeriod string   `json:"rotation_period"`
-- 	OrbitalPeriod  string   `json:"orbital_period"`
-- 	Diameter       string   `json:"diameter"`
-- 	Climate        string   `json:"climate"`
-- 	Gravity        string   `json:"gravity"`
-- 	Terrain        string   `json:"terrain"`
-- 	SurfaceWater   string   `json:"surface_water"`
-- 	Population     string   `json:"population"`
-- 	Residents      []string `json:"residents"`
-- 	Films          []string `json:"films"`
-- 	Created        string   `json:"created"`
-- 	Edited         string   `json:"edited"`
-- 	URL            string   `json:"url"`
-- }
--
-- type Species struct {
-- 	Name            string   `json:"name"`
-- 	Classification  string   `json:"classification"`
-- 	Designation     string   `json:"designation"`
-- 	AverageHeight   string   `json:"average_height"`
-- 	SkinColors      string   `json:"skin_colors"`
-- 	HairColors      string   `json:"hair_colors"`
-- 	EyeColors       string   `json:"eye_colors"`
-- 	AverageLifespan string   `json:"average_lifespan"`
-- 	Homeworld       string   `json:"homeworld"`
-- 	Language        string   `json:"language"`
-- 	People          []string `json:"people"`
-- 	Films           []string `json:"films"`
-- 	Created         string   `json:"created"`
-- 	Edited          string   `json:"edited"`
-- 	URL             string   `json:"url"`
-- }
--
-- type Starship struct {
-- 	Name                 string   `json:"name"`
-- 	Model                string   `json:"model"`
-- 	Manufacturer         string   `json:"manufacturer"`
-- 	CostInCredits        string   `json:"cost_in_credits"`
-- 	Length               string   `json:"length"`
-- 	MaxAtmospheringSpeed string   `json:"max_atmosphering_speed"`
-- 	Crew                 string   `json:"crew"`
-- 	Passengers           string   `json:"passengers"`
-- 	CargoCapacity        string   `json:"cargo_capacity"`
-- 	Consumables          string   `json:"consumables"`
-- 	HyperdriveRating     string   `json:"hyperdrive_rating"`
-- 	MGLT                 string   `json:"MGLT"`
-- 	StarshipClass        string   `json:"starship_class"`
-- 	Pilots               []string `json:"pilots"`
-- 	Films                []string `json:"films"`
-- 	Created              string   `json:"created"`
-- 	Edited               string   `json:"edited"`
-- 	URL                  string   `json:"url"`
-- }
--
-- type Vehicle struct {
-- 	Name                 string   `json:"name"`
-- 	Model                string   `json:"model"`
-- 	Manufacturer         string   `json:"manufacturer"`
-- 	CostInCredits        string   `json:"cost_in_credits"`
-- 	Length               string   `json:"length"`
-- 	MaxAtmospheringSpeed string   `json:"max_atmosphering_speed"`
-- 	Crew                 string   `json:"crew"`
-- 	Passengers           string   `json:"passengers"`
-- 	CargoCapacity        string   `json:"cargo_capacity"`
-- 	Consumables          string   `json:"consumables"`
-- 	VehicleClass         string   `json:"vechicle_class"`
-- 	Pilots               []string `json:"pilots"`
-- 	Films                []string `json:"films"`
-- 	Created              string   `json:"created"`
-- 	Edited               string   `json:"edited"`
-- 	URL                  string   `json:"url"`
-- }

