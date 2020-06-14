package main

import (
	"database/sql"
	"fmt"
	"log"
	"redditclone/pkg/post"
	"redditclone/pkg/server"
	"redditclone/pkg/user"
	"time"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysql, err := createMySQLConnect()
	if err != nil {
		log.Fatal(err)
	}
	mongoDB, err := createMongoDBDatabase()
	if err != nil {
		log.Fatal(err)
	}
	// createPosts(mongoDB)
	s := server.New(mysql, mongoDB)
	if err := s.Start(":9000"); err != nil {
		log.Fatal(err)
	}
}

func createMySQLConnect() (*sql.DB, error) {
	dsn := fmt.Sprintf("root:password@tcp(localhost:3306)/golang?%s&%s", "charset=utf8", "interpolateParams=true")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func createMongoDBDatabase() (*mgo.Database, error) {
	sess, err := mgo.Dial("mongodb://root:password@localhost")
	if err != nil {
		return nil, err
	}
	return sess.DB("golang"), nil
}

func createPosts(db *mgo.Database) {
	collection := db.C("posts")
	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&post.Post{
			ID:       bson.NewObjectId(),
			Type:     "text",
			Category: "music",
			Title:    "Test text Post",
			Text:     "hello, world!",
			Author: &user.User{
				ID:       "4a8f0b7cc134d99602b3c9af",
				Username: "user",
			},
			Comments: []*post.Comment{
				{
					ID:   bson.NewObjectId(),
					Body: "this is comment",
					Author: &user.User{
						ID:       "188d2fe9ad7154960f279d8d",
						Username: "goodnick",
					},
					Created: time.Now(),
				},
			},
			Votes:   nil,
			Created: time.Now(),
		})
		collection.Insert(&post.Post{
			ID:       bson.NewObjectId(),
			Type:     "text",
			Category: "music",
			Title:    "Test text Post 2",
			Text:     "examplex text",
			Author: &user.User{
				ID:       "4a8f0b7cc134d99602b3c9af",
				Username: "user",
			},
			Comments: []*post.Comment{
				{
					ID:   bson.NewObjectId(),
					Body: "this is comment",
					Author: &user.User{
						ID:       "188d2fe9ad7154960f279d8d",
						Username: "goodnick",
					},
					Created: time.Now(),
				},
				{
					ID:   bson.NewObjectId(),
					Body: "hello!",
					Author: &user.User{
						ID:       "4a8f0b7cc134d99602b3c9af",
						Username: "user",
					},
					Created: time.Now(),
				},
			},
			Votes:   nil,
			Created: time.Now(),
		})
		collection.Insert(&post.Post{
			ID:       bson.NewObjectId(),
			Type:     "link",
			Category: "programming",
			Title:    "article about programming",
			URL:      "https://google.com",
			Author: &user.User{
				ID:       "188d2fe9ad7154960f279d8d",
				Username: "goodnick",
			},
			Comments: []*post.Comment{
				{
					ID:   bson.NewObjectId(),
					Body: "great!",
					Author: &user.User{
						ID:       "4a8f0b7cc134d99602b3c9af",
						Username: "user",
					},
					Created: time.Now(),
				},
				{
					ID:   bson.NewObjectId(),
					Body: "thx!",
					Author: &user.User{
						ID:       "188d2fe9ad7154960f279d8d",
						Username: "goodnick",
					},
					Created: time.Now(),
				},
			},
			Votes:   nil,
			Created: time.Now(),
		})
	}
}
