package services

import (
	"context"
	"time"

	"check-in/api/internal/database"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
)

type LocationService struct {
	db database.DB
}

func (service LocationService) GetCheckInsEntriesDay(
	checkIns []*models.CheckIn,
	schools []*models.School,
) map[int64]*dtos.CheckInsLocationEntryRaw {
	schoolsIDNameMap, _ := getSchoolMaps(schools)

	checkInEntries := make(map[int64]*dtos.CheckInsLocationEntryRaw)

	var lastEntry *dtos.CheckInsLocationEntryRaw
	for _, checkIn := range checkIns {
		schoolName := schoolsIDNameMap[checkIn.SchoolID]

		_, schoolsMap := getSchoolMaps(schools)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacity: checkIn.Capacity,
			Schools:  schoolsMap,
		}

		checkInEntry.Schools[schoolName]++

		if lastEntry != nil {
			checkInEntry.Schools[schoolName] += lastEntry.Schools[schoolName]
		}

		checkInEntries[checkIn.CreatedAt.Unix()] = checkInEntry
		lastEntry = checkInEntries[checkIn.CreatedAt.Unix()]
	}

	return checkInEntries
}

func (service LocationService) GetCheckInsEntriesRange(
	startDate time.Time,
	endDate time.Time,
	checkIns []*models.CheckIn,
	schools []*models.School,
) map[int64]*dtos.CheckInsLocationEntryRaw {
	schoolsIDNameMap, _ := getSchoolMaps(schools)

	checkInEntries := make(map[int64]*dtos.CheckInsLocationEntryRaw)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dVal := helpers.StartOfDay(&d)

		_, schoolsMap := getSchoolMaps(schools)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacity: 0,
			Schools:  schoolsMap,
		}

		checkInEntries[dVal.Unix()] = checkInEntry
	}

	for _, checkIn := range checkIns {
		datetime := helpers.StartOfDay(&checkIn.CreatedAt)
		schoolName := schoolsIDNameMap[checkIn.SchoolID]

		checkInEntry := checkInEntries[datetime.Unix()]

		checkInEntry.Schools[schoolName]++

		if checkIn.Capacity > checkInEntry.Capacity {
			checkInEntry.Capacity = checkIn.Capacity
		}
	}

	return checkInEntries
}

func getSchoolMaps(
	schools []*models.School,
) (map[int64]string, map[string]int) {
	schoolsIDNameMap := make(map[int64]string)
	schoolsMap := make(map[string]int)
	for _, school := range schools {
		schoolsIDNameMap[school.ID] = school.Name
		schoolsMap[school.Name] = 0
	}

	return schoolsIDNameMap, schoolsMap
}

func (service LocationService) GetTotalCount(ctx context.Context) (*int64, error) {
	query := `
		SELECT COUNT(*)
		FROM locations
	`

	var total *int64

	err := service.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return nil, handleError(err)
	}

	return total, nil
}

func (service LocationService) GetAll(ctx context.Context) ([]*models.Location, error) {
	query := `
		SELECT id, name, capacity
		FROM locations
	`

	rows, err := service.db.Query(ctx, query)
	if err != nil {
		return nil, handleError(err)
	}

	locations := []*models.Location{}

	for rows.Next() {
		var location models.Location

		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.Capacity,
		)

		if err != nil {
			return nil, handleError(err)
		}

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError(err)
	}

	return locations, nil
}

