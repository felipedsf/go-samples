package service

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const (
	TIMEOUT = 10 * time.Millisecond
	INSERT  = "insert into exchange (code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) values (?,?,?,?,?,?,?,?,?,?,?)"
)

type Exchange struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type ExchangeResult struct {
	Bid string `json:"bid"`
}

type ExchangeService struct {
	Db *sql.DB
}

func (s ExchangeService) InsertExchange(exchange Exchange) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	stmt, err := s.Db.Prepare(INSERT)
	if err != nil {
		log.Fatal("error to prepare statement: ", err)
		return
	}

	_, err = stmt.ExecContext(ctx, exchange.Usdbrl.Code, exchange.Usdbrl.Codein, exchange.Usdbrl.Name, exchange.Usdbrl.High, exchange.Usdbrl.Low, exchange.Usdbrl.VarBid, exchange.Usdbrl.PctChange, exchange.Usdbrl.Bid, exchange.Usdbrl.Ask, exchange.Usdbrl.Timestamp, exchange.Usdbrl.CreateDate)
	if err != nil {
		log.Fatal("error to execute insert: ", err)
		return
	}
}
