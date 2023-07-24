//nolint:gochecknoglobals //global var for access
package main

import "check-in/api/internal/models"

var allRoles = []models.Roles{
	models.AdminRole,
	models.ManagerRole,
	models.DefaultRole,
}

var adminRole = []models.Roles{
	models.AdminRole,
}

var defaultRole = []models.Roles{
	models.DefaultRole,
}

var managerAndAdminRole = []models.Roles{
	models.ManagerRole,
	models.AdminRole,
}
