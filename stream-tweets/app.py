import json
import os

import boto3
from TwitterAPI import TwitterAPI

kinesis = boto3.client(
    'kinesis',
    aws_access_key_id=os.environ['AWS_KEY_SERVERLESS'],
    aws_secret_access_key=os.environ['AWS_SECRET_SERVERLESS'],
    region_name=os.environ['AWS_REGION']
)

twitter = TwitterAPI(
    os.environ['TWITTER_CONSUMER_KEY'],
    os.environ['TWITTER_CONSUMER_SECRET'],
    os.environ['TWITTER_ACCESS_TOKEN_KEY'],
    os.environ['TWITTER_ACCESS_TOKEN_SECRET']
)

params = {'track': 'bitcoin'}
stream = twitter.request('statuses/filter', params)

for tweet in stream:
    print(json.dumps(tweet))
