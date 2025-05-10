package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/configs"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/event"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/event/handler"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/infra/database"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/infra/graph"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/migrations"

	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/infra/grpc/pb"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/infra/grpc/service"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/infra/web"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/infra/web/webserver"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/usecase"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig("../..")
	if err != nil {
		panic(err)
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName)
	fmt.Println("Connecting to database")
	db, err := sql.Open(configs.DBDriver, connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if !configs.SkipMigrations {
		fmt.Println("Running migrations...")
		migrationFolder := "internal/infra/database/migrations"
		err := migrations.MigrateUp(db, migrationFolder)
		if err != nil {
			panic(err)
		}
	}

	rabbitMQChannel := getRabbitMQChannel(configs.RabbitMQUser, configs.RabbitMQPass, configs.RabbitMQHost, configs.RabbitMQPort)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	repository := database.NewOrderRepository(db)
	orderCreateEvent := event.NewOrderCreated()

	useCaseContainer := usecase.NewOrderContainer(repository, eventDispatcher, orderCreateEvent)

	// webserver config
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(*useCaseContainer)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/orders", webOrderHandler.Get)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*useCaseContainer)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// graphql server config
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			OrderContainer: *useCaseContainer,
		}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(user, password, host, port string) *amqp.Channel {
	fmt.Println("Connecting to RabbitMQ")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
