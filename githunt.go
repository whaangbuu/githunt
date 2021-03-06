package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	flag "github.com/ogier/pflag"
	"github.com/olekukonko/tablewriter"
	"github.com/whaangbuu/go-loading/loading"
)

// flags ...
var (
	user string
)

const (
	apiURL       = "https://api.github.com"
	userEndpoint = "/users/"
)

// User is our simple blue-print.
type User struct {
	Login             string      `json:"login"`
	ID                int         `json:"id"`
	AvatarURL         string      `json:"avatar_url"`
	GravatarID        string      `json:"gravatar_id"`
	URL               string      `json:"url"`
	HTMLURL           string      `json:"html_url"`
	FollowersURL      string      `json:"followers_url"`
	FollowingURL      string      `json:"following_url"`
	GistsURL          string      `json:"gists_url"`
	StarredURL        string      `json:"starred_url"`
	SubscriptionsURL  string      `json:"subscriptions_url"`
	OrganizationsURL  string      `json:"organizations_url"`
	ReposURL          string      `json:"repos_url"`
	EventsURL         string      `json:"events_url"`
	ReceivedEventsURL string      `json:"received_events_url"`
	Type              string      `json:"type"`
	SiteAdmin         bool        `json:"site_admin"`
	Name              string      `json:"name"`
	Company           string      `json:"company"`
	Blog              string      `json:"blog"`
	Location          string      `json:"location"`
	Email             string      `json:"email"`
	Hireable          interface{} `json:"hireable"`
	Bio               string      `json:"bio"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

func main() {

	flag.Parse()

	if flag.NFlag() == 0 {
		printUsage()
	}

	users := strings.Split(user, ",")
	fmt.Printf("Searching user(s): %s\n", users)
	loading.StartNew("Fetching")
	printTabularData(users)
}

func init() {
	flag.StringVarP(&user, "user", "u", "", "Search Users")
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

func printTabularData(users []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"User ID", "Username", "Name", "Bio", "Followers", "Following"})

	for _, u := range users {
		result, err := GetUserByUsername(u)
		// Convert string to int.
		IDstr := strconv.Itoa(result.ID)
		followingStr := strconv.Itoa(result.Following)
		followerStr := strconv.Itoa(result.Followers)

		if err != nil {
			log.Printf("ERROR DUE TO: %s", err.Error())
			return
		}

		data := [][]string{
			[]string{IDstr, result.Login, result.Name, result.Bio, followerStr, followingStr},
		}
		for _, v := range data {
			table.Append(v)
		}
	}
	fmt.Println()
	table.Render()
}

func GetUserByUsername(username string) (*User, error) {
	res, err := http.Get(apiURL + userEndpoint + username)
	var user User

	if err != nil {
		return nil, fmt.Errorf("ERROR DUE TO: %s", err.Error())
	}
	// Dont forget to close
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ERROR DUE TO: %s", err.Error())
	}
	json.Unmarshal(body, &user)
	return &user, nil
}
