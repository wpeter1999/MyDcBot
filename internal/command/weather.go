package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

const cwaForecastEndpoint = "https://opendata.cwa.gov.tw/api/v1/rest/datastore/F-C0032-001"

// 台灣 22 個縣市（用於 choices）
var taiwanLocations = []string{
	"臺北市", "新北市", "桃園市", "臺中市", "臺南市", "高雄市",
	"基隆市", "新竹市", "嘉義市",
	"新竹縣", "苗栗縣", "彰化縣", "南投縣", "雲林縣", "嘉義縣", "屏東縣",
	"宜蘭縣", "花蓮縣", "臺東縣", "澎湖縣", "金門縣", "連江縣",
}

// WeatherCommand 定義 /weather 指令
var WeatherCommand = &BotCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "weather",
		Description: "查詢台灣各縣市今明 36 小時天氣預報",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "location",
				Description: "請選擇縣市",
				Required:    true,
				Choices:     buildLocationChoices(),
			},
		},
	},
	Handler: weatherCommandHandler,
}

// 建立 choices
func buildLocationChoices() []*discordgo.ApplicationCommandOptionChoice {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(taiwanLocations))
	for i, loc := range taiwanLocations {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  loc,
			Value: loc,
		}
	}
	return choices
}

type weatherForecast struct {
	LocationName    string
	Weather         string
	RainProbability string
	MinTemperature  string
	MaxTemperature  string
}

type weatherFetcher interface {
	FetchForecast(apiKey, location string) (weatherForecast, error)
}

type cwaWeatherFetcher struct {
	client  *http.Client
	baseURL string
}

var (
	errNoWeatherData                = errors.New("no weather data")
	weatherClient    weatherFetcher = &cwaWeatherFetcher{
		client:  http.DefaultClient,
		baseURL: cwaForecastEndpoint,
	}
	now = time.Now
)

type cwaResponse struct {
	Records struct {
		Location []struct {
			LocationName   string `json:"locationName"`
			WeatherElement []struct {
				ElementName string `json:"elementName"`
				Time        []struct {
					StartTime string `json:"startTime"`
					EndTime   string `json:"endTime"`
					Parameter struct {
						ParameterName string `json:"parameterName"`
					} `json:"parameter"`
				} `json:"time"`
			} `json:"weatherElement"`
		} `json:"location"`
	} `json:"records"`
}

func (f *cwaWeatherFetcher) FetchForecast(apiKey, location string) (weatherForecast, error) {
	apiURL, err := f.buildURL(apiKey, location)
	if err != nil {
		return weatherForecast{}, err
	}

	client := f.client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return weatherForecast{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weatherForecast{}, fmt.Errorf("cwa api returned status %d", resp.StatusCode)
	}

	var data cwaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weatherForecast{}, err
	}

	return parseWeatherForecast(data, location)
}

func (f *cwaWeatherFetcher) buildURL(apiKey, location string) (string, error) {
	baseURL := f.baseURL
	if baseURL == "" {
		baseURL = cwaForecastEndpoint
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("Authorization", apiKey)
	query.Set("format", "JSON")
	query.Set("locationName", location)
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func parseWeatherForecast(data cwaResponse, location string) (weatherForecast, error) {
	if len(data.Records.Location) == 0 {
		return weatherForecast{}, fmt.Errorf("%w: %s", errNoWeatherData, location)
	}

	loc := data.Records.Location[0]
	forecast := weatherForecast{LocationName: loc.LocationName}

	for _, elem := range loc.WeatherElement {
		if len(elem.Time) == 0 {
			continue
		}
		value := elem.Time[0].Parameter.ParameterName
		switch elem.ElementName {
		case "Wx":
			forecast.Weather = value
		case "PoP":
			forecast.RainProbability = value
		case "MinT":
			forecast.MinTemperature = value
		case "MaxT":
			forecast.MaxTemperature = value
		}
	}

	return forecast, nil
}

func formatWeatherForecast(forecast weatherForecast, currentTime time.Time) string {
	msg := fmt.Sprintf("**%s** 今明 36 小時天氣預報：\n", forecast.LocationName)

	if forecast.Weather != "" {
		msg += fmt.Sprintf("天氣現象：%s\n", forecast.Weather)
	}
	if forecast.RainProbability != "" {
		msg += fmt.Sprintf("降雨機率：%s%%\n", forecast.RainProbability)
	}
	if forecast.MinTemperature != "" {
		msg += fmt.Sprintf("最低氣溫：%s°C\n", forecast.MinTemperature)
	}
	if forecast.MaxTemperature != "" {
		msg += fmt.Sprintf("最高氣溫：%s°C\n", forecast.MaxTemperature)
	}

	msg += fmt.Sprintf("\n資料時間：%s", currentTime.Format("2006-01-02 15:04"))
	return msg
}

func weatherCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 取得使用者選擇的縣市（必填）
	var location string
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "location" {
			location = opt.StringValue()
		}
	}

	apiKey := os.Getenv("CWA_API_KEY")
	if apiKey == "" {
		respond(s, i, "尚未設定 CWA_API_KEY 環境變數，無法查詢天氣資料。")
		return
	}

	forecast, err := weatherClient.FetchForecast(apiKey, location)
	if errors.Is(err, errNoWeatherData) {
		respond(s, i, fmt.Sprintf("找不到「%s」的天氣資料。", location))
		return
	}
	if err != nil {
		log.Printf("查詢氣象資料失敗: %v", err)
		respond(s, i, "查詢天氣資料失敗，請稍後再試。")
		return
	}

	respond(s, i, formatWeatherForecast(forecast, now()))
}
