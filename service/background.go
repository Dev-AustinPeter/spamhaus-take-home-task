package service

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/Dev-AustinPeter/spamhaus-take-home-task/constants"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/types"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/utils"
)

func StartBackgroundFetch() {
	ticker := time.NewTicker(time.Duration(constants.FETCH_INTERVAL) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("[INFO] Running background fetch...")
		var urls []*types.URLData
		utils.URLStore.Range(func(_, value interface{}) bool {
			urls = append(urls, value.(*types.URLData))

			return true
		})

		sort.Slice(urls, func(i, j int) bool { return urls[i].Count > urls[j].Count })
		if len(urls) > 10 {
			urls = urls[:10]
		}

		var wg sync.WaitGroup
		for _, urlData := range urls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				utils.FetchURL(url)
			}(urlData.URL)
		}
		wg.Wait()
		log.Println("[INFO] Background fetch completed")
	}
}
