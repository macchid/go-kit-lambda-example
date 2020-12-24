package user

type User struct {
	Email     string `json:"email", bson:"email"`
	FirstName string `json:"firstName", bson:"firstName"`
	LastName  string `json:"lastName", bson:"lastName"`
}
