package types

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"FirstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}
