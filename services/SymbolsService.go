package services

import (
	"fmt"
	"sync"
	"uacc-backend/integrations"
)

type SymbolsResponse map[string]string

type SymbolsService interface {
	GetSymbols() (SymbolsResponse, error)
}

type SymbolsServiceImpl struct {
	agents []integrations.ProxyAgent
}

func NewSymbolsService(agents []integrations.ProxyAgent) *SymbolsServiceImpl {
	return &SymbolsServiceImpl{agents}
}

func (service *SymbolsServiceImpl) GetSymbols() (SymbolsResponse, error) {
	var wg sync.WaitGroup
	result := make(SymbolsResponse)
	resultsChannel := make(chan integrations.SymbolsResponse, len(service.agents))

	for _, agent := range service.agents {
		wg.Add(1)
		go func(agent integrations.ProxyAgent) {
			defer wg.Done()
			symbols, err := agent.GetSymbols()
			if err != nil {
				fmt.Println("Error getting symbols:", err)
				return
			}
			resultsChannel <- symbols
		}(agent)
	}

	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	for symbols := range resultsChannel {
		for key, value := range symbols {
			result[key] = value
		}
	}

	return result, nil
}
