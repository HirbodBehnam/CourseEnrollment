package course

// DepartmentID represents a department ID
type DepartmentID uint8

// Departments is a map of department id to department name
// We don't add departments while the app is running so no lock for you!
type Departments map[DepartmentID]string
