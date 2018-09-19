package main

import (
	"github.com/Tweeter2/src/domain"
	"github.com/Tweeter2/src/service"
	"github.com/abiosoft/ishell"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var tweetManager *service.TweetManager

func main() {

	quit := make(chan bool)

	fileTweetWriter := service.NewMemoryTweetWriter()
	tweetWriter := service.NewChannelTweetWriter(fileTweetWriter)

	tweetManager = service.NewTweetManager(tweetWriter)

	shell := ishell.New()
	shell.SetPrompt("Tweeter >> ")
	shell.Print("Type 'help' to know commands\n")

	shell.AddCmd(&ishell.Cmd{
		Name: "publishTweet",
		Help: "Publishes a tweet",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			c.Print("Type your username: ")

			user := c.ReadLine()

			c.Print("Type your tweet: ")

			text := c.ReadLine()

			tweet := domain.NewTextTweet(user, text)

			id, err := tweetManager.PublishTweet(tweet, quit)

			if err == nil {
				c.Printf("Tweet sent with id: %v\n", id)
			} else {
				c.Print("Error publishing tweet:", err)
			}

			return
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "publishImageTweet",
		Help: "Publishes a tweet with an image",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			c.Print("Type your username: ")

			user := c.ReadLine()

			c.Print("Type your tweet: ")

			text := c.ReadLine()

			c.Print("Type the url of your image: ")

			url := c.ReadLine()

			tweet := domain.NewImageTweet(user, text, url)

			id, err := tweetManager.PublishTweet(tweet, quit)

			if err == nil {
				c.Printf("Tweet sent with id: %v\n", id)
			} else {
				c.Print("Error publishing tweet:", err)
			}

			return
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "publishQuoteTweet",
		Help: "Publishes a tweet with a quote",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			c.Print("Type your username: ")

			user := c.ReadLine()

			c.Print("Type your tweet: ")

			text := c.ReadLine()

			c.Print("Type the id of the tweet you want to quote: ")

			id, _ := strconv.Atoi(c.ReadLine())

			quoteTweet := tweetManager.GetTweetById(id)

			tweet := domain.NewQuoteTweet(user, text, quoteTweet)

			id, err := tweetManager.PublishTweet(tweet, quit)

			if err == nil {
				c.Printf("Tweet sent with id: %v\n", id)
			} else {
				c.Print("Error publishing tweet:", err)
			}

			return
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "showTweet",
		Help: "Shows the last tweet",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			tweet := tweetManager.GetTweet()

			c.Println(tweet)

			return
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "showTweets",
		Help: "Shows all the tweets",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			tweets := tweetManager.GetTweets()

			c.Println(tweets)

			return
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "showTweetById",
		Help: "Shows the tweet with the provided id",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			c.Print("Type the id: ")

			id, _ := strconv.Atoi(c.ReadLine())

			tweet := tweetManager.GetTweetById(id)

			c.Println(tweet)

			return
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "countTweetsByUser",
		Help: "Counts the tweets published by the user",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			c.Print("Type the user: ")

			user := c.ReadLine()

			count := tweetManager.CountTweetsByUser(user)

			c.Println(count)

			return
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "showTweetsByUser",
		Help: "Shows the tweets published by the user",
		Func: func(c *ishell.Context) {

			defer c.ShowPrompt(true)

			c.Print("Type the user: ")

			user := c.ReadLine()

			tweets := tweetManager.GetTweetsByUser(user)

			c.Println(tweets)

			return
		},
	})

	go routeInit()

	shell.Run()

}

func routeInit() {
	router := gin.Default()

	router.GET("/showTweet", showTweet)
	router.GET("/showTweets", showTweets)
	router.GET("/showTweetById/:id", showTweetById)
	router.POST("/publishTweet", publishTweet)
	router.POST("/publishImageTweet", publishImageTweet)
	router.POST("/publishQuoteTweet", publishQuoteTweet)
	router.GET("/countTweetsByUser/:user", countTweetsByUser)
	router.GET("/showTweetsByUser/:user", showTweetsByUser)

	router.Run(":8090")
}

func showTweet(c *gin.Context) {
	c.String(http.StatusOK, tweetManager.GetTweet().PrintableTweet())
}

func publishTweet(c *gin.Context) {
	var t *domain.TextTweet
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quit := make(chan bool)
	t = domain.NewTextTweet(t.GetUser(),t.GetText())
	tweetManager.PublishTweet(t, quit)
	<-quit
	c.String(http.StatusOK, "Tweet submited")
}

func showTweets(c *gin.Context) {
	c.JSON(http.StatusOK, tweetManager.GetTweets())
}

func showTweetById(c *gin.Context){
	id, _ := strconv.Atoi(c.Param("id"))
	c.JSON(http.StatusOK, tweetManager.GetTweetById(id))
}

func publishImageTweet(c *gin.Context) {
	var t *domain.ImageTweet
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tweet := domain.NewImageTweet(t.User, t.Text, t.URL)
	quit := make(chan bool)
	tweetManager.PublishTweet(tweet, quit)
	<-quit
	c.String(http.StatusOK, "Tweet submited")
}

func publishQuoteTweet(c *gin.Context) {
	var t *domain.QuoteTweet
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quoteTweet := tweetManager.GetTweetById(t.IdQuotedTweet)

	tweet := domain.NewQuoteTweet(t.User, t.Text, quoteTweet)
	quit := make(chan bool)
	tweetManager.PublishTweet(tweet, quit)
	<-quit
	c.String(http.StatusOK, "Tweet submited")
}

func countTweetsByUser(c *gin.Context) {
	user := c.Param("user")
	count := tweetManager.CountTweetsByUser(user)
	c.String(http.StatusOK, "Tweets count: %d", count)
}

func showTweetsByUser(c *gin.Context) {
	user := c.Param("user")
	tweets := tweetManager.GetTweetsByUser(user)
	c.JSON(http.StatusOK,tweets)
}