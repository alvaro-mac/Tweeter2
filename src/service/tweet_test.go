package service

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/Tweeter2/src/domain"
)

var user string
var tweetManager *TweetManager

func thereIsAUser(arg1 string) error {
	//Agregar "Alvaro" al slice de usuarios
	user = arg1
	return nil
}

func theUserPublishesATweet(arg1 string) error {
	//Publicar el Tweet
	fileTweetWriter := NewMemoryTweetWriter()
	tweetWriter := NewChannelTweetWriter(fileTweetWriter)
	tweetManager = NewTweetManager(tweetWriter)

	var tweet domain.Tweet

	text := arg1

	tweet = domain.NewTextTweet(user, text)

	quit := make(chan bool)

	// Operation
	tweetManager.PublishTweet(tweet, quit)

	<-quit

	return nil
}

func thereShouldBeTweet(arg1 int) error {
	//Controlar que haya 1 tweet
	num := tweetManager.CountTweetsByUser(user)

	if num == arg1 {
		return nil
	}else{
		return fmt.Errorf("there are %d tweets published", num)
	}
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^there is a user "([^"]*)"$`, thereIsAUser)
	s.Step(`^the user publishes a tweet "([^"]*)"$`, theUserPublishesATweet)
	s.Step(`^there should be (\d+) tweet$`, thereShouldBeTweet)
}
