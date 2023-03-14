package db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"emperror.dev/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

// CreateInstanceAndDatabase is a tools that, well, create instance and database on cloud spanner
// Warning: local and test development only
// nolint
func CreateInstanceAndDatabase(ctx context.Context, projectID, instanceID, databaseID string) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	logger.Debug(fmt.Sprintf("creating instance projects/%s/instances/%s", projectID, instanceID))
	if err := createInstance(ctx, projectID, instanceID); err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("creating database projects/%s/instances/%s/databases/%s", projectID, instanceID, databaseID))
	if err := createDatabase(ctx, projectID, instanceID, databaseID); err != nil {
		return err
	}

	return nil
}

// StartEmulator will start spanner emulator and set SPANNER_EMULATOR_HOST environment
// run this as SetupSuite on all repository tests. No tear down required, testcontainers handles that
func StartEmulator() error {
	// if SPANNER_EMULATOR_HOST already been set
	// it means there's an emulator ready
	if os.Getenv("SPANNER_EMULATOR_HOST") != "" {
		return nil
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "gcr.io/cloud-spanner-emulator/emulator:latest",
		ExposedPorts: []string{"9010/tcp", "9020/tcp"},
		WaitingFor:   wait.ForLog("gRPC server listening at"),
	}
	emu, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return errors.WrapIf(err, "failed to start container")
	}

	host, err := emu.Host(ctx)
	if err != nil {
		return errors.WrapIf(err, "failed to get host")
	}

	mappedPort, err := emu.MappedPort(ctx, "9010")
	if err != nil {
		return errors.WrapIf(err, "failed to get mappedPort")
	}

	uri := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())
	_ = os.Setenv("SPANNER_EMULATOR_HOST", strings.TrimPrefix(uri, "http://"))

	return nil
}

// nolint
var dbCount struct {
	sync.Mutex
	count int
}

func SetupTestTable(tablePrefixName string) (*spanner.Client, error) {
	projectID := "test-project-id"
	instanceID := "test-instance-id"
	ctx := context.Background()

	// ensure that each test is using different database. no cleanup required
	dbCount.Lock()
	defer dbCount.Unlock()
	dbCount.count++
	databaseID := fmt.Sprintf("test_%s_%d", tablePrefixName, dbCount.count)

	err := CreateInstanceAndDatabase(ctx, projectID, instanceID, databaseID)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to CreateInstanceAndDatabase")
	}

	err = Migrate(projectID, instanceID, databaseID)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to migrate")
	}

	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, databaseID)
	db, err := spanner.NewClient(ctx, dsn)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to start spanner client")
	}

	return db, nil
}

// nolint
func createInstance(ctx context.Context, projectID, instanceID string) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	instanceAdminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}

	instanceName := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	_, err = instanceAdminClient.GetInstance(ctx, &instancepb.GetInstanceRequest{
		Name: instanceName,
	})

	if err != nil && spanner.ErrCode(err) != codes.NotFound {
		return err
	}

	if err == nil {
		logger.Debug(fmt.Sprintf("instance %s already created, skipping...", instanceName))

		return nil
	}

	_, err = instanceAdminClient.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
	})

	if err != nil && spanner.ErrCode(err) != codes.AlreadyExists {
		return err
	}

	return nil
}

// nolint
func createDatabase(ctx context.Context, projectID, instanceID, databaseID string) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	databaseAdminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}

	databaseName := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, databaseID)
	_, err = databaseAdminClient.GetDatabase(ctx, &databasepb.GetDatabaseRequest{Name: databaseName})

	if err != nil && spanner.ErrCode(err) != codes.NotFound {
		return err
	}

	if err == nil {
		logger.Debug(fmt.Sprintf("database %s already created, skipping...", databaseName))

		return nil
	}

	op, err := databaseAdminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", databaseID),
	})
	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("waiting for database %s to be created", databaseName))

	if _, err = op.Wait(ctx); err != nil {
		return err
	}

	return err
}
