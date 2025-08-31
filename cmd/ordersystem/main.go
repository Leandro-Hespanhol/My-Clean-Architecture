package main

import (
	"database/sql"
	"fmt"

	"net/http"

	"MyCleanArchitecture/configs"
	"MyCleanArchitecture/internal/event"

	"MyCleanArchitecture/internal/infra/database"
	// "MyCleanArchitecture/internal/infra/grpc/pb"
	// "MyCleanArchitecture/internal/infra/grpc/service"
	"MyCleanArchitecture/internal/infra/web"
	"MyCleanArchitecture/internal/infra/web/webserver"
	"MyCleanArchitecture/internal/usecase"
	"MyCleanArchitecture/pkg/events"

	"github.com/streadway/amqp"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// Using nil for now - you'll need to implement an in-memory repository or connect to a real database
	var db *sql.DB = nil

	// rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	// eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
	// 	RabbitMQChannel: rabbitMQChannel,
	// })

	orderRepository := database.NewOrderRepository(db)
	orderCreatedEvent := event.NewOrderCreated()

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderCreatedEvent, eventDispatcher)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	// Use the variables to avoid "declared and not used" errors
	_ = createOrderUseCase
	_ = listOrdersUseCase

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, orderCreatedEvent)
	webserver.AddHandler("/order", "POST", webOrderHandler.Create)
	webserver.AddHandler("/order", "GET", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	// gRPC Server (commented out due to protobuf issues)
	// grpcServer := grpc.NewServer()
	// orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	// pb.RegisterOrderServiceServer(grpcServer, orderService)
	// reflection.Register(grpcServer)

	fmt.Println("gRPC server disabled - fix protobuf files to enable")
	// fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	// if err != nil {
	// 	panic(err)
	// }
	// go grpcServer.Serve(lis)

	// GraphQL Server (simplified for demonstration)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GraphQL Playground - Use gqlgen to generate proper schema"))
	})
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GraphQL Query endpoint - Use gqlgen to generate proper resolvers"))
	})

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
