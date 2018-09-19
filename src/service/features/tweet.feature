Feature: Tweet Manager
  In order to publish a tweet
  As a user
  I need to able to publish chat messages

  Rules:
  -Only text, image and quote tweets are allowed

  Scenario: Publishing a single text tweet
    Given there is a user "alvaro"
    When the user publishes a tweet "es un tweet"
    Then there should be 1 tweet