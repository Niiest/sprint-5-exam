package main

// consumer_example implements a consumer using the non-channel Poll() API
// to retrieve messages and events.

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
)

type KafkaConfiguration struct {
	BootstrapServers []string `json:"bootstrap_servers"`
	Topic            string   `json:"topic"`
	User             string   `json:"user"`
	Password         string   `json:"password"`
}

type SchemaRegistryConfiguration struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Configuration struct {
	Kafka          KafkaConfiguration          `json:"kafka"`
	SchemaRegistry SchemaRegistryConfiguration `json:"schema_registry"`
}

func main() {

	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Printf("error:%s", err)
		os.Exit(1)
	}

	group := "group"
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(configuration.Kafka.BootstrapServers, ","),
		"group.id":           group,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
		"security.protocol":  "SASL_SSL",
		"ssl.ca.location":    "../../YandexInternalRootCA.crt",
		"sasl.mechanism":     "SCRAM-SHA-512",
		"sasl.username":      configuration.Kafka.User,
		"sasl.password":      configuration.Kafka.Password,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	schemaRegistryConfig := schemaregistry.NewConfigWithBasicAuthentication(
		configuration.SchemaRegistry.Url,
		configuration.SchemaRegistry.User,
		configuration.SchemaRegistry.Password)

	client, err := schemaregistry.NewClient(schemaRegistryConfig)

	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	deser, err := jsonschema.NewDeserializer(client, serde.ValueSerde, jsonschema.NewDeserializerConfig())

	if err != nil {
		fmt.Printf("Failed to create deserializer: %s\n", err)
		os.Exit(1)
	}

	err = c.SubscribeTopics([]string{configuration.Kafka.Topic}, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				value := User{}
				err := deser.DeserializeInto(*e.TopicPartition.Topic, e.Value, &value)
				if err != nil {
					fmt.Printf("Failed to deserialize payload: %s\n", err)
				} else {
					fmt.Printf("%% Message on %s:\n%+v\n", e.TopicPartition, value)
				}
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}

// User is a simple record example
type User struct {
	Name           string `json:"name"`
	FavoriteNumber int64  `json:"favorite_number"`
	FavoriteColor  string `json:"favorite_color"`
}
