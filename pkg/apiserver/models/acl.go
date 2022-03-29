package models

// ACLScheme meta model
type ACLScheme struct {
	Application *ApplicationScheme `bson:"application,omitempty" json:"application"`
	Permissions *PermissionScheme  `bson:"permissions,omitempty" json:"permissions"`
}
