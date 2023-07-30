package dtos

import "check-in/api/internal/validator"

type CheckInDto struct {
	SchoolID int64  `json:"schoolId"`
	TimeZone string `json:"timeZone"`
} //	@name	CheckInDto

func ValidateCheckInDto(v *validator.Validator, checkInDto CheckInDto) {
	v.Check(checkInDto.SchoolID > 0, "schoolId", "must be greater than zero")
	v.Check(len(checkInDto.TimeZone) > 0, "timeZone", "must be provided")
}
