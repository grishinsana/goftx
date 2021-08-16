package models

import "github.com/shopspring/decimal"

type AccountInformation struct {
	BackstopProvider             bool            `json:"backstopProvider"`
	Collateral                   decimal.Decimal `json:"collateral"`
	FreeCollateral               decimal.Decimal `json:"freeCollateral"`
	InitialMarginRequirement     decimal.Decimal `json:"initialMarginRequirement"`
	Liquidating                  bool            `json:"liquidating"`
	MaintenanceMarginRequirement decimal.Decimal `json:"maintenanceMarginRequirement"`
	MakerFee                     decimal.Decimal `json:"makerFee"`
	MarginFraction               decimal.Decimal `json:"marginFraction"`
	OpenMarginFraction           decimal.Decimal `json:"openMarginFraction"`
	TakerFee                     decimal.Decimal `json:"takerFee"`
	TotalAccountValue            decimal.Decimal `json:"totalAccountValue"`
	TotalPositionSize            decimal.Decimal `json:"totalPositionSize"`
	Username                     string          `json:"username"`
	Leverage                     decimal.Decimal `json:"leverage"`
	Positions                    []Position      `json:"positions"`
}

type Position struct {
	Cost                         decimal.Decimal `json:"cost"`
	EntryPrice                   decimal.Decimal `json:"entryPrice"`
	EstimatedLiquidationPrice    decimal.Decimal `json:"estimatedLiquidationPrice"`
	Future                       string          `json:"future"`
	InitialMarginRequirement     decimal.Decimal `json:"initialMarginRequirement"`
	LongOrderSize                decimal.Decimal `json:"longOrderSize"`
	MaintenanceMarginRequirement decimal.Decimal `json:"maintenanceMarginRequirement"`
	NetSize                      decimal.Decimal `json:"netSize"`
	OpenSize                     decimal.Decimal `json:"openSize"`
	RealizedPnl                  decimal.Decimal `json:"realizedPnl"`
	ShortOrderSize               decimal.Decimal `json:"shortOrderSize"`
	Side                         string          `json:"side"`
	Size                         decimal.Decimal `json:"size"`
	UnrealizedPnl                decimal.Decimal `json:"unrealizedPnl"`
	CollateralUsed               decimal.Decimal `json:"collateralUsed"`
	RecentAverageOpenPrice       decimal.Decimal `json:"recentAverageOpenPrice"`
}
