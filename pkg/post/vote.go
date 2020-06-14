package post

type Vote struct {
	UserID string `json:"user" bson:"user"`
	Vote   int    `json:"vote" bson:"vote"`
}

func NewVote(userID string, value int) *Vote {
	return &Vote{
		UserID: userID,
		Vote:   value,
	}
}
