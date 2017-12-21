package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/orisano/gopl/ch04/ex11/github"
)

func main() {
	var user, repo string
	flag.StringVar(&user, "user", "orisano", "github user")
	flag.StringVar(&repo, "repo", "gopl", "github repository")
	flag.Parse()

	token := os.Getenv("GITHUB_TOKEN")
	client, err := github.NewClient(token)

	if err != nil {
		log.Fatal(err)
	}

	cmds := flag.Args()
	switch cmds[0] {
	case "get":
		if len(cmds) >= 2 {
			number, err := strconv.Atoi(cmds[1])
			if err != nil {
				log.Fatalln(err)
			}
			item, err := client.GetIssue(user, repo, number)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
		} else {
			issues, err := client.GetIssues(user, repo)
			if err != nil {
				log.Fatalln(err)
			}
			for _, item := range issues {
				fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
			}
		}
	case "create":
		if len(cmds) >= 3 {
			title := cmds[1]
			body := cmds[2]
			item, err := client.CreateIssue(user, repo, title, body)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
		}
	case "edit":
		if len(cmds) >= 4 {
			number, err := strconv.Atoi(cmds[1])
			if err != nil {
				log.Fatalln(err)
			}
			title := cmds[2]
			body := cmds[3]
			item, err := client.EditIssue(user, repo, number, title, body)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
		}
	case "close":
		if len(cmds) >= 2 {
			number, err := strconv.Atoi(cmds[1])
			if err != nil {
				log.Fatalln(err)
			}
			item, err := client.CloseIssue(user, repo, number)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
		}
	}

}
