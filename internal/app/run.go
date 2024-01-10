package app

import (
	"context"
	"fmt"
	"log"
	grpcside "merger-adapter/internal/api/grpc_merger_client"
	"merger-adapter/internal/api/telegrambot"
	"merger-adapter/internal/common/msgs"
	"merger-adapter/internal/config"
	"merger-adapter/internal/service/runnable"
)

func Run(ctx context.Context, cfg *config.Config) error {

	app, errCh := newApplication(ctx)

	server, err := grpcside.InitGrpcMergerClient(grpcside.Config{
		Host: cfg.ServerHost,
	})
	if err != nil {
		return fmt.Errorf("init grpc merger server: %s", err)
	}
	log.Println(msgs.InitGrpcMergerClientInitialized)

	tgbClient, err := telegrambot.InitClient(telegrambot.Config{
		Token:  cfg.TelegramBotToken,
		ChatID: cfg.TelegramChat,
		Server: server,
		ApiKey: cfg.XApiKey,
	})
	if err != nil {
		return fmt.Errorf("tg client initialization: %s", err)
	}
	log.Println(msgs.TelegramAdapterInitialized)

	go app.run(tgbClient, "telegram adapter")

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