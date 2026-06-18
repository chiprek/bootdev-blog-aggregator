package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/chiprek/bootdev-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var feed RSSFeed
	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return &feed, err
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("%s needs 1 argument, the time between requests", cmd.name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeed(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("%s needs 2 arguments, name and url, only got %d", cmd.name, len(cmd.args))
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{ID: uuid.New(), Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})

	return err
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("* %s (%s) by %s\n", feed.Name, feed.Url, feed.UserName.String)
	}
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Followed by %s\n", user.Name)

	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}

func handlerFeedFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("%s needs 1 argument, the url to follow for the current user", cmd.name)
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}

	fmt.Printf("%s follows %s\n", follow.UserName, follow.FeedName)

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("%s needs 1 argument, the url to unfollow", cmd.name)
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	return s.db.RemoveFollow(context.Background(), database.RemoveFollowParams{UserID: user.ID, FeedID: feed.ID})
}

func scrapeFeed(s *state) error {
	next, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), next.ID)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), next.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Fetched feed %s\n", feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		params := database.CreatePostParams{
			ID:     uuid.New(),
			Title:  item.Title,
			Url:    item.Link,
			FeedID: next.ID,
		}

		if item.Description != "" {
			params.Description = sql.NullString{Valid: true, String: item.Description}
		}

		pub, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			params.PublishedAt = sql.NullTime{Valid: true, Time: pub}
		}
		s.db.CreatePost(context.Background(), params)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var err error
	limit := 2
	if len(cmd.args) > 0 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}

	posts, err := s.db.GetPostsByUser(context.Background(), database.GetPostsByUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("%+v\n", post)
	}

	return nil
}
