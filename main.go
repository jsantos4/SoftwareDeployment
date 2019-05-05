package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    cron "github.com/robfig/cron"
    loggly "github.com/jamespearly/loggly"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "os"
    "os/signal"
    "time"
)

type Person struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Gender  string `json:"gender"`
	Region  string `json:"region"`
}


type PersonName struct {
  Name string
  Surname string
  Gender string
  Region string
  TimeStamp string
}

func apiRequest() *http.Response {
    url := "https://uinames.com/api/"


   req, err := http.NewRequest("GET", url, nil)

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
       fmt.Println("Error on response.\n[ERRO] -", err)
   }

   return resp
}

func printStats(Person *Person) {
    fmt.Println("Stats")
    fmt.Println("Name: ", Person.Name, Person.Surname)
    fmt.Println("Region:", Person.Region)

}


func run() {

    response := apiRequest()

    defer response.Body.Close()

    byteValue, _ := ioutil.ReadAll(response.Body)

    var stats Person

    json.Unmarshal(byteValue, &stats)

    client := loggly.New("Name")

    logMessage := "Full name: " + stats.Name + stats.Surname
    logContent := client.Send("info", logMessage)

    fmt.Println(logContent)

    ses, err := session.NewSession(&aws.Config{
	Region: aws.String("us-east-1")},
    )
    if err != nil {
	fmt.Println("Error creating session:")
	fmt.Println(err.Error())
	os.Exit(1)
    }

    // Create DynamoDB client
    svc := dynamodb.New(ses)

		  res := stats

	e := PersonName{
        Name: res.Name,
        Surname: res.Surname,
        Gender: res.Gender,
        Region: res.Region,
        TimeStamp: time.Now().Format("2019-01-01 12:05:28"),
	  }

    av, err := dynamodbattribute.MarshalMap(e)

    input := &dynamodb.PutItemInput{
	Item:      av,
	TableName: aws.String("JinYang"),
    }

    _, err = svc.PutItem(input)

    if err != nil {
	fmt.Println("Got error calling PutItem:")
	fmt.Println(err.Error())
	os.Exit(1)
    }

    fmt.Println("success")

}

func main() {
	c := cron.New()
	c.AddFunc("@every 1h", func() { run() })
	c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
