package parse

import (
	"fmt"
	"log"
	"os"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	awsCreds "github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/rcrowley/go-metrics"
	"github.com/sclasen/go-metrics-cloudwatch/config"
	"github.com/sclasen/go-metrics-cloudwatch/reporter"
)

var timers map[string]metrics.Timer
var counters map[string]metrics.Meter

func updateTimer(name string, ts time.Time) {
	if timers == nil {
		timers = map[string]metrics.Timer{}
	}

	if timers[name] == nil {
		timers[name] = metrics.NewTimer()
		metrics.GetOrRegister(fmt.Sprintf("parse-%s-timer", name), timers[name]) //todo change this string naming
	}
	timers[name].UpdateSince(ts)
}

func incrementCounter(name string, value int64) {
	if timers == nil {
		counters = map[string]metrics.Meter{}
	}

	if counters[name] == nil {
		counters[name] = metrics.NewMeter()
		metrics.GetOrRegister(fmt.Sprintf("parse-%s-counter", name), counters[name]) //todo change this string naming
	}
	counters[name].Mark(value)

}

// SetupMetricExternalLogging will setup sending the counter to cloudwatch
func SetupMetricExternalLogging(cloudwatchAccessKey string, cloudwatchAccessSecret string, namespace string) {

	creds := awsCreds.NewStaticCredentials(cloudwatchAccessKey, cloudwatchAccessSecret, "" /*token empty*/)
	awsConfig := awssdk.NewConfig()
	awsConfig.Credentials = creds
	awsConfig.Region = awssdk.String("us-east-1")
	sess, err := awssession.NewSession(awsConfig)
	if err != nil {
		fmt.Printf("failed to create session, %v", err)
		return
	}

	metricsConf := &config.Config{
		Client:            cloudwatch.New(sess),
		Namespace:         namespace,
		Filter:            &config.NoFilter{},
		ReportingInterval: 1 * time.Minute,
		StaticDimensions:  map[string]string{"name": "value"},
	}
	go reporter.Cloudwatch(metrics.DefaultRegistry, metricsConf)
}

// SetupMetricFileLogging will setup sending the counter to a file
func SetupMetricFileLogging(metricFilePath string) {
	if metricFilePath != "" {
		var metricsLog *os.File
		metricsLog, err := os.OpenFile(metricFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			fmt.Printf("error open metric file, no metrics will be logged.")
			fmt.Printf(err.Error())
		} else {
			go metrics.Log(metrics.DefaultRegistry, 1*time.Minute, log.New(metricsLog, "metrics: ", log.Lmicroseconds))
		}

	}
}