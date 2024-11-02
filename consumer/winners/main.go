package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaBroker = "my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092" // Kafka broker
	topic       = "winners"                                                 // Tópico a consumir
	redisAddr   = "my-release-redis-master.default:6379"                    // Dirección del servidor Redis
	redisClient *redis.Client
)

func initKafkaReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topic,
		GroupID:  "my-consumer-group-winners",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}

func extractFaculty(message string) string {
	parts := strings.Split(message, ",")
	for _, part := range parts {
		if strings.Contains(part, "Faculty") {
			return strings.TrimSpace(strings.Split(part, ":")[1]) // Extract the faculty name
		}
	}
	return ""
}

func extractDiscipline(message string) string {
	parts := strings.Split(message, ",")
	for _, part := range parts {
		if strings.Contains(part, "Discipline") {
			return strings.TrimSpace(strings.Split(part, ":")[1]) // Extract the discipline
		}
	}
	return ""
}

// Map discipline ID to name
func getDisciplineName(disciplineID string) string {
	switch disciplineID {
	case "1":
		return "Natacion"
	case "2":
		return "Atletismo"
	case "3":
		return "Boxeo"
	default:
		return ""
	}
}

func consumeFromKafka(reader *kafka.Reader) {
	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		msg := string(message.Value)
		log.Printf("Received message: %s", msg)

		// Extraer la facultad del mensaje
		faculty := extractFaculty(msg)
		if faculty == "" {
			log.Printf("Faculty not found in message: %s", msg)
			continue
		}

		// Incrementar el contador total de estudiantes de la facultad
		_, err = redisClient.HIncrBy(context.Background(), "Total-Estudiantes", faculty, 1).Result()
		if err != nil {
			log.Printf("Failed to increment student counter for faculty %s: %v", faculty, err)
		} else {
			log.Printf("Student counter for faculty %s incremented successfully", faculty)
		}

		// Extraer la disciplina del mensaje
		disciplineID := extractDiscipline(msg)
		if disciplineID == "" {
			log.Printf("Discipline not found in message: %s", msg)
			continue
		}

		disciplineName := getDisciplineName(disciplineID)
		if disciplineName == "" {
			log.Printf("Unknown discipline ID %s in message: %s", disciplineID, msg)
			continue
		}

		_, err = redisClient.HIncrBy(context.Background(), "Ganadores-Disciplinas", disciplineName, 1).Result()
		if err != nil {
			log.Printf("Failed to increment winners counter for discipline %s: %v", disciplineName, err)
		} else {
			log.Printf("Winners counter for discipline %s incremented successfully", disciplineName)
		}

		if err := reader.CommitMessages(context.Background(), message); err != nil {
			log.Printf("Failed to commit message: %v", err)
		} else {
			log.Println("Message committed successfully")
		}
	}
}

func main() {
	flag.Parse()

	initRedisClient()
	defer redisClient.Close()

	reader := initKafkaReader()
	defer reader.Close()

	log.Printf("Kafka consumer for winners started, listening to topic %s", topic)

	// Consumir mensajes en un bucle
	consumeFromKafka(reader)
}
