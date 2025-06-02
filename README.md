# gator cli tool for rss feed
This tool is a part of the boot.dev curriculum. The program creates a postgres
database storing user, rss feed, and rss post information. It also acts as a
command line interface to view posts.

## Installation Guide
You will need to install the following:
1. Go programming language
2. Postgres
You will also need to create a config file and do some initial set up on your
local postgres installation.

## Installing Go
To install Go follow the instructions at the Go website: [Go.dev](https://go.dev/doc/install).
I would also recommend looking up operating system specific instructions on Google.
Once you have installed go, verify it is working by running the command:
`go version`
The code was written with go 1.23.6 so if you are using an older version and it doesn't work,
try updating your go version.

## Installing Postgres and setting up the database
Follow the Postgres installation and database setup guide provided by Boot.dev for this project at:
[Boot.dev](https://www.boot.dev/lessons/74bea1f2-19cd-4ea9-966e-e2ca9dd1dfa9).
If you don't have a Boot dev account just close the promp that you are in guest mode.

If you are using a *fancy* distro of Linux, install Postgres following your preferred distro's
guide. Once Postgres is installed, follow the Boot.dev instructions from step 5 onwards.

## Setting up your config file
In your home directory create a file with the name ".gatorconfig.json".
It needs to look like the following
`
    {
        "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
        "current_user_name": "<your_name_here>"
    }
`
The general form of the string after "db_url": is:
` "postgres://<postgres_username>:<postgres_password>@localhost:<port_number>/gator?sslmode=disable" `
So if you chose a different username, password, or port for you Postgres installation change them 
accordingly.

## Getting databases up and running
Once you have created a database and have a config file it's to create the 
tables you will need. Open a terminal use command:
`psql "postgres://postgres:@localhost:5432/gator"`
to access the gator database using psql. Then copy and paste the create table
statements found in sql/schema portion of the project. Do them in numerical
order. So `CREATE TABLE users ...` first, then feeds, feed_follows, and 
finally posts.

Once you've done that in psql run command `\dt` and you should see the tables
you just created. You can then enter `\q` to exit out of psql.


## Running gator and a fixing a problem that *may* occur
So you've installed Go, and Postgres, and then you've created your config file and database.
Now all there is left to do is to `go build` or `go install` and the use `./gator` for the
program to start putting users, RSS feeds and posts into the database. 

## gator commands
1. `register <username>` Add a user to the database.
2. `login <username>` Set user to current_user_name.
3. `reset` Wipes database.
4. `users` List all users.
5. `addfeed <feed_name> <feed_url>` Add an RSS feed to database.
6. `feeds` Show all feeds in database.
7. `follow <feed_url>` Add feed to current users followed feeds.
8. `following` Show all feeds current user is following.
9. `unfollow <feed_url>` Removes feed from user's followed feeds
10. `browse <int_count optional>` Browse latest updated feeds.
11. No command given will run gator aggregate for the current user. It will 
    update feeds based on which feeds are the oldest/haven't been updated at 
    all.

## Example usage
`./gator register steve`
`./gator login steve`
`./gator addfeed "HackerNews" "https://hnrss.org/newest"`
`./gator follow "https://hnrss.org/newest"`
`./gator` This will go get posts for user steve from Hackernews every 15s, use CTRL+C to stop.
`./gator browse` Will return two posts from Hackernews.


