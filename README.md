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
guide. Once Postgres is installed pickup from step 5 of the Boot.dev instructions.

## Setting up your config file
In your home directory create a file with the name ".gatorconfig.json".
It needs to look like the following
`
    {
        "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
        "current_user_name": "<your_name_here>"
    }
`
The general form of the string after "db_ur": is:
` "postgres://<postgres_username>:<postgres_password>@localhost:<port_number>/gator?sslmode=disable" `
So if you chose a different username, password, or port for you Postgres installation change them 
accordingly.


## Running Gator and a fixing a problem that *may* occur
So you've installed Go, and Postgres, and then you've created your config file and database.
Now all there is left to do is to `go build` or `go install` and the use `./gator` for the
program to start putting users, RSS feeds and posts into the database. 




