package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lvoytek/discourse_client_go/pkg/discourse"
)

// Cache data used to avoid unnecessary Discourse API calls
var (
	// Cache of topics labelled by category slug and topic id
	topicsCache = map[string]map[int]*discourse.TopicData{}
)

var (
	timeLayout = time.RFC3339
)

func IntervalCollect(discourseClient *discourse.Client, categoryList []string, interval time.Duration) {
	for {
		collectData(discourseClient, categoryList)
		time.Sleep(interval)
	}
}

func collectData(discourseClient *discourse.Client, categoryList []string) {
	for _, categorySlug := range categoryList {
		go collectTopicsFromCategory(discourseClient, categorySlug)
	}

}

func collectTopicsFromCategory(discourseClient *discourse.Client, categorySlug string) {
	topics, ok := topicsCache[categorySlug]

	if !ok {
		topics = map[int]*discourse.TopicData{}
	}

	//TODO: Collect additional topics when there are more than 30
	categoryData, err := discourse.GetCategoryContentsBySlug(discourseClient, categorySlug)

	if err != nil {
		log.Println("Category data collection error for", categorySlug, "-", err)
		return
	}

	for _, topicOverview := range categoryData.TopicList.Topics {
		cachedTopic, topicExists := topics[topicOverview.ID]

		// If cached topic data exists, check if it actually needs to be updated
		if topicExists {
			cachedUpdateTime, err := time.Parse(timeLayout, cachedTopic.LastPostedAt)

			if err != nil {
				log.Println("Cached topic time parse error", err)
				continue
			}

			newTopicUpdateTime, err := time.Parse(timeLayout, topicOverview.LastPostedAt)

			if err != nil {
				log.Println("Downloaded topic time parse error", err)
				continue
			}

			fmt.Println(newTopicUpdateTime, cachedUpdateTime)

			if newTopicUpdateTime.Compare(cachedUpdateTime) > 0 {
				fmt.Println("Cached", cachedTopic.Title)
				continue
			}
		}

		// Get a new copy of the topic
		updatedTopic, err := discourse.GetTopicByID(discourseClient, topicOverview.ID)

		if err == nil {
			topics[topicOverview.ID] = updatedTopic
		} else {
			log.Println("Download topic error:", err)
		}

		fmt.Println(categorySlug, topicOverview.Title)
	}

	fmt.Println(topics)

	if len(topics) > 0 {
		topicsCache[categorySlug] = topics
	}
}
