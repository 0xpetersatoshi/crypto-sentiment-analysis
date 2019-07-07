import json
import os
import logging

import boto3
from TwitterAPI import TwitterAPI

logging.basicConfig(
    format='[%(asctime)s - %(name)s - %(levelname)s] - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

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

tracks = [
    'bitcoin',
    'etherium',
    'litecoin',
    'ripple',
    'btc',
    'eth',
    'ltc',
    'xrp',
    'crypto',
    'cryptocurrency',
    'cryptocurrencies'
    ]
params = {'track': tracks}
stream = twitter.request('statuses/filter', params)

tweets_processed = 0
for tweet in stream:
    try:
        if tweet['retweeted'] or tweet['text'].startswith('RT '):
            continue
    except KeyError as e:
        logger.error('Exception ocurred', exc_info=True)
    
    tweets_processed += 1
    if tweets_processed % 20 == 0:
        logger.info(f'{tweets_processed} tweets proccessed')

    kinesis.put_record(
        StreamName='twitter',
        Data=json.dumps(tweet),
        PartitionKey=tweet['lang']
        )
