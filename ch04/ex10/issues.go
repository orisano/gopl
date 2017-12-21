package main

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/orisano/gopl/ch04/ex10/github"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	var month, year, after []*github.Issue
	now := time.Now()
	for _, item := range result.Items {
		hours := now.Sub(item.CreatedAt).Hours()
		if hours < 24*31 {
			month = append(month, item)
			continue
		}
		if hours < 24*365 {
			year = append(year, item)
			continue
		}
		after = append(after, item)
	}

	fmt.Printf("%d issues:\n", result.TotalCount)
	fmt.Println("======= 一ヶ月以内 =========")
	for _, item := range month {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
	fmt.Println("======= 一年以内 =========")
	for _, item := range year {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
	fmt.Println("======= 一年以上 =========")
	for _, item := range after {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}
