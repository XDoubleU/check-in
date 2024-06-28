package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/tools"
	"github.com/jackc/pgx/v5"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

type LocationRepository struct {
	db       postgres.DB
	schools  SchoolRepository
	checkins CheckInRepository
}

func (repo LocationRepository) GetCheckInsEntriesDay(
	checkIns []*models.CheckIn,
	schools []*models.School,
) *orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw] {
	schoolsIDNameMap, _ := repo.schools.GetSchoolMaps(schools)

	checkInEntries := orderedmap.New[string, *dtos.CheckInsLocationEntryRaw]()

	_, lastEntrySchoolsMap := repo.schools.GetSchoolMaps(schools)
	capacities := orderedmap.New[string, int64]()
	for _, checkIn := range checkIns {
		schoolName := schoolsIDNameMap[checkIn.SchoolID]

		// Used to deep copy schoolsMap
		var schoolsMap dtos.SchoolsMap
		data, _ := json.Marshal(lastEntrySchoolsMap)
		_ = json.Unmarshal(data, &schoolsMap)

		capacities.Set(checkIn.LocationID, checkIn.Capacity)

		// Used to deep copy capacities
		var capacitiesCopy *orderedmap.OrderedMap[string, int64]
		data, _ = json.Marshal(capacities)
		_ = json.Unmarshal(data, &capacitiesCopy)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacities: capacitiesCopy,
			Schools:    schoolsMap,
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

func (repo LocationRepository) GetCheckInsEntriesRange(
	startDate time.Time,
	endDate time.Time,
	checkIns []*models.CheckIn,
	schools []*models.School,
) *orderedmap.OrderedMap[string, *dtos.CheckInsLocationEntryRaw] {
	schoolsIDNameMap, _ := repo.schools.GetSchoolMaps(schools)

	checkInEntries := orderedmap.New[string, *dtos.CheckInsLocationEntryRaw]()
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dVal := tools.StartOfDay(d)

		_, schoolsMap := repo.schools.GetSchoolMaps(schools)

		checkInEntry := &dtos.CheckInsLocationEntryRaw{
			Capacities: orderedmap.New[string, int64](),
			Schools:    schoolsMap,
		}

		checkInEntries.Set(dVal.Format(time.RFC3339), checkInEntry)
	}

	for i := range checkIns {
		datetime := tools.StartOfDay(checkIns[i].CreatedAt.Time)
		schoolName := schoolsIDNameMap[checkIns[i].SchoolID]

		checkInEntry, _ := checkInEntries.Get(datetime.Format(time.RFC3339))

		schoolValue, _ := checkInEntry.Schools.Get(schoolName)
		schoolValue++
		checkInEntry.Schools.Set(schoolName, schoolValue)

		capacity, present := checkInEntry.Capacities.Get(checkIns[i].LocationID)
		if !present {
			capacity = 0
		}

		if checkIns[i].Capacity > capacity {
			capacity = checkIns[i].Capacity
		}

		checkInEntry.Capacities.Set(checkIns[i].LocationID, capacity)
	}

	return checkInEntries
}

func (repo LocationRepository) GetTotalCount(ctx context.Context) (*int64, error) {
	query := `
		SELECT COUNT(*)
		FROM locations
	`

	var total *int64

	err := repo.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return nil, postgres.HandleError(err)
	}

	return total, nil
}

func (repo LocationRepository) GetAll(ctx context.Context) ([]*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		ORDER BY name ASC
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.HandleError(err)
	}

	locations := []*models.Location{}
	for rows.Next() {
		var location models.Location

		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.Capacity,
			&location.TimeZone,
			&location.UserID,
		)
		if err != nil {
			return nil, postgres.HandleError(err)
		}

		err = location.NormalizeName()
		if err != nil {
			return nil, postgres.HandleError(err)
		}

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.HandleError(err)
	}

	for i, location := range locations {
		err = repo.prepareLocation(ctx, location)
		if err != nil {
			return nil, postgres.HandleError(err)
		}

		locations[i] = location
	}

	return locations, nil
}

