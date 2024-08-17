package models

import (
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/jackc/pgx/v5/pgtype"
)

type Location struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	NormalizedName     string             `json:"normalizedName"`
	Available          int64              `json:"available"`
	Capacity           int64              `json:"capacity"`
	AvailableYesterday int64              `json:"availableYesterday"`
	CapacityYesterday  int64              `json:"capacityYesterday"`
	YesterdayFullAt    pgtype.Timestamptz `json:"yesterdayFullAt"    swaggertype:"string"`
	TimeZone           string             `json:"timeZone"`
	UserID             string             `json:"userId"`
} //	@name	Location

func (location *Location) SetFields(
	checkInsToday []*CheckIn,
	checkInsYesterday []*CheckIn,
) error {
	location.SetCheckInRelatedFields(
		checkInsToday,
		checkInsYesterday,
	)

	return location.NormalizeName()
}

func (location *Location) SetCheckInRelatedFields(
	allCheckInsToday []*CheckIn,
	allCheckInsYesterday []*CheckIn,
) {
	checkInsToday := []*CheckIn{}
	for _, checkIn := range allCheckInsToday {
		if checkIn.LocationID == location.ID {
			checkInsToday = append(checkInsToday, checkIn)
		}
	}

	checkInsYesterday := []*CheckIn{}
	for _, checkIn := range allCheckInsYesterday {
		if checkIn.LocationID == location.ID {
			checkInsYesterday = append(checkInsYesterday, checkIn)
		}
	}

	location.Available = location.Capacity - int64(len(checkInsToday))
	location.CapacityYesterday = 0
	//nolint:exhaustruct //other fields are optional
	location.YesterdayFullAt = pgtype.Timestamptz{}

	var lastCheckInYesterday *CheckIn
	switch {
	case len(checkInsYesterday) > 0:
		lastCheckInYesterday = checkInsYesterday[len(checkInsYesterday)-1]

		location.CapacityYesterday = lastCheckInYesterday.Capacity
		location.AvailableYesterday = location.CapacityYesterday - int64(
			len(checkInsYesterday),
		)
	case len(checkInsToday) > 0:
		location.CapacityYesterday = checkInsToday[0].Capacity
		location.AvailableYesterday = checkInsToday[0].Capacity
	default:
		location.CapacityYesterday = location.Capacity
		location.AvailableYesterday = location.Capacity
	}

	if location.AvailableYesterday == 0 && lastCheckInYesterday != nil {
		location.YesterdayFullAt = lastCheckInYesterday.CreatedAt
	}
}

func (location *Location) NormalizeName() error {
	output, err := normalize(location.Name)
	if err != nil {
		return err
	}

	location.NormalizedName = *output

	return nil
}

func (location *Location) CompareNormalizedName(name string) (bool, error) {
	err := location.NormalizeName()
	if err != nil {
		return false, err
	}

	normalizedName, err := normalize(name)
	if err != nil {
		return false, err
	}

	if location.NormalizedName != *normalizedName {
		return false, nil
	}

	return true, nil
}

func normalize(str string) (*string, error) {
	re1 := regexp2.MustCompile(`\s`, 0)
	re2 := regexp2.MustCompile(`^-+|[^a-z0-9-]|(?<!-)-+$`, 0)

	lower := strings.ToLower(str)
	re1Result, err := re1.Replace(lower, "-", -1, -1)
	if err != nil {
		return nil, err
	}

	output, err := re2.Replace(re1Result, "", -1, -1)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
