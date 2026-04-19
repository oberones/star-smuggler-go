package services

type TradeResult struct {
	Succeeded bool
	Message   string
}

func SuccessfulTrade(message string) TradeResult {
	return TradeResult{
		Succeeded: true,
		Message:   message,
	}
}

func FailedTrade(message string) TradeResult {
	return TradeResult{
		Succeeded: false,
		Message:   message,
	}
}
