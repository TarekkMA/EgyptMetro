package station

type Station struct {
	ID        int8     `db:"id"`
	Name      string   `db:"name"`
	Lat       *float64 `db:"lat"`
	Lng       *float64 `db:"lng"`
	isTransit bool
}

type Line struct {
	ID   int8   `db:"id"`
	Name string `db:"name"`
}

type Joint struct {
	FromStation int8  `db:"station_from"`
	ToStation   int8  `db:"station_to"`
	Time        int64 `db:"time_sec"`
}

type StationLine struct {
	StationID int8 `db:"station_id"`
	LineID    int8 `db:"line_id"`
}
