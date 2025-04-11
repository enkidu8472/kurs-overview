package config

// https://eodhd.com/
// i need it to download on which markets a ISIN is present
// with free token, i have 20 request per day (500 calls for new token)
var ApiTokenEodhd = "put your token here"

// to store the info of the EOD meta call
var EodSearchFileName = "eod_search_info.json"

// to store the section of eod info that was used last time to download time series
var EodLastUsedSymbolFileName = "eod_last_used.json"

// www.alphavantage.co/ function=TIME_SERIES_DAILY
// with this token daily stock data are downloaded
// > standard API rate limit is 25 requests per day.
var ApiTokenAlpha = "put your token here"

// a file name fragement to strore the information of the EOD meta call
var AlphaTimeSerieFileName = "alphaV_time_serie.csv"

// strore custom addtitional information on ISIN level
var CustomIsinInfoFileName = "customInfo.json"

var DataPath = "/home/rhytm/abr/0_proc/kurs-overview/data/"

var DataIsinPath = DataPath + "isin/"

// to store the transactions that we get from Trade Republic
var TrLogFile = "tr-transaction-log.csv"
