package controllers

import (
	accounts "main/server/contracts/contracts.accounts"
	auth "main/server/contracts/contracts.auth"
	marketdata "main/server/contracts/contracts.marketdata"
	ping "main/server/contracts/contracts.ping"
	shares "main/server/contracts/contracts.shares"
	trade "main/server/contracts/contracts.trade"
)

type Server struct {
	accounts.UnimplementedAccountsServer
	auth.UnimplementedAuthServer
	marketdata.UnimplementedMarketDataServer
	shares.UnimplementedSharesServer
	trade.UnimplementedTradeServer
	ping.UnimplementedPingServer
}