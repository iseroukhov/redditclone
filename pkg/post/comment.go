package post

import (
	"redditclone/pkg/user"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Comment struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Body    string        `json:"body" bson:"body"`
	Author  *user.User    `json:"author" bson:"author"`
	Created time.Time     `json:"created" bson:"created"`
}

func NewComment(body string, author *user.User) *Comment {
	return &Comment{
		ID:      bson.NewObjectId(),
		Body:    body,
		Author:  author,
		Created: time.Now(),
	}
}
