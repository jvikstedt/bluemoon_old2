package room

type User struct {
	id int
}

func NewUser(id int) *User {
	return &User{
		id: id,
	}
}

func (p *User) ID() int {
	return p.id
}
