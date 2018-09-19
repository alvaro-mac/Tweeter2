package service

import (
	"github.com/Tweeter2/src/domain"
	"os"
)

type MemoryTweetWriter struct {
	Tweets []domain.Tweet
}

type FileMemoryTweetWriter interface {
	GetTweets() []domain.Tweet
	AddTweet(tweet domain.Tweet)
}

type ChannelTweetWriter struct {
	writer FileMemoryTweetWriter
}

type FileTweetWriter struct {
	file *os.File
}

func NewFileTweetWriter() *FileTweetWriter {

	file, _ := os.OpenFile(
		"tweets.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)

	writer := new(FileTweetWriter)
	writer.file = file

	return writer
}

func (*FileTweetWriter) GetTweets() []domain.Tweet {
	return nil
}

func (mw *MemoryTweetWriter) GetTweets() []domain.Tweet {
	return mw.Tweets
}


func NewMemoryTweetWriter() *MemoryTweetWriter {
	writer := MemoryTweetWriter{}
	return &writer
}

func NewChannelTweetWriter(memwriter FileMemoryTweetWriter) *ChannelTweetWriter {
	writer := ChannelTweetWriter{writer:memwriter}
	return &writer
}

func (mw *MemoryTweetWriter) AddTweet (t domain.Tweet){
	mw.Tweets = append(mw.Tweets, t)
}

func (writer *FileTweetWriter) AddTweet(tweet domain.Tweet) {

	if writer.file != nil {
		byteSlice := []byte(tweet.PrintableTweet() + "\n")
		writer.file.Write(byteSlice)
	}
}

func (cw *ChannelTweetWriter) WriteTweet(tChan chan domain.Tweet, quit chan bool) {
	for t := range tChan {

		cw.writer.AddTweet(t)
	}
	quit <- true
	return
}