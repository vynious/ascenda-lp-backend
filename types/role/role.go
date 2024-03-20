package role

type Role struct {
	RoleId   string
	RoleName string
}

type RolePermission struct {
	RoleId    string
	CanCreate bool
	CanRead   bool
	CanUpdate bool
	CanDelete bool
	Resource  string
}
