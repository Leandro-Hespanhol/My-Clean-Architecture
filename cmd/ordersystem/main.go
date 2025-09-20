package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"MyCleanArchitecture/configs"
	"MyCleanArchitecture/internal/event"

	"MyCleanArchitecture/internal/infra/database"
	"MyCleanArchitecture/internal/infra/grpc/pb"
	"MyCleanArchitecture/internal/infra/grpc/service"
	"MyCleanArchitecture/internal/infra/web"
	"MyCleanArchitecture/internal/infra/web/webserver"
	"MyCleanArchitecture/internal/usecase"
	"MyCleanArchitecture/pkg/events"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	fmt.Println("Successfully connected to database")

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

	// gRPC Server
	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// GraphQL Server (simplified implementation)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>GraphQL Playground</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>GraphQL API</h1>
        <p>Use POST requests to /query endpoint with GraphQL queries</p>
        
        <h2>Available Operations:</h2>
        
        <div class="endpoint">
            <h3>Create Order (Mutation)</h3>
            <pre>
mutation {
  createOrder(input: {
    id: "123e4567-e89b-12d3-a456-426614174000"
    price: 100.50
    tax: 10.05
  }) {
    id
    price
    tax
    finalPrice
  }
}
            </pre>
        </div>
        
        <div class="endpoint">
            <h3>List Orders (Query)</h3>
            <pre>
query {
  orders {
    id
    price
    tax
    finalPrice
  }
}
            </pre>
        </div>
    </div>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Simple GraphQL query parsing
		if req.Query == "query { orders { id price tax finalPrice } }" {
			// List orders
			output, err := listOrdersUseCase.Execute()
			if err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"errors": []map[string]interface{}{
						{"message": err.Error()},
					},
				})
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"orders": output,
				},
			})
		} else if req.Query[:8] == "mutation" && req.Variables != nil {
			// Create order mutation
			if input, ok := req.Variables["input"].(map[string]interface{}); ok {
				dto := usecase.OrderInputDTO{
					ID:    input["id"].(string),
					Price: input["price"].(float64),
					Tax:   input["tax"].(float64),
				}

				output, err := createOrderUseCase.Execute(dto)
				if err != nil {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"errors": []map[string]interface{}{
							{"message": err.Error()},
						},
					})
					return
				}

				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": map[string]interface{}{
						"createOrder": output,
					},
				})
			} else {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"errors": []map[string]interface{}{
						{"message": "Invalid input variables"},
					},
				})
			}
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": []map[string]interface{}{
					{"message": "Query not supported"},
				},
			})
		}
	})

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	fmt.Println("GraphQL Playground available at http://localhost:" + configs.GraphQLServerPort + "/")
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
