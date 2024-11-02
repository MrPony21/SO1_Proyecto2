package main

import (
	"context"

	pb "go-client/proto"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Aqui elegimos los servidores y su diciplina
var (
	servers = map[int]string{
		1: "natacion-service:50051",
		2: "atletismo-service:50052",
		3: "boxeo-service:50053",
	}
)

type Student struct {
	Student    string `json:"student"`
	Age        int    `json:"age"`
	Faculty    string `json:"faculty"`
	Discipline int    `json:"discipline"`
}

func sendData(fiberCtx *fiber.Ctx) error {
	var body Student
	if err := fiberCtx.BodyParser(&body); err != nil {
		return fiberCtx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Verificamos si la disciplina está en el mapa de servidores
	addr, ok := servers[body.Discipline]
	if !ok {
		return fiberCtx.Status(400).JSON(fiber.Map{
			"error": "discipline debe ser 1, 2 o 3 para enviar al servidor gRPC",
		})
	}

	// Establecer conexión al servidor
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStudentClient(conn)

	// Crear un canal para recibir la respuesta y el error
	responseChan := make(chan *pb.StudentResponse)
	errorChan := make(chan error)
	go func() {
		// Contactar al servidor y obtener su respuesta
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.GetStudent(ctx, &pb.StudentRequest{
			Student:    body.Student,
			Age:        int32(body.Age),
			Faculty:    body.Faculty,
			Discipline: pb.Discipline(body.Discipline),
		})

		if err != nil {
			errorChan <- err
			return
		}

		responseChan <- r
	}()

	select {
	case response := <-responseChan:
		return fiberCtx.JSON(fiber.Map{
			"message": response.GetSuccess(),
		})
	case err := <-errorChan:
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	case <-time.After(5 * time.Second):
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": "timeout",
		})
	}
}

func main() {
	app := fiber.New()
	app.Post("/Agronomia", sendData)

	err := app.Listen(":8080")
	if err != nil {
		log.Println(err)
		return
	}
}
