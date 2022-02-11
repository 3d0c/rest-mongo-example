package models

// ACLScheme meta model
type ACLScheme struct {
	Application ApplicationScheme `bson:"application,omitempty"`
	Permissions PermissionScheme  `bson:"permissions,omitempty"`
}