func (service LocationService) GetAllPaginated(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]*models.Location, error) {
	query := `
		SELECT id, name, capacity
		FROM locations
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := service.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, handleError(err)
	}

	locations := []*models.Location{}

	for rows.Next() {
		var location models.Location

		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.Capacity,
		)

		if err != nil {
			return nil, handleError(err)
		}

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError(err)
	}

	return locations, nil
}

func (service LocationService) GetByID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	query := `
		SELECT name, capacity, user_id,
		(capacity - (SELECT COUNT(*) 
								FROM check_ins 
								WHERE DATE(check_ins.created_at) = DATE(NOW()) 
								AND check_ins.location_id = locations.id))
		FROM locations
		WHERE locations.id = $1
	`

	location := models.Location{
		ID: id,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(&location.Name, &location.Capacity, &location.UserID, &location.Available)

	if err != nil {
		return nil, handleError(err)
	}

	err = location.NormalizeName()
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (service LocationService) GetByUserID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	query := `
		SELECT locations.id, locations.name, locations.capacity, 
		(capacity - (SELECT COUNT(*) 
								FROM check_ins 
								WHERE DATE(check_ins.created_at) = DATE(NOW()) 
								AND check_ins.location_id = locations.id)),
		(SELECT MAX(check_ins.created_at) 
				FROM check_ins
				INNER JOIN (
					SELECT location_id, COUNT(*) AS total_check_ins, MAX(capacity) AS max_capacity
					FROM check_ins
					WHERE DATE(created_at) = (DATE(NOW()) - INTERVAL '1' DAY)
					GROUP BY location_id
				) daily_stats ON check_ins.location_id = daily_stats.location_id
				WHERE check_ins.location_id = locations.id
				AND DATE(check_ins.created_at) = (DATE(NOW()) - INTERVAL '1' DAY)
				AND daily_stats.total_check_ins >= daily_stats.max_capacity
			)
		FROM locations
		WHERE locations.user_id = $1
	`

	location := models.Location{
		UserID: id,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		id).Scan(
		&location.ID,
		&location.Name,
		&location.Capacity,
		&location.Available,
		&location.YesterdayFullAt,
	)

	if err != nil {
		return nil, handleError(err)
	}

	err = location.NormalizeName()
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (service LocationService) GetByName(
	ctx context.Context,
	name string,
) (*models.Location, error) {
	locations, err := service.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	for _, location := range locations {
		var output bool

		output, err = location.CompareNormalizedName(name)
		if err != nil {
			return nil, handleError(err)
		}

		if output {
			return location, nil
		}
	}

	return nil, ErrRecordNotFound
}

func (service LocationService) Create(
	ctx context.Context,
	name string,
	capacity int64,
	userID string,
) (*models.Location, error) {
	query := `
		INSERT INTO locations (name, capacity, user_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	location := models.Location{
		Name:      name,
		Capacity:  capacity,
		Available: capacity,
		UserID:    userID,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		name,
		capacity,
		userID,
	).Scan(&location.ID)

	if err != nil {
		return nil, handleError(err)
	}

	err = location.NormalizeName()
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (service LocationService) Update(
	ctx context.Context,
	location *models.Location,
	user *models.User,
	updateLocationDto dtos.UpdateLocationDto,
) error {
	if updateLocationDto.Name != nil {
		location.Name = *updateLocationDto.Name
	}

	if updateLocationDto.Capacity != nil {
		diff := *updateLocationDto.Capacity - location.Capacity
		location.Available += diff

		if location.Available < 0 {
			location.Available = 0
		}

		location.Capacity = *updateLocationDto.Capacity
	}

	if updateLocationDto.Username != nil {
		user.Username = *updateLocationDto.Username
	}

	if updateLocationDto.Password != nil {
		passwordHash, _ := models.HashPassword(*updateLocationDto.Password)
		user.PasswordHash = passwordHash
	}

	tx, err := service.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return handleError(err)
	}

	queryLocation := `
		UPDATE locations
		SET name = $2, capacity = $3
		WHERE id = $1
	`

	resultLocation, err := tx.Exec(
		ctx,
		queryLocation,
		location.ID,
		location.Name,
		location.Capacity,
	)

	if err != nil {
		return handleError(err)
	}

	rowsAffected := resultLocation.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	queryUser := `
		UPDATE users
		SET username = $2, password_hash = $3
		WHERE id = $1 AND role = 'default'
	`

	resultUser, err := tx.Exec(
		ctx,
		queryUser,
		user.ID,
		user.Username,
		user.PasswordHash,
	)

	if err != nil {
		return handleError(err)
	}

	rowsAffected = resultUser.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	err = location.NormalizeName()
	if err != nil {
		return err
	}

	return nil
}

func (service LocationService) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM locations
		WHERE id = $1
	`

	result, err := service.db.Exec(ctx, query, id)
	if err != nil {
		return handleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
