package app

import (
	"context"
	"fmt"
	"log"
	grpcside "merger-adapter/internal/api/grpc_merger_client"
	"merger-adapter/internal/api/telegrambot"
	"merger-adapter/internal/api/vkontaktebot"
	"merger-adapter/internal/common/msgs"
	"merger-adapter/internal/component/sqlite"
	"merger-adapter/internal/service/filestore"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/runnable"
)

func Run(ctx context.Context, cfg *Config) error {
	app, errCh := newApplication(ctx)

	// filestore
	files := filestore.NewLocal("files")

	// init db:
	db, err := sqlite.InitSqlite(sqlite.Config{
		DataSourceName: cfg.DbFile,
	})
	if err != nil {
		return fmt.Errorf("init databse: %s", err)
	}
	log.Println(msgs.DatabaseInitialized)

	messagesMap := kvstore.NewSqliteMessagesMap(db)
	log.Println(msgs.MessagesMapCreated)

	server, err := grpcside.InitGrpcMergerClient(grpcside.Config{
		Host: cfg.ServerHost,
	})
	if err != nil {
		return fmt.Errorf("init grpc merger server: %s", err)
	}
	log.Println(msgs.InitGrpcMergerClientInitialized)

	tgbClient, err := telegrambot.InitClient(telegrambot.Deps{
		Token:       cfg.TgBotToken,
		ApiKey:      cfg.TgXApiKey,
		ChatID:      cfg.TgChat,
		Server:      server,
		MessagesMap: messagesMap,
		Files:       files,
	})
	if err != nil {
		return fmt.Errorf("tg client initialization: %s", err)
	}
	log.Println(msgs.TelegramAdapterInitialized)
	go app.run(tgbClient, "vkontakte adapter")

	vkbClient, err := vkontaktebot.InitClient(vkontaktebot.Deps{
		Token:       cfg.VkBotToken,
		ApiKey:      cfg.VkXApiKey,
		PeerId:      cfg.VkPeer,
		Server:      server,
		MessagesMap: messagesMap,
		Files:       files,
	})
	if err != nil {
		return fmt.Errorf("vk client initialization: %s", err)
	}
	log.Println(msgs.VkontakteAdapterInitialized)

	go app.run(vkbClient, "vkontakte adapter")

	log.Println(msgs.ApplicationStarted)

	return app.gracefulShutdownApplication(errCh)
}

func (a *application) gracefulShutdownApplication(errCh <-chan error) error {
	var err error
	select {
	case <-a.ctx.Done():
		log.Println(msgs.ApplicationReceiveCtxDone)
	case err = <-errCh:
		a.cancelFunc()
		log.Println(msgs.ApplicationReceiveInternalError)
	}
	a.wg.Wait()
	return err
}

func (a *application) run(r runnable.Runnable, name string) {
	a.wg.Add(1)
	defer a.wg.Done()
	log.Println(msgs.RunRunnable, name)
	err := r.Run(a.ctx)
	if err != nil {
		a.errorf("%s run: %s", name, err)
	}
	log.Println(msgs.StoppedRunnable, name)
}

func (a *application) errorf(format string, args ...any) {
	select {
	case a.errCh <- fmt.Errorf(format, args...):
	default:
	}
}
