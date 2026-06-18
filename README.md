#Gator
A simple Blog aggregator written in Go

## REQUIREMENTS
- Go 1.24.4 (or later)
- Postgres
- Goose
- sqlc

## Installation


First:
Setup Posgres (Follow instructions that you can find from Postgres to do so) and write the location (ip address or localhost) of the machine you are using to run it.

If you are not running this command on the same machine you will need to use the ip address and the port that you configured if changed from the default port of ":5432"

Grab the this repo using the golang install command shown below

```bash
go install github.com/chiprek/bootdev-blog-aggregator
```
Then:
now using psql on the machine you have the database installed on run the following command to create the database `gator`

If not on the same machine you will need to pass psql the same db_url you will put in the config to tell psql where to run the creation.

```bash
psql -c "CREATE DATABASE gator;"
```

next navigate to the `/sql/schema/` directory and run 
```bash 
goose postgres "postgres://<db_username>:<password>@<hostadd>:5432/gator" up
```

lastly in the root directory eg: `/bootdev-blog-aggregator/` of the project run sqlc to generate the required database related backend code. 


## Configuration

To use Gator, First make a `.gatorconfig.json` within your /home/ directory with the following congif

```json
{
  "db_url": "posgres://<user>:<password>@<location of machine running postgress eg localhost or ip address>:<port. default is 5432>/gator?sslmode=disable"
}
```
## Commands
To start register your first user using command below 
```bash
go run . register <username>
```
this will make a user and set the current logged into user to them.

if multiple users are present you must change to them with 
```bash
go run . login <username>
```
if hell froze over and you need to start anew run 
```bash
go run . reset
```
to list all users within the database run 
```bash
go run . users
```
to add a feed run 
```bash
#this also follows the feed for aggregation.
go run . addfeed <feed name> <feed URL>
```
to list available feeds in the database run
```bash
go run . feeds
```
to follow a feed run 
```bash
go run . follow <url of feed to follow>
```
to list the feeds followed run
```bash
go run . following
```
to unfollow a feed run
```bash
go run . unfollow <url to unfollow>
```
to pull the latest feed run
```bash
go run . agg <time between polling the feeds eg 5m for five minute intervals or 24h for 1 day intervals> 
```
> [!CAUTION]
Running too short of an interval will probably DDOS the machines that are holding your feed. Do not do that. It is considered not proper. Just note I a random person on the Internet will consider you a dingus if you do. dont be a dingus.

to browse your feeds run 
```bash
#if not specified this will only pull 2 feeds
go run . browse <NumOfFeeds>
```

Thank you for using my project. if you want to make things similar to this yet dont know the pratical knowlage to do so sign up at [boot.dev](https://www.boot.dev) it has finally gotten me this far.
