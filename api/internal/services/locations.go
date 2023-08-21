package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"check-in/api/internal/database"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
)

type LocationService struct {
	db      database.DB
	schools SchoolService
}

const availableQuery = `
	(	capacity - (	
			SELECT COUNT(*) 
			FROM check_ins 
			WHERE DATE(check_ins.created_at AT TIME ZONE locations.time_zone) 
			= DATE(NOW() AT TIME ZONE locations.time_zone) 
			AND check_ins.location_id = locations.id
		)
	)
`

const yesterdayFullAtQuery = `
	(	SELECT MAX(check_ins.created_at AT TIME ZONE locations.time_zone)
		FROM check_ins
		INNER JOIN (
			SELECT location_id, COUNT(*) AS total_check_ins, MAX(capacity) AS max_capacity
			FROM check_ins
			WHERE DATE(created_at AT TIME ZONE locations.time_zone) 
			= (DATE(NOW() AT TIME ZONE locations.time_zone) - INTERVAL '1' DAY)
			GROUP BY location_id
		) daily_stats 
		ON check_ins.location_id = daily_stats.location_id
		WHERE check_ins.location_id = locations.id
		AND DATE(check_ins.created_at AT TIME ZONE locations.time_zone) 
		= (DATE(NOW() AT TIME ZONE locations.time_zone) - INTERVAL '1' DAY)
		AND daily_stats.total_check_ins >= daily_stats.max_capacity
	)
`

func (service LocationService) GetCheckInsEntriesDay(
	checkIns []*models.CheckIn,
	schools []*models.School,
) *orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw] {
	schoolsIDNameMap, _ := service.schools.GetSchoolMaps(schools)

	checkInEntries := orderedmap.New[string, *dtos.CheckInsLocationEntryRaw]()

	_, lastEntrySchoolsMap := service.schools.GetSchoolMaps(schools)
	for _, checkIn := range checkIns {
		schoolName := schoolsIDNameMap[checkIn.SchoolID]

		var schoolsMap dtos.SchoolsMap

		data, _ := json.Marshal(lastEntrySchoolsMap)
		_ = json.Unmarshal(data, &schoolsMap)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacity: checkIn.Capacity,
			Schools:  schoolsMap,
		}

		schoolValue, _ := checkInEntry.Schools.Get(schoolName)
		schoolValue++
		checkInEntry.Schools.Set(schoolName, schoolValue)

		checkInEntries.Set(
			checkIn.CreatedAt.Time.Format(time.RFC3339),
			checkInEntry,
		)
		lastEntrySchoolsMap = checkInEntry.Schools
	}

	return checkInEntries
}

func (service LocationService) GetCheckInsEntriesRange(
	startDate *time.Time,
	endDate *time.Time,
	checkIns []*models.CheckIn,
	schools []*models.School,
) *orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw] {
	schoolsIDNameMap, _ := service.schools.GetSchoolMaps(schools)

	checkInEntries := orderedmap.New[string, *dtos.CheckInsLocationEntryRaw]()
	for d := *startDate; !d.After(*endDate); d = d.AddDate(0, 0, 1) {
		dVal := helpers.StartOfDay(&d)

		_, schoolsMap := service.schools.GetSchoolMaps(schools)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacity: 0,
			Schools:  schoolsMap,
		}

		checkInEntries.Set(dVal.Format(time.RFC3339), checkInEntry)
	}

	for _, checkIn := range checkIns {
		datetime := helpers.StartOfDay(&checkIn.CreatedAt.Time)
		schoolName := schoolsIDNameMap[checkIn.SchoolID]

		checkInEntry, _ := checkInEntries.Get(datetime.Format(time.RFC3339))

		schoolValue, _ := checkInEntry.Schools.Get(schoolName)
		schoolValue++
		checkInEntry.Schools.Set(schoolName, schoolValue)

		if checkIn.Capacity > checkInEntry.Capacity {
			checkInEntry.Capacity = checkIn.Capacity
		}
	}

	return checkInEntries
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
		SELECT id, name, capacity, %s, %s
		FROM locations
	`

	query = fmt.Sprintf(query, availableQuery, yesterdayFullAtQuery)

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
			&location.Available,
			&location.YesterdayFullAt,
		)

		if err != nil {
			return nil, handleError(err)
		}

		err = location.NormalizeName()
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
		SELECT id, name, capacity, user_id, time_zone, %s, %s
		FROM locations
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	query = fmt.Sprintf(query, availableQuery, yesterdayFullAtQuery)

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
			&location.UserID,
			&location.TimeZone,
			&location.Available,
			&location.YesterdayFullAt,
		)

		if err != nil {
			return nil, handleError(err)
		}

		err = location.NormalizeName()
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
		SELECT name, capacity, user_id, time_zone, %s, %s
		FROM locations
		WHERE locations.id = $1
	`

	query = fmt.Sprintf(query, availableQuery, yesterdayFullAtQuery)

	location := models.Location{
		ID: id,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&location.Name,
		&location.Capacity,
		&location.UserID,
		&location.TimeZone,
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

func (service LocationService) GetByUserID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, %s, %s
		FROM locations
		WHERE user_id = $1
	`

	query = fmt.Sprintf(query, availableQuery, yesterdayFullAtQuery)

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
		&location.TimeZone,
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
	timeZone string,
	username string,
	password string,
) (*models.Location, error) {
	tx, err := service.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return nil, handleError(err)
	}

	var user *models.User
	var location *models.Location

	user, err = createUser(ctx, tx, username, password)
	if err != nil {
		return nil, err
	}

	location, err = createLocation(ctx, tx, name, capacity, timeZone, user.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func createUser(
	ctx context.Context,
	tx pgx.Tx,
	username string,
	password string,
) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, 'default')
		RETURNING id
	`

	user := models.User{
		Username: username,
		Role:     models.DefaultRole,
	}

	passwordHash, err := models.HashPassword(password)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		ctx,
		query,
		username,
		passwordHash,
	).Scan(&user.ID)

	if err != nil {
		return nil, handleError(err)
	}

	return &user, nil
}

