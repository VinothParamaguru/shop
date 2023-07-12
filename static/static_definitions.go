package static

// This is bad programming, needs change

type TokenParams struct {
	id       int
	validity int
}

func GetSessionTokenType() int {
	return 1
}

func GetSessionTokenParams() TokenParams {
	return TokenParams{id: 1, validity: 7200}
}
