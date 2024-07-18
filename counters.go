package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	postEdits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discourse_post_edits_total",
			Help: "The total number of edits on a specific post.",
		},
		[]string{
			"post_id",
		},
	)

	topicComments = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "discourse_topic_comments_total",
			Help: "The total number of comments on a specific topic.",
		},
		[]string{
			"topic_id",
		},
	)
)
