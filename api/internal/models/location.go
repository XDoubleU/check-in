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

func (location *Location) SetCheckInRelatedFields(
	checkInsToday []*CheckIn,
	checkInsYesterday []*CheckIn,
	lastCheckInYesterday *CheckIn,
) {
	location.Available = location.Capacity - int64(len(checkInsToday))
	location.CapacityYesterday = 0

	if lastCheckInYesterday != nil {
		location.CapacityYesterday = lastCheckInYesterday.Capacity
	}

	location.AvailableYesterday = location.CapacityYesterday - int64(
		len(checkInsYesterday),
	)
	location.YesterdayFullAt = pgtype.Timestamptz{}

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
