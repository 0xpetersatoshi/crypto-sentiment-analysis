# Crypto Sentiment Analysis

This is a fun side project with the goal of pulling twitter posts related to crypto (mainly Bitcoin), running sentiment analysis, and then comparing that sentiment to crypto prices to see if there are any major correlations in twitter sentiment to price movements.

### Launching Python Twitter Stream from Docker

##### Build image

`docker build -t pbegle/twitter-streammer .`

##### Run Container

Need to pass in environment variables
```
docker run \
-e AWS_KEY_SERVERLESS=$AWS_KEY_SERVERLESS \
-e AWS_SECRET_SERVERLESS=$AWS_SECRET_SERVERLESS \
-e AWS_REGION=$AWS_REGION \
-e TWITTER_CONSUMER_KEY=$TWITTER_CONSUMER_KEY \
-e TWITTER_CONSUMER_SECRET=$TWITTER_CONSUMER_SECRET \
-e TWITTER_ACCESS_TOKEN_KEY=$TWITTER_ACCESS_TOKEN_KEY \
-e TWITTER_ACCESS_TOKEN_SECRET=$TWITTER_ACCESS_TOKEN_SECRET \
pbegle/twitter-streammer
```