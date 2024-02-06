package app

import (
	"context"
	"fmt"
	"log"
	grpcside "merger-adapter/internal/api/grpc_merger_client"
	"merger-adapter/internal/api/telegrambot"
	"merger-adapter/internal/component/sqlite"
	"merger-adapter/internal/repository/messages_repository"
	"merger-adapter/internal/service/blobstore"
	"time"
)

func Run(ctx context.Context, cfg *Config) error {
	app, errCh := newApplication(ctx)

	// blobstore
	files, err := blobstore.InitRedis(blobstore.Config{
		FilesLifetime: time.Minute * 10,
		RedisUrl:      cfg.RedisUrl,
	})
	if err != nil {
		return fmt.Errorf("init redis as blobstore: %s", err)
	}

	// init db:
	db, err := sqlite.InitSqlite(sqlite.Config{
		DataSourceName: cfg.DbFile,
	})
	if err != nil {
		return fmt.Errorf("init databse: %s", err)
	}
	log.Println("DatabaseInitialized")

	repo := messages_repository.NewMessagesRepositorySqlite(db)
	log.Println("messages repository created")

	server, err := grpcside.InitGrpcMergerClient(grpcside.Config{
		Host: cfg.ServerHost,
	})
	if err != nil {
		return fmt.Errorf("init grpc merger server: %s", err)
	}
	log.Println("InitGrpcMergerClientInitialized")

	tgbClient, err := telegrambot.InitClient(telegrambot.Deps{
		Token:  cfg.TgBotToken,
		ApiKey: cfg.TgXApiKey,
		ChatID: cfg.TgChat,
		Server: server,
		Files:  files,
		Repo:   repo,
	})
	if err != nil {
		return fmt.Errorf("tg client initialization: %s", err)
	}
	log.Println("TelegramAdapterInitialized")
	go app.run(tgbClient, "telegram adapter")

	log.Println("Application start is over, waiting when ctx done")

	return app.gracefulShutdownApplication(errCh)
}

func (a *application) gracefulShutdownApplication(errCh <-chan error) error {
	var err error
	select {
	case <-a.ctx.Done():
		log.Println("Application receive ctx.Done signal")
	case err = <-errCh:
		a.cancelFunc()
		log.Println("ApplicationReceiveInternalError")
	}
	a.wg.Wait()
	return err
}

func (a *application) run(r Runnable, name string) {
	a.wg.Add(1)
	defer a.wg.Done()
	log.Println("Run Runnable", name)
	err := r.Run(a.ctx)
	if err != nil {
		a.errorf("%s run: %s", name, err)
	}
	log.Println("Stopped Runnable", name)
}

func (a *application) errorf(format string, args ...any) {
	select {
	case a.errCh <- fmt.Errorf(format, args...):
	default:
	}
}

type Runnable interface {
	Run(ctx context.Context) error
}
