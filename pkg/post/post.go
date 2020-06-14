package post

import (
	"errors"
	"redditclone/pkg/user"
	"time"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

type Post struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	Type             string        `json:"type" bson:"type"`
	Category         string        `json:"category" bson:"category"`
	Title            string        `json:"title" bson:"title"`
	Text             string        `json:"text,omitempty" bson:"text"`
	URL              string        `json:"url,omitempty" bson:"url"`
	Score            int           `json:"score" bson:"score"`
	Views            int           `json:"views" bson:"views"`
	UpvotePercentage int           `json:"upvotePercentage" bson:"upvote_percentage"`
	Author           *user.User    `json:"author" bson:"author"`
	Comments         []*Comment    `json:"comments" bson:"comments"`
	Votes            []*Vote       `json:"votes" bson:"votes"`
	Created          time.Time     `json:"created" bson:"created"`
}

func NewPost() *Post {
	return &Post{
		Comments: make([]*Comment, 0, 5),
		Votes:    make([]*Vote, 0, 5),
	}
}

func (p *Post) refresh() {
	var positive, score int
	for _, v := range p.Votes {
		score += v.Vote
		if v.Vote > 0 {
			positive++
		}
	}
	p.Score = score
	positive = positive * 100
	if positive > 0 {
		p.UpvotePercentage = positive / len(p.Votes)
	} else {
		p.UpvotePercentage = 0
	}
}

// --------------

type Repository struct {
	DB *mgo.Collection
}

func NewRepository(database *mgo.Collection) *Repository {
	return &Repository{
		DB: database,
	}
}

var (
	NotFound          = errors.New("post not found")
	InvalidID         = errors.New("invalid post id")
	InvalidCategoryID = errors.New("invalid category id")
	InvalidCommentID  = errors.New("invalid comment id")
)

func (r *Repository) IncViews(p *Post) error {
	count, err := r.DB.Find(bson.M{"_id": p.ID}).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return InvalidID
	}
	p.Views++
	err = r.DB.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAll() ([]*Post, error) {
	var posts []*Post
	err := r.DB.Find(bson.M{}).All(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *Repository) GetByID(id string) (*Post, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, InvalidID
	}
	_id := bson.ObjectIdHex(id)
	pst := &Post{}
	err := r.DB.Find(bson.M{"_id": _id}).One(&pst)
	if err != nil {
		return nil, err
	}
	return pst, nil
}

func (r *Repository) Create(p *Post) (*Post, error) {
	p.ID = bson.NewObjectId()
	p.Created = time.Now()
	p.refresh()
	err := r.DB.Insert(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) Update(p *Post) (*Post, error) {
	count, err := r.DB.Find(bson.M{"_id": p.ID}).Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, InvalidID
	}
	p.refresh()
	err = r.DB.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) Delete(id bson.ObjectId) error {
	err := r.DB.Remove(bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByAuthor(username string) ([]*Post, error) {
	posts := make([]*Post, 0, 5)
	err := r.DB.Find(bson.M{
		"author.username": username,
	}).All(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *Repository) GetByCategory(slug string) ([]*Post, error) {
	posts := make([]*Post, 0, 5)
	err := r.DB.Find(bson.M{
		"category": slug,
	}).All(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
