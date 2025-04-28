package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/google/uuid"
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

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(configuration.Kafka.BootstrapServers, ","),
		"security.protocol": "SASL_SSL",
		"ssl.ca.location":   "../../YandexInternalRootCA.crt",
		"sasl.mechanism":    "SCRAM-SHA-512",
		"sasl.username":     configuration.Kafka.User,
		"sasl.password":     configuration.Kafka.Password,
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	schemaRegistryConfig := schemaregistry.NewConfigWithBasicAuthentication(
		configuration.SchemaRegistry.Url,
		configuration.SchemaRegistry.User,
		configuration.SchemaRegistry.Password)

	client, err := schemaregistry.NewClient(schemaRegistryConfig)

	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	ser, err := jsonschema.NewSerializer(client, serde.ValueSerde, jsonschema.NewSerializerConfig())

	if err != nil {
		fmt.Printf("Failed to create serializer: %s\n", err)
		os.Exit(1)
	}

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)

	value := User{
		Name:           "First user",
		FavoriteNumber: 42,
		FavoriteColor:  "blue",
	}

	ser.Conf.UseLatestVersion = true
	ser.Conf.AutoRegisterSchemas = false
	payload, err := ser.Serialize(configuration.Kafka.Topic, &value)

	if err != nil {
		fmt.Printf("Failed to serialize payload: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Serialized\n")

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &configuration.Kafka.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(uuid.New().String()),
		Value:          payload,
		Headers:        []kafka.Header{},
	}, deliveryChan)
	if err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		os.Exit(1)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}

// User is a simple record example
type User struct {
	Name           string `json:"name"`
	FavoriteNumber int64  `json:"favorite_number"`
	FavoriteColor  string `json:"favorite_color"`
}
