package project

import (
	"context"
	"fmt"
	"github.com/devpies/saas-core/internal/project/res"
	"github.com/nats-io/nats.go"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/devpies/saas-core/internal/project/config"
	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/handler"
	"github.com/devpies/saas-core/internal/project/repository"
	"github.com/devpies/saas-core/internal/project/service"
	"github.com/devpies/saas-core/pkg/log"
	"github.com/devpies/saas-core/pkg/msg"

	"github.com/ardanlabs/conf"
	"go.uber.org/zap"
)

// Run contains the app setup.
func Run() error {
	var (
		cfg     config.Config
		logger  *zap.Logger
		logPath = "log/out.log"
		err     error
	)

	if err = conf.Parse(os.Args[1:], "PROJECT", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("PROJECT", &cfg)
			if err != nil {
				logger.Error("error generating config usage", zap.Error(err))
				return err
			}
			fmt.Println(usage)
			return nil
		}
		logger.Error("error parsing config", zap.Error(err))
		return err
	}

	if cfg.Web.Production {
		logger, err = log.NewProductionLogger(logPath)
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	defer logger.Sync()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	pg, Close, err := db.NewPostgresDatabase(logger, cfg)
	if err != nil {
		return err
	}
	defer Close()

	// Execute latest migration in production.
	if cfg.Web.Production {
		if err = res.MigrateUp(pg.URL.String()); err != nil {
			logger.Error("error connecting to user database", zap.Error(err))
			return err
		}
	}

	jetStream := msg.NewStreamContext(logger, shutdown, cfg.Nats.Address, cfg.Nats.Port)

	_ = jetStream.Create(msg.StreamProjects)

	// Initialize 3-layered architecture.
	taskRepo := repository.NewTaskRepository(logger, pg)
	columnRepo := repository.NewColumnRepository(logger, pg)
	projectRepo := repository.NewProjectRepository(logger, pg)
	membershipRepo := repository.NewMembershipRepository(logger, pg)

	taskService := service.NewTaskService(logger, taskRepo)
	columnService := service.NewColumnService(logger, columnRepo)
	projectService := service.NewProjectService(logger, projectRepo, membershipRepo)
	membershipService := service.NewMembershipService(logger, membershipRepo)

	taskHandler := handler.NewTaskHandler(logger, taskService, columnService)
	columnHandler := handler.NewColumnHandler(logger, columnService)
	projectHandler := handler.NewProjectHandler(logger, jetStream, projectService, columnService, taskService)

	opts := []nats.SubOpt{nats.DeliverAll(), nats.ManualAck()}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("listener panic: %v", r)
				logger.Error(fmt.Sprintf("%s", debug.Stack()), zap.Error(err))
			}
		}()

		// Listen to membership events to save a redundant copy in the database.
		jetStream.Listen(
			string(msg.TypeMembershipCreated),
			msg.SubjectMembershipCreated,
			"membership_created_consumer",
			membershipService.CreateMembershipCopyFromEvent,
			opts...,
		)
		jetStream.Listen(
			string(msg.TypeMembershipUpdated),
			msg.SubjectMembershipUpdated,
			"membership_updated_consumer",
			membershipService.UpdateMembershipCopyFromEvent,
			opts...,
		)
		jetStream.Listen(
			string(msg.TypeMembershipDeleted),
			msg.SubjectMembershipDeleted,
			"membership_deleted_consumer",
			membershipService.DeleteMembershipCopyFromEvent,
			opts...,
		)
		jetStream.Listen(
			string(msg.TypeTeamAssignedEventType),
			msg.SubjectProjectTeamAssigned,
			"project_assigned_consumer",
			projectService.AssignProjectTeamFromEvent,
			opts...,
		)
	}()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, taskHandler, columnHandler, projectHandler, cfg),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting project service on %s:%s", cfg.Web.Address, cfg.Web.Port))
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err = <-serverErrors:
		logger.Error("error on startup", zap.Error(err))
		return err
	case sig := <-shutdown:
		logger.Info(fmt.Sprintf("Start shutdown due to %s signal", sig))

		// Give on going tasks a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			err = srv.Close()
		}

		switch {
		case sig == syscall.SIGSTOP:
			logger.Error("error on integrity issue caused shutdown", zap.Error(err))
			return err
		case err != nil:
			logger.Error("error on gracefully shutdown", zap.Error(err))
			return err
		}
	}

	return err
}
