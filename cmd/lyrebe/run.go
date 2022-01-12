package lyrebe

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
	"github.com/teal-seagull/lyre-be-v4/pkg/log"
)

const (
	envPrefix = "LYREBE"
)

var (
	globalCtx context.Context
	globalWG  *sync.WaitGroup
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Lyre-be API Server",
	Long:  `runs Lyre-be API Server`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		log.InitLogger(config.TheConfig().Logger)
		log.TheLogger().Debug("lyrebe component",
			zap.String("config", fmt.Sprintf("%#v", config.TheConfig())))

		runProcesses()
		globalWG.Wait()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	viper.SetEnvPrefix(envPrefix)

	var cancelFn func()
	globalCtx, cancelFn = context.WithCancel(context.Background())
	globalWG = &sync.WaitGroup{}

	globalWG.Add(1)
	go signalHandler(cancelFn)
}

func signalHandler(fn func()) {
	defer globalWG.Done()
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.TheLogger().Info("stop execution", zap.String("signal", sig.String()))
	fn()
	close(sigs)
}

func runProcesses() {
	var (
		apiSrv *apiserver.APIHTTPServer
		err    error
	)

	globalWG.Add(1)
	defer globalWG.Done()

	if apiSrv, err = apiserver.NewAPIHTTPServer(config.TheConfig().Server); err != nil {
		log.TheLogger().Fatal("error initializing API server", zap.Error(err))
	}

	apiSrv.Run(globalCtx)
}
