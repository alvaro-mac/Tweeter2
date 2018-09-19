package service

import (
	"errors"
	"github.com/Tweeter2/src/domain"
)

type TweetManager struct {
	writer *ChannelTweetWriter
	idCounter int
}

func (tm *TweetManager) PublishTweet(s domain.Tweet, quit chan bool) (int, error) {
	if s.GetUser() == "" {
		return 0, errors.New("user is required");
	}

	if s.GetText() == "" {
		return 0, errors.New("text is required");
	}

	if len(s.GetText()) > 140 {
		return 0, errors.New("text exceeds 140 characters");
	}

	tweetsToWrite := make(chan domain.Tweet)

	go tm.writer.WriteTweet(tweetsToWrite, quit)
	tm.idCounter++
	s.SetId(tm.idCounter)
	tweetsToWrite <- s

	close(tweetsToWrite)

	return s.GetId(), nil
}

func (tm *TweetManager) GetTweet() domain.Tweet{
	return tm.writer.writer.GetTweets()[0]
}

func NewTweetManager(wrt *ChannelTweetWriter) *TweetManager {
	tweetManager := TweetManager{writer:wrt,idCounter:0}
	return &tweetManager
}

func (tm *TweetManager) GetTweets() []domain.Tweet {
	return tm.writer.writer.GetTweets()
}

func (tm *TweetManager) GetTweetById(id int) domain.Tweet {
	var result domain.Tweet
	for _, value := range tm.writer.writer.GetTweets() {
		if value.GetId() == id {
			result = value
			break
		}
	}

	return result
}

func (tm *TweetManager) CountTweetsByUser(user string) int {
	var result int = 0
	for _, value := range tm.writer.writer.GetTweets() {
		if value.GetUser() == user {
			result++
		}
	}

	return result
}

func (tm *TweetManager) GetTweetsByUser(user string) []domain.Tweet {
	var result map[string][]domain.Tweet
	result = make(map[string][]domain.Tweet)
	for _, value := range tm.writer.writer.GetTweets() {
		result[value.GetUser()] = append(result[value.GetUser()], value)
	}

	return result[user]
}