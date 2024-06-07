package screeners_macd

type MacdScreener struct {
}

func NewMacdScreener() MacdScreener {
	return MacdScreener{}
}

func (ms *MacdScreener) StartScanning() error {

	// Взять макс инструментов что сможем сфетчить

	// TODO: Сделать класс, лимитирующий подключения, в котором был бы метод типо Reserve() чтобы резервировать N лимитов под нужные вжи, например при запуске стратегии
	
	return nil
}