//nolint:gochecknoglobals //global var for access
package main

import "check-in/api/internal/models"

var allRoles = []models.Role{
	models.AdminRole,
	models.ManagerRole,
	models.DefaultRole,
}

var adminRole = []models.Role{
	models.AdminRole,
}

var defaultRole = []models.Role{
	models.DefaultRole,
}

var managerAndAdminRole = []models.Role{
	models.ManagerRole,
	models.AdminRole,
}
