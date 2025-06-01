package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/itlightning/dateparse"
	"github.com/janiv/gator/internal/config"
	"github.com/janiv/gator/internal/database"
)

type State struct {
	db  *database.Queries
	cfg *config.Config
}

func NewState() *State {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &State{
		cfg: &c,
	}
}

type Command struct {
	name string
	args []string
}

type Commands struct {
	commands map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error {
	val, exists := c.commands[cmd.name]
	if !exists {
		return errors.New("yo you missing a command")
	}
	return val(s, cmd)
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commands[name] = f
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing args")
	}
	name := cmd.args[0]
	usr_check, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("%s is not in database\n", name)
		return err
	}
	set_err := s.cfg.SetUser(usr_check.Name)
	if set_err != nil {
		return set_err
	}
	fmt.Println("user set")
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing args")
	}
	fmt.Printf("args= %s\n", cmd.args[0])
	name := cmd.args[0]
	usr_check, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("%s not in database, attempting to create\n", name)
	} else {
		fmt.Printf("%s already exists\n", usr_check.Name)
		return errors.New("user already exists")
	}
	curr_time := time.Now()
	db_params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: curr_time,
		UpdatedAt: curr_time,
		Name:      name,
	}
	fmt.Println(db_params)
	usr, err := s.db.CreateUser(context.Background(), db_params)
	if err != nil {
		return err
	}
	handlerLogin(s, cmd)
	fmt.Printf("User %s was created\n", usr.Name)
	return nil
}

func handlerReset(s *State, cmd Command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("database reset")
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	// ok this is a bit stupid but we doing it anyway
	curr := users[0]

	for _, usr := range users {
		if usr.UpdatedAt.After(curr.UpdatedAt) {
			curr = usr
		}
	}
	for _, usr := range users {
		fmt.Printf("* %s ", usr.Name)
		if usr == curr {
			fmt.Print("(current)")
		}
		fmt.Print("\n")
	}
	return nil
}

func handlerAgg(time_between_reqs string) error {
	timeBetweenRequests, time_err := time.ParseDuration(time_between_reqs)
	if time_err != nil {
		return time_err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		fmt.Println("Making a request!")
		scrapeFeeds()
		fmt.Println("WAITING TO MAKE NEXT REQUEST")
	}
}

func scrapeFeeds() error {
	s := NewState()
	db, err := sql.Open("postgres", s.cfg.DbURL)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.db = dbQueries
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	curr_time := time.Now()
	db_params := database.MarkFeedFetchedParams{
		LastFetched: sql.NullTime{
			Time:  curr_time,
			Valid: true,
		},
		ID: next_feed.ID,
	}
	mar_err := s.db.MarkFeedFetched(context.Background(), db_params)
	if mar_err != nil {
		return mar_err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	feed, err := FetchFeed(ctx, next_feed.Url)
	if err != nil {
		return err
	}
	ch := &feed.Channel
	items := ch.Item
	post_time := time.Now()
	for _, i := range items {
		pretty_title := html.UnescapeString(i.Title)
		desc := html.UnescapeString(i.Description)
		pub_date := i.PubDate
		pub_date_parsed, err := dateparse.ParseLocal(pub_date)
		if err != nil {
			return err
		}
		db_params_posts := database.CreatePostParams{
			CreatedAt: post_time,
			UpdatedAt: post_time,
			Title: sql.NullString{
				String: pretty_title,
				Valid:  true,
			},
			PostUrl: i.Link,
			PostDescription: sql.NullString{
				String: desc,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  pub_date_parsed,
				Valid: true,
			},
			FeedID: next_feed.ID,
		}
		post, err := s.db.CreatePost(context.Background(), db_params_posts)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Post: %s was created\n", post.PostUrl)
	}

	return nil
}

func handlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("missing args")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	curr_time := time.Now()
	db_params := database.CreateFeedParams{
		CreatedAt: curr_time,
		UpdatedAt: curr_time,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), db_params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed %s was created\n", feed.Name)
	curr_time = time.Now()
	feed_follow_params := database.CreateFeedFollowParams{
		CreatedAt: curr_time,
		UpdatedAt: curr_time,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	res, err := s.db.CreateFeedFollow(context.Background(), feed_follow_params)
	if err != nil {
		return err
	}
	fmt.Printf("Updated db to show %s following %s", res.UserID, res.FeedName)
	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, f := range feeds {
		usr_name, usr_err := s.db.GetUserByID(context.Background(), f.UserID)
		if usr_err != nil {
			return err
		}
		fmt.Printf("%s %s %s", f.Name, f.Url, usr_name)
	}
	return nil
}

func handlerFollow(s *State, cmd Command, user database.User) error {
	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	curr_time := time.Now()
	db_params := database.CreateFeedFollowParams{
		CreatedAt: curr_time,
		UpdatedAt: curr_time,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	res, err := s.db.CreateFeedFollow(context.Background(), db_params)
	if err != nil {
		return err
	}
	fmt.Printf("%s sucessfully followed a feed\n", res.UserName)
	fmt.Printf("User: %s is following %s\n", user.Name, feed.Name)
	return nil

}

func handlerFollowing(s *State, cmd Command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s is following:\n", user.Name)
	for _, f := range feeds {
		fmt.Printf("%s\n", f.FeedName)
	}
	return nil
}

func handlerUnfollow(s *State, cmd Command, user database.User) error {
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	db_params := database.DeleteByIDsParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}
	delete_err := s.db.DeleteByIDs(context.Background(), db_params)
	if delete_err != nil {
		return delete_err
	}
	return nil
}

func handlerBrowse(s *State, cmd Command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) != 0 {
		temp_limit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = int32(temp_limit)
	}

	db_params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), db_params)
	if err != nil {
		return err
	}
	for _, p := range posts {
		var title string
		var desc string
		if p.Title.Valid {
			title = p.Title.String
		}
		if p.PostDescription.Valid {
			desc = p.PostDescription.String
		}
		fmt.Printf("%s\n%s\n", title, desc)

	}
	return nil
}