func createLocation(
	ctx context.Context,
	tx pgx.Tx,
	name string,
	capacity int64,
	timeZone string,
	userID string,
) (*models.Location, error) {
	query := `
		INSERT INTO locations (name, capacity, time_zone, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	location := models.Location{
		Name:      name,
		Capacity:  capacity,
		Available: capacity,
		TimeZone:  timeZone,
		UserID:    userID,
	}

	err := tx.QueryRow(
		ctx,
		query,
		name,
		capacity,
		timeZone,
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
	locationChanged := false
	userChanged := false

	if updateLocationDto.Name != nil {
		locationChanged = true

		location.Name = *updateLocationDto.Name
	}

	if updateLocationDto.Capacity != nil {
		locationChanged = true

		diff := *updateLocationDto.Capacity - location.Capacity
		location.Available += diff

		if location.Available < 0 {
			location.Available = 0
		}

		location.Capacity = *updateLocationDto.Capacity
	}

	if updateLocationDto.TimeZone != nil {
		locationChanged = true

		location.TimeZone = *updateLocationDto.TimeZone
	}

	if updateLocationDto.Username != nil {
		userChanged = true

		user.Username = *updateLocationDto.Username
	}

	if updateLocationDto.Password != nil {
		userChanged = true

		passwordHash, _ := models.HashPassword(*updateLocationDto.Password)
		user.PasswordHash = passwordHash
	}

	tx, err := service.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return handleError(err)
	}

	if locationChanged {
		err = updateLocation(ctx, tx, location)
		if err != nil {
			return err
		}
	}

	if userChanged {
		err = updateUser(ctx, tx, user)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func updateLocation(ctx context.Context, tx pgx.Tx, location *models.Location) error {
	queryLocation := `
			UPDATE locations
			SET name = $2, capacity = $3, time_zone = $4
			WHERE id = $1
		`

	resultLocation, err := tx.Exec(
		ctx,
		queryLocation,
		location.ID,
		location.Name,
		location.Capacity,
		location.TimeZone,
	)

	if err != nil {
		return handleError(err)
	}

	rowsAffected := resultLocation.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	err = location.NormalizeName()
	if err != nil {
		return err
	}

	return nil
}

func updateUser(ctx context.Context, tx pgx.Tx, user *models.User) error {
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

	rowsAffected := resultUser.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (service LocationService) Delete(
	ctx context.Context,
	location *models.Location,
) error {
	tx, err := service.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return handleError(err)
	}

	err = deleteLocation(ctx, tx, location.ID)
	if err != nil {
		return err
	}

	err = deleteUser(ctx, tx, location.UserID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func deleteLocation(ctx context.Context, tx pgx.Tx, id string) error {
	query := `
		DELETE FROM locations
		WHERE id = $1
	`

	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return handleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func deleteUser(ctx context.Context, tx pgx.Tx, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1 AND role = 'default'
	`

	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return handleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
