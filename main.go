package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/limits"
)

var (

	//Binance
	binanceAPIKey    = "<your binance api key>"
	binanceSecretKey = "<your binance secret key>"
)

func humanTimeToUnixNanoTime(input time.Time) int64 {
	return int64(time.Nanosecond) * input.UnixNano() / int64(time.Millisecond)
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(20)

	// Up some limits.
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}

	client := binance.NewClient(binanceAPIKey, binanceSecretKey)

	targetSymbol := "BNBUSDT"

	// reference https://golang.org/pkg/time/#Date
	// func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time

	// NOTE : Binance only accepts unix nano timestamp, not unix time stamp
	// For example
	// 1527616800  <----- will return empty result or wrong date range
	// 1527616899999 <--- ok

	//startTimeHumanReadable := time.Date(2018, time.October, 10, 0, 0, 0, 0, time.UTC)
	//start := humanTimeToUnixNanoTime(startTimeHumanReadable)

	//endTimeHumanReadable := time.Date(2018, time.October, 24, 0, 0, 0, 0, time.UTC)
	//end := humanTimeToUnixNanoTime(endTimeHumanReadable)

	targetInterval := "1d" /// case sensitive, 1D != 1d or 30M != 30m

	data, err := client.NewKlinesService().Symbol(targetSymbol).Interval(targetInterval).Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	totalGain := 0.0
	totalLoss := 0.0

	for i := 1; i < len(data); i++ {

		previous := data[i].Close
		current := data[i-1].Close

		// convert string to float64
		previousClose, _ := strconv.ParseFloat(previous, 64)
		currentClose, _ := strconv.ParseFloat(current, 64)

		difference := currentClose - previousClose

		if difference >= 0 {
			totalGain += difference
		} else {
			totalLoss -= difference
		}
	}

	rs := totalGain / math.Abs(totalLoss)
	rsi := 100 - 100/(1+rs)
	fmt.Println("RSI for "+targetSymbol+" : ", rsi)
}
