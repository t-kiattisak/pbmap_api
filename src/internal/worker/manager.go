package worker

import (
	"pbmap_api/src/config"
	"pbmap_api/src/internal/adapter/pokeapi"
	"pbmap_api/src/internal/scheduler"
	"pbmap_api/src/internal/usecase"
)

func StartBackgroundJobs(cfg *config.Config) func() {
	fcmService, _ := usecase.NewFCMService(cfg)
	pokeClient := pokeapi.NewPokeAPIClient()
	dataSyncService := usecase.NewDataSyncService(fcmService, pokeClient)

	appScheduler := scheduler.NewScheduler(dataSyncService)
	appScheduler.Start()

	return func() {
		appScheduler.Stop()
	}
}
