package bugle

type user struct {
	Name  string
	Email string
}

func User(email, token string) user {
	// TODO: derive name and validate token
	return user{"Connor", email}
}
