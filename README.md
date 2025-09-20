# MyCleanArchitecture - Order Management System

Este projeto implementa um sistema de gerenciamento de pedidos seguindo os princípios da Clean Architecture, baseado nos projetos 14-gRPC e 20-CleanArch.

## Funcionalidades

O sistema oferece as seguintes funcionalidades:

- **Criar Order**: Criar um novo pedido
- **Listar Orders**: Listar todos os pedidos existentes

## APIs Disponíveis

### 1. REST API
- **Porta**: 8000
- **Endpoints**:
  - `POST /order` - Criar um novo pedido
  - `GET /order` - Listar todos os pedidos

### 2. gRPC Service
- **Porta**: 50051
- **Services**:
  - `CreateOrder` - Criar um novo pedido
  - `ListOrders` - Listar todos os pedidos

### 3. GraphQL
- **Porta**: 8080
- **Endpoints**:
  - `POST /query` - Endpoint GraphQL
  - `GET /` - GraphQL Playground
- **Operations**:
  - `Mutation createOrder` - Criar um novo pedido
  - `Query orders` - Listar todos os pedidos

## Estrutura do Projeto

```
MyCleanArchitecture/
├── cmd/
│   └── ordersystem/
│       └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── entity/                     # Entidades de domínio
│   │   ├── order.go
│   │   └── interface.go
│   ├── usecase/                    # Casos de uso
│   │   ├── create_order.go
│   │   └── list_orders.go
│   ├── infra/                      # Camada de infraestrutura
│   │   ├── database/
│   │   │   └── order_repository.go
│   │   ├── web/
│   │   │   ├── order_handler.go
│   │   │   └── webserver/
│   │   ├── grpc/
│   │   │   ├── pb/                 # Arquivos gerados do protobuf
│   │   │   └── service/
│   │   └── graph/                  # GraphQL
│   │       ├── schema.graphqls
│   │       ├── resolver.go
│   │       └── schema.resolvers.go
│   └── event/                      # Eventos de domínio
│       ├── order_created.go
│       └── handler/
├── pkg/
│   └── events/                     # Event dispatcher
├── configs/
│   └── config.go                   # Configurações
├── migrations/
│   └── 001_create_orders_table.sql # Migração do banco
├── api/
│   └── orders.http                 # Exemplos de requisições
├── proto/
│   └── order.proto                 # Definição do protobuf
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Pré-requisitos

- Go 1.21+
- Docker e Docker Compose
- MySQL 8.0

## Como Executar

### Opção 1: Com Docker (Recomendado)

1. Clone o repositório:
```bash
git clone <repository-url>
cd MyCleanArchitecture
```

2. Execute o Docker Compose:
```bash
docker compose up --build
```

Isso irá:
- Subir o banco de dados MySQL na porta 3306
- Executar as migrações automaticamente
- Subir a aplicação com todos os serviços

### Opção 2: Execução Local

1. Certifique-se de que o MySQL está rodando localmente

2. Configure as variáveis de ambiente no arquivo `app_config.env`:
```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders
WEB_SERVER_PORT=8000
GRPC_SERVER_PORT=50051
GRAPHQL_SERVER_PORT=8080
```

3. Execute as migrações no MySQL:
```sql
CREATE DATABASE IF NOT EXISTS orders;
USE orders;
SOURCE migrations/001_create_orders_table.sql;
```

4. Instale as dependências:
```bash
go mod tidy
```

5. Execute a aplicação:
```bash
go run cmd/ordersystem/main.go
```

## Testando a Aplicação

### Usando o arquivo api/orders.http

O projeto inclui um arquivo `api/orders.http` com exemplos de requisições para todas as APIs. Você pode usar extensões como REST Client no VS Code para executar essas requisições.

### GraphQL Playground

O GraphQL está totalmente implementado e funcional. Acesse o GraphQL Playground em http://localhost:8080/ para testar as queries e mutations interativamente.

### Exemplos de Uso

#### REST API

**Criar Order:**
```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "price": 100.50,
    "tax": 10.05
  }'
```

**Listar Orders:**
```bash
curl http://localhost:8000/order
```

#### GraphQL

Acesse o GraphQL Playground em: http://localhost:8080

**Criar Order:**
```graphql
mutation CreateOrder($input: OrderInput!) {
  createOrder(input: $input) {
    id
    price
    tax
    finalPrice
  }
}
```

Variables:
```json
{
  "input": {
    "id": "123e4567-e89b-12d3-a456-426614174001",
    "price": 200.75,
    "tax": 20.08
  }
}
```

**Listar Orders:**
```graphql
query {
  orders {
    id
    price
    tax
    finalPrice
  }
}
```

#### gRPC

Para testar o gRPC, você pode usar ferramentas como `grpcurl` ou `BloomRPC`:

```bash
# Listar serviços disponíveis
grpcurl -plaintext localhost:50051 list

# Criar Order
grpcurl -plaintext -d '{
  "id": "123e4567-e89b-12d3-a456-426614174002",
  "price": 150.25,
  "tax": 15.03
}' localhost:50051 pb.OrderService/CreateOrder

# Listar Orders
grpcurl -plaintext -d '{}' localhost:50051 pb.OrderService/ListOrders
```

## Arquitetura

O projeto segue os princípios da Clean Architecture:

- **Entities**: Regras de negócio fundamentais (`internal/entity`)
- **Use Cases**: Regras de negócio específicas da aplicação (`internal/usecase`)
- **Interface Adapters**: Conversores entre casos de uso e frameworks (`internal/infra`)
- **Frameworks & Drivers**: Detalhes externos como banco de dados, web, etc.

## Tecnologias Utilizadas

- **Go**: Linguagem de programação
- **MySQL**: Banco de dados
- **gRPC**: Comunicação entre serviços
- **GraphQL**: API flexível para consultas
- **Chi**: Router HTTP
- **GORM**: ORM para Go
- **Docker**: Containerização
- **RabbitMQ**: Message broker para eventos
- **Viper**: Gerenciamento de configurações

## Padrões Implementados

- **Clean Architecture**: Separação clara de responsabilidades
- **Repository Pattern**: Abstração da camada de dados
- **Event-Driven Architecture**: Eventos de domínio
- **Dependency Injection**: Inversão de dependências
- **CQRS**: Separação entre comandos e consultas

## Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request
