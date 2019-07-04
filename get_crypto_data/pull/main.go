package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type cryptoData struct {
	Response   string `json:"Response"`
	Type       int    `json:"Type"`
	Aggregated bool   `json:"Aggregated"`
	FromSymbol string `json:"FromSymbol"`
	ToSymbol   string `json:"ToSymbol"`
	Data       []struct {
		Time       int     `json:"time"`
		Close      float64 `json:"close"`
		High       float64 `json:"high"`
		Low        float64 `json:"low"`
		Open       float64 `json:"open"`
		Volumefrom float64 `json:"volumefrom"`
		Volumeto   float64 `json:"volumeto"`
	} `json:"Data"`
	TimeTo            int  `json:"TimeTo"`
	TimeFrom          int  `json:"TimeFrom"`
	FirstValueInArray bool `json:"FirstValueInArray"`
	ConversionType    struct {
		Type             string `json:"type"`
		ConversionSymbol string `json:"conversionSymbol"`
	} `json:"ConversionType"`
	RateLimit struct {
	} `json:"RateLimit"`
	HasWarning bool `json:"HasWarning"`
}

var (
	region     = os.Getenv("AWS_REGION_PB")
	bucket     = os.Getenv("S3_BUCKET_DATA_LAKE")
	apiKey     = os.Getenv("CRYPTO_COMPARE_KEY")
	base       = "https://min-api.cryptocompare.com"
	path       = "data/histominute"
	prefix     = "%s/raw/crypto-api-data/crypto-compare/minute-data/%s/%s_%s"
	stage      = os.Getenv("STAGE")
	t          = time.Now().Format("2006-01-02-150405")
	filename   = "/tmp/response.json"
	currencies = []string{"BTC", "ETH", "XRP", "LTC"}
	limit      = getEnv("CC_RESPONSE_LIMIT", "15")
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) {

	for _, v := range currencies {
		queryStringParams := buildQueryStringParams(v, limit)
		apiURL := buildURL(base, path, queryStringParams)

		data := apiResponseToStruct(apiURL, cryptoData{})
		data.FromSymbol = queryStringParams["fsym"]
		data.ToSymbol = queryStringParams["tsym"]
		log.Println("Writting response to file...")
		writeToJSON(data, filename)

		file, err := os.Open(filename)
		if err != nil {
			log.Fatal("Failed to open file ", err.Error())
		}
		defer file.Close()

		// Set up S3 connection
		config := aws.Config{
			Region: aws.String(region),
			// Credentials: credentials.NewSharedCredentials("", "personal"),
		}
		sess := session.New(&config)
		svc := s3manager.NewUploader(sess)

		log.Println("Uploading file to S3...")
		s3Prefix := formatS3Prefix(prefix, stage, data.FromSymbol, t, filename)
		log.Println(s3Prefix)
		result, err := svc.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(s3Prefix),
			Body:   file,
		})
		if err != nil {
			log.Fatal("Error uploading to S3:", err.Error())
		}
		log.Printf("Successfully uploaded %s to %s\n", filepath.Base(filename), result.Location)
	}

}

func main() {
	lambda.Start(Handler)
}

func buildURL(base, path string, queryParams map[string]string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		log.Fatal("Malformed url:", err.Error())
	}

	// Add endpoint
	baseURL.Path += path

	// Add query string parameters
	params := url.Values{}
	for key, val := range queryParams {
		params.Add(key, val)
	}

	// Add params to the url
	baseURL.RawQuery = params.Encode()

	return baseURL.String()
}

func apiResponseToStruct(url string, payload cryptoData) cryptoData {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error getting a response ", err.Error())
	}

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding response ", err.Error())
	}

	return payload
}

func writeToJSON(payload cryptoData, filepath string) {
	file, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		log.Fatal("Error marshaling to json ", err.Error())
	}

	err = ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		log.Fatal("Error writing to file ", err.Error())
	}
}

func formatS3Prefix(prefix, stage, fsym, t, filename string) string {
	return fmt.Sprintf(prefix, stage, fsym, t, filepath.Base(filename))
}

func buildQueryStringParams(fromSymbol, limit string) map[string]string {
	// Define query string params
	queryStringParams := make(map[string]string)
	queryStringParams["fsym"] = fromSymbol
	queryStringParams["tsym"] = "USD"
	queryStringParams["limit"] = limit
	queryStringParams["api_key"] = apiKey

	return queryStringParams
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
