package main

import (
	"context"
	"fmt"
	"github.com/sukhajata/devicetwin.git/api"
	"github.com/sukhajata/devicetwin.git/internal/dbclient/nosql"
	"github.com/sukhajata/devicetwin.git/internal/dbclient/sql"
	"github.com/sukhajata/devicetwin.git/internal/messageprocessor"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sukhajata/devicetwin.git/internal/consistency"
	"github.com/sukhajata/devicetwin.git/internal/core"
	"github.com/sukhajata/devicetwin.git/internal/dataapi"
	"github.com/sukhajata/devicetwin.git/internal/dbclient"
	"github.com/sukhajata/devicetwin.git/pkg/authhelper"
	"github.com/sukhajata/devicetwin.git/pkg/db"
	"github.com/sukhajata/devicetwin.git/pkg/errorhelper"
	"github.com/sukhajata/devicetwin.git/pkg/grpchelper"
	"github.com/sukhajata/devicetwin.git/pkg/loggerhelper"
	"github.com/sukhajata/devicetwin.git/pkg/ppmqtt"
	pbAuth "github.com/sukhajata/ppauth"
	pb "github.com/sukhajata/ppconfig"
	pbConnection "github.com/sukhajata/ppconnection"
	pbLogger "github.com/sukhajata/pplogger"
	"github.com/sukhajata/ppmessage/ppdownlink"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	useCouchbase              = getEnv("useCouchbase", "false")
	mqttBroker                = getEnv("mqttBroker", "ssl://mosquitto:8883")
	mqttUsername              = getEnv("mqttUsername", "admin")
	mqttPassword              = getEnv("mqttPassword", "admin")
	mqttDownlinkTopic         = getEnv("mqttDownlinkTopic", "application/powerpilot/downlink/config")
	mqttUplinkTopic           = getEnv("mqttUplinkTopic", "$share/devicetwin/application/powerpilot/uplink/config/#")
	mqttConnectionUpdateTopic = getEnv("mqttConnectionsTopic", "$share/devicetwin/application/powerpilot/connections")

	couchbaseBucketName       = getEnv("couchbaseBucketName", "test")
	couchbaseBucketNameShared = getEnv("couchbaseBucketNameShared", "shared")
	couchbaseUsername         = getEnv("couchbaseUsername", "admin")
	couchbasePassword         = getEnv("couchbasePassword", "admin")
	couchbaseServerAddress    = getEnv("couchbaseServerAddress", "couchbase.dev.svc")

	dataServiceAddress      = getEnv("dataServiceAddress", "http://postgrest-api:3000")
	dataToken               = getEnv("dataToken", "")
	connectionsDataViewName = getEnv("connectionsDataViewName", "POWER_BI_CONNECTIONS_DATA")
	repeatCheckSchedule     = getEnv("repeatCheckSchedule", "15_30_45")
	psqlURL                 = getEnv("psqlURL", "postgresql://admin:admin@dev-timescale.dev.svc/test?sslmode=require")

	minutesRunConsistencyCheck = getEnv("minutesRunConsistencyCheck", "1440")
	configServicePort          = getEnv("configServicePort", "9090")
	authServiceAddress         = getEnv("authServiceAddress", "auth-service:9030")
	loggerServiceAddress       = getEnv("loggerServiceAddress", "logger-service:9031")
	connectionServiceAddress   = getEnv("connectionServiceAddress", "connection-service:9010")
	serviceKey                 = getEnv("serviceKey", "test")
	adminRole                  = getEnv("rolePowerpilotAdmin", "powerpilot-admin")
	installerRole              = getEnv("rolePowerpilotInstaller", "powerpilot-installer")
	superuserRole              = getEnv("rolePowerpilotSuperuser", "powerpilot-superuser")
	//receiveChan                = make(chan *ppuplink.ConfigUplinkMessage, 2)
	transmitChan = make(chan *ppdownlink.ConfigDownlinkMessage, 2)

	grpcAuthClient       pbAuth.AuthServiceClient
	grpcLoggerClient     pbLogger.LoggerServiceClient
	grpcConnectionClient pbConnection.ConnectionServiceClient
	dbClient             dbclient.Client
	mqttClient           *ppmqtt.PPClient
	configService        core.ConfigHandler
	loggerHelper         loggerhelper.Helper
	errorChan            chan *pbLogger.ErrorMessage
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func setupScheduledConsistencyCheck(consistencyService *consistency.Service) {
	mins, err := strconv.Atoi(minutesRunConsistencyCheck)
	if err != nil {
		mins = 1
	}
	ticker := time.NewTicker(time.Duration(mins) * time.Minute)

	for range ticker.C {
		consistencyService.RunScheduledConsistencyCheck()
	}

}

// PublishDownlink to mqtt
func PublishDownlink(downlink *ppdownlink.ConfigDownlinkMessage) error {
	bytes, err := proto.Marshal(downlink)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Publishing message deviceeui %v index %v slot %v value %v", downlink.Deviceeui, downlink.Index, downlink.Slot, downlink.Value)
	loggerhelper.WriteToLog(message)

	topic := fmt.Sprintf("%s/%s", mqttDownlinkTopic, downlink.Deviceeui)
	err = mqttClient.Publish(ppmqtt.Message{
		Topic:   topic,
		Payload: bytes,
	})

	return err
}

func connectMQTT() {
	var err error
	mqttClient, err = ppmqtt.NewClient(mqttBroker, mqttUsername, mqttPassword, "config-service")
	errorhelper.PanicOnError(err)

	// subscribe
	err = mqttClient.Subscribe(mqttUplinkTopic)
	errorhelper.PanicOnError(err)
	err = mqttClient.Subscribe(mqttConnectionUpdateTopic)
	errorhelper.PanicOnError(err)

	loggerhelper.WriteToLog("Connected to mqtt broker")

	// listen for mqtt errors and reconnect
	go func(errorChan chan error) {
		err := <-errorChan
		loggerhelper.WriteToLog(err.Error())
		connectMQTT()
	}(mqttClient.ErrorChan)

	// listen for mqtt messages and process
	messageProcessor := messageprocessor.NewMessageProcessor(configService, dbClient, errorChan)
	go func(messageChan chan ppmqtt.Message) {
		for msg := range messageChan {
			go messageProcessor.ProcessMessage(msg)
		}
	}(mqttClient.ReceiveChan)

	// listen for messages to send from internal services
	go func(downlinkChan chan *ppdownlink.ConfigDownlinkMessage) {
		for msg := range downlinkChan {
			err := PublishDownlink(msg)
			if err != nil {
				loggerHelper.LogError("listenDownlinks", err.Error(), pbLogger.ErrorMessage_FATAL)
				return
			}
		}
	}(transmitChan)

}

func sendError(msg *pbLogger.ErrorMessage) {
	ctx, cancel := authhelper.GetContextWithAuth(serviceKey)
	defer cancel()
	_, err := grpcLoggerClient.LogError(ctx, msg)
	if err != nil {
		loggerhelper.WriteToLog(err.Error())
	}
}

func sendDeviceEvent(msg *pbLogger.DeviceLogMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := grpcLoggerClient.LogDeviceEvent(ctx, msg)
	if err != nil {
		loggerHelper.LogError("SetDesired4", err.Error(), pbLogger.ErrorMessage_SEVERE)
	} else {
		loggerhelper.WriteToLog(fmt.Sprintf("Sent log report: %s", msg.Message))
	}
}

// main entry
func main() {
	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}

	// logger client
	conn, err := grpchelper.ConnectGRPC(loggerServiceAddress)
	errorhelper.PanicOnError(err)
	defer func() {
		err := conn.Close()
		if err != nil {
			loggerhelper.WriteToLog(err.Error())
		}
	}()
	grpcLoggerClient = pbLogger.NewLoggerServiceClient(conn)

	// helper
	loggerHelper = loggerhelper.NewHelper(grpcLoggerClient, "config-service")

	// auth client
	conn2, err := grpchelper.ConnectGRPC(authServiceAddress)
	errorhelper.PanicOnError(err)
	defer func() {
		err := conn2.Close()
		if err != nil {
			loggerHelper.LogError("CloseAuthClienConnection", err.Error(), pbLogger.ErrorMessage_SEVERE)
		}
	}()
	grpcAuthClient = pbAuth.NewAuthServiceClient(conn2)

	// connection service client
	conn3, err := grpchelper.ConnectGRPC(connectionServiceAddress)
	errorhelper.PanicOnError(err)
	defer func() {
		err := conn3.Close()
		if err != nil {
			loggerHelper.LogError("CloseConnectionClientConnection", err.Error(), pbLogger.ErrorMessage_SEVERE)
		}
	}()
	grpcConnectionClient = pbConnection.NewConnectionServiceClient(conn3)

	// error chan
	errorChan = make(chan *pbLogger.ErrorMessage, 3)
	go func(errorChan chan *pbLogger.ErrorMessage) {
		for msg := range errorChan {
			loggerhelper.WriteToLog(msg.Message)
			sendError(msg)
		}
	}(errorChan)

	// device event chan
	deviceEventChan := make(chan *pbLogger.DeviceLogMessage, 3)
	go func(deviceEventChan <-chan *pbLogger.DeviceLogMessage) {
		for msg := range deviceEventChan {
			sendDeviceEvent(msg)
		}
	}(deviceEventChan)

	// database connection
	if useCouchbase == "true" {
		dbEngine, err := db.NewCouchbaseEngine(couchbaseServerAddress, couchbaseUsername, couchbasePassword, couchbaseBucketName, couchbaseBucketNameShared)
		errorhelper.PanicOnError(err)
		dbClient = nosql.NewCouchbaseClient(dbEngine, couchbaseBucketName, couchbaseBucketNameShared, loggerHelper)
	} else {
		dbEngine, err := db.NewTimescaleEngine(psqlURL)
		errorhelper.PanicOnError(err)
		dbClient = sql.NewTimescaleClient(dbEngine, errorChan)
	}

	// data API
	dataAPIClient := dataapi.NewClient(dataServiceAddress, dataToken, connectionsDataViewName, &http.Client{})

	// consistency service
	consistencyService := consistency.NewService(dbClient, dataAPIClient, transmitChan, repeatCheckSchedule, loggerHelper)
	go setupScheduledConsistencyCheck(consistencyService)

	// config service
	configService = core.NewService(
		dbClient,
		grpcConnectionClient,
		grpcAuthClient,
		consistencyService,
		serviceKey,
		loggerHelper,
		transmitChan,
		errorChan,
		deviceEventChan,
		adminRole,
		installerRole,
		superuserRole,
	)

	// mqtt broker
	connectMQTT()

	// grpc server
	configServiceServer := api.NewGRPCConfigServer(configService, consistencyService, loggerHelper)

	// http server
	api.NewHTTPServer(configService)

	loggerhelper.WriteToLog("Connected to services")

	// setup up gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configServicePort))
	errorhelper.PanicOnError(err)
	/*var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			panic(err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}*/
	grpcServer := grpc.NewServer()
	pb.RegisterConfigServiceServer(grpcServer, configServiceServer)

	err = grpcServer.Serve(lis)
	errorhelper.PanicOnError(err)

}
