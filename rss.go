package main

import _ "github.com/lib/pq"
import (
	"context"
	"database/sql"
	"encoding/xml"
	"net/http"
	"io"
	"log"
	"time"
	"sync"

	"github.com/lcphutchinson/database"
)

type RSSItem struct {
	Title		string		`xml:"title"`
	Link		string		`xml:"link"`
	Description	string		`xml:"description"`
}

type RSSFeed struct {
	Title		string		`xml:"title"`
	Description	string		`xml:"description"`
	Items		[]RSSItem	`xml:"item"`
}

type RSSBody struct {
	XMLName		xml.Name	`xml:"rss"`
	Channel		RSSFeed 	`xml:"channel"`
}

func fetchPosts(url string) (RSSBody, error){
	rss := RSSBody{}
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error in http.Get:")
		return rss, err
	}
	
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error in io.Read:")
		return rss, err
	}

	err = xml.Unmarshal(body, &rss)
	if err != nil {
		log.Println("Error in Unmarshal:")
		return rss, err
	}
	return rss, nil
}

type worker struct {
	DB		*database.Queries
	batchSize	int32
	loopInterval	time.Duration
}

func (w *worker) Work(){
	var batch sync.WaitGroup
	for true{
		feeds, err := w.DB.GetNextNFeeds(context.TODO(), w.batchSize)
		if err != nil {
			log.Printf("DB fetch failed with error { %v \n }\nContinuing...\n", err)
			continue
		}
		for _, feed := range feeds{
			batch.Add(1)
			go func(){
				defer batch.Done()
				data, err := fetchPosts(feed.Url)
				if err != nil {
					log.Println(err)
					return
				}
				for _, post := range data.Channel.Items {
					var nullDesc sql.NullString
					nullDesc.Scan(post.Description)
					params := database.CreatePostParams{
						Title:		post.Title,
						Url:		post.Link,
						Description:	nullDesc,
						FeedID:		feed.ID,
					}
					w.DB.CreatePost(context.TODO(), params)
				}
			}()
		}
		time.Sleep(w.loopInterval)
	}
}

