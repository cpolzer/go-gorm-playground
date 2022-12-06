package user

type UserDto struct {
	ID        uint
	FirstName string
	LastName  string
	Address   AddressDto
}

type AddressDto struct {
	ID         uint
	Street     string
	PostalCode string
}