func (repo LocationRepository) GetAllPaginated(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := repo.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, postgres.HandleError(err)
	}

	locations := []*models.Location{}
	for rows.Next() {
		var location models.Location

		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.Capacity,
			&location.TimeZone,
			&location.UserID,
		)
		if err != nil {
			return nil, postgres.HandleError(err)
		}

		err = location.NormalizeName()
		if err != nil {
			return nil, postgres.HandleError(err)
		}

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.HandleError(err)
	}

	for i, location := range locations {
		err = repo.prepareLocation(ctx, location)
		if err != nil {
			return nil, postgres.HandleError(err)
		}

		locations[i] = location
	}

	return locations, nil
}

func (repo LocationRepository) GetByID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	return repo.getBy(ctx, "WHERE locations.id = $1", id)
}

func (repo LocationRepository) GetByUserID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	return repo.getBy(ctx, "WHERE user_id = $1", id)
}

func (repo LocationRepository) getBy(
	ctx context.Context,
	whereQuery string,
	value string,
) (*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		%s
	`

	query = fmt.Sprintf(query, whereQuery)

	location := models.Location{}

	err := repo.db.QueryRow(
		ctx,
		query,
		value).Scan(
		&location.ID,
		&location.Name,
		&location.Capacity,
		&location.TimeZone,
		&location.UserID,
	)
	if err != nil {
		return nil, postgres.HandleError(err)
	}

	err = repo.prepareLocation(ctx, &location)
	if err != nil {
		return nil, postgres.HandleError(err)
	}

	return &location, nil
}

func (repo LocationRepository) prepareLocation(
	ctx context.Context,
	location *models.Location,
) error {
	var checkInsToday []*models.CheckIn
	var checkInsYesterday []*models.CheckIn
	var err error

	loc, _ := time.LoadLocation(location.TimeZone)
	today := time.Now().In(loc)
	yesterday := today.AddDate(0, 0, -1)

	checkInsToday, err = repo.checkins.GetAllOfDay(ctx, location.ID, today)
	if err != nil {
		return err
	}

	checkInsYesterday, err = repo.checkins.GetAllOfDay(ctx, location.ID, yesterday)
	if err != nil {
		return err
	}

	location.SetCheckInRelatedFields(
		checkInsToday,
		checkInsYesterday,
	)

	err = location.NormalizeName()
	if err != nil {
		return err
	}

	return nil
}

func (repo LocationRepository) GetByName(
	ctx context.Context,
	name string,
) (*models.Location, error) {
	locations, err := repo.GetAll(ctx)
	if err != nil {
		return nil, postgres.HandleError(err)
	}

	for _, location := range locations {
		var output bool

		output, err = location.CompareNormalizedName(name)
		if err != nil {
			return nil, postgres.HandleError(err)
		}

		if output {
			return location, nil
		}
	}

	return nil, httptools.ErrRecordNotFound
}

func (repo LocationRepository) Create(
	ctx context.Context,
	name string,
	capacity int64,
	timeZone string,
	username string,
	password string,
) (*models.Location, error) {
	tx, err := repo.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return nil, postgres.HandleError(err)
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
		return nil, postgres.HandleError(err)
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
		return nil, postgres.HandleError(err)
	}

	err = location.NormalizeName()
	if err != nil {
		return nil, err
	}

	location.SetCheckInRelatedFields([]*models.CheckIn{}, []*models.CheckIn{})

	return &location, nil
}

func (repo LocationRepository) Update(
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

	tx, err := repo.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return postgres.HandleError(err)
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
		return postgres.HandleError(err)
	}

	rowsAffected := resultLocation.RowsAffected()
	if rowsAffected == 0 {
		return httptools.ErrRecordNotFound
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
		return postgres.HandleError(err)
	}

	rowsAffected := resultUser.RowsAffected()
	if rowsAffected == 0 {
		return httptools.ErrRecordNotFound
	}

	return nil
}

func (repo LocationRepository) Delete(
	ctx context.Context,
	location *models.Location,
) error {
	tx, err := repo.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return postgres.HandleError(err)
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
		return postgres.HandleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return httptools.ErrRecordNotFound
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
		return postgres.HandleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return httptools.ErrRecordNotFound
	}

	return nil
}
