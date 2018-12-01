package station

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// Stations [station id] -> station
var stations []Station

// Lines [line id] -> line
var lines []Line

// Joints [jointId id] joint
var joints []Joint

// StationLine [StationLineId Id] StationLine
var stationLine []StationLine

func LoadDateFromDB(db *sqlx.DB) {
	if err := db.Select(&stations, "SELECT * FROM Station"); err != nil {
		log.Fatalf("Cannot read stations %+v", err)
	}
	if err := db.Select(&joints, "SELECT * FROM Joints"); err != nil {
		log.Fatalf("Cannot read joints %+v", err)
	}
	if err := db.Select(&stationLine, "SELECT * FROM StationLine"); err != nil {
		log.Fatalf("Cannot read stationLines %+v", err)
	}
	if err := db.Select(&lines, "SELECT * FROM Lines"); err != nil {
		log.Fatalf("Cannot read lines %+v", err)
	}
}

// Getters

func GetStations() []Station {
	cp := make([]Station, len(stations))
	copy(cp, stations)
	return cp
}

func GetStation(stationId int8) Station {
	return stations[stationId]
}

func GetLine(lineId int8) Line {
	return lines[lineId]
}

func GetJoint(jointId int8) Joint {
	return joints[jointId]
}

func GetStationLine(stationLineId int8) StationLine {
	return stationLine[stationLineId]
}

// Utils

func GetStationById(stationId int8) *Station {
	for _, v := range stations {
		if v.ID == stationId {
			return &v
		}
	}
	return nil
}

func GetStationIdByName(stationName string) int8 {
	for _, v := range stations {
		if v.Name == stationName {
			return v.ID
		}
	}
	return -1
}

func GetJoinnedStationsToExcluding(stationId int8, excludeToId int8) []*Station {
	result := make([]*Station, 0)

	for _, v := range joints {
		if v.FromStation == stationId && v.ToStation != excludeToId {
			s := GetStationById(v.ToStation)
			result = append(result, s)
		}
	}

	return result
}

func goFromTo(last int8, from int8, to int8, stations []Station) []Station {

	connectedStations := GetJoinnedStationsToExcluding(from, last)

	station := *GetStationById(from)
	station.isTransit = len(connectedStations) != 1
	stations = append(stations, station)

	if from == to {
		return stations
	}

	for _, v := range connectedStations {
		stationsCopy := make([]Station, len(stations))
		copy(stationsCopy, stations)
		res := goFromTo(from, v.ID, to, stationsCopy)
		if res[len(res)-1].ID == to {
			return res
		}
	}

	return stations
}

func GoFromTo(from int8, to int8) []Station {

	back := make([]Station, 0)
	return goFromTo(-1, from, to, back)

}
