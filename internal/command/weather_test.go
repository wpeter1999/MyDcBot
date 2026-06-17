package command

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TestWeatherCommand_Definition 測試 WeatherCommand 的基本定義
func TestWeatherCommand_Definition(t *testing.T) {
	if WeatherCommand.Command.Name != "weather" {
		t.Errorf("WeatherCommand 名稱應為 'weather'，實際為 '%s'", WeatherCommand.Command.Name)
	}
	if WeatherCommand.Command.Description == "" {
		t.Error("WeatherCommand 應該有描述")
	}
	if WeatherCommand.Handler == nil {
		t.Error("WeatherCommand 必須有 Handler")
	}
}

// TestWeatherCommand_HasRequiredLocation 測試 WeatherCommand 是否有必填的 location 參數
func TestWeatherCommand_HasRequiredLocation(t *testing.T) {
	if len(WeatherCommand.Command.Options) == 0 {
		t.Fatal("WeatherCommand 應該有至少一個參數")
	}

	opt := WeatherCommand.Command.Options[0]
	if opt.Name != "location" {
		t.Errorf("第一個參數應為 'location'，實際為 '%s'", opt.Name)
	}
	if !opt.Required {
		t.Error("location 參數應該是必填的")
	}
	if len(opt.Choices) != 22 {
		t.Errorf("location 應該有 22 個 choices，實際有 %d 個", len(opt.Choices))
	}
}

// TestWeatherCommand_IsRegistered 測試 WeatherCommand 是否已註冊
func TestWeatherCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range CommandRegistry {
		if cmd.Command.Name == "weather" {
			found = true
			break
		}
	}
	if !found {
		t.Error("WeatherCommand 未註冊到 CommandRegistry")
	}
}

// TestBuildLocationChoices 測試台灣縣市 choices 的數量與 name/value 對應
func TestBuildLocationChoices(t *testing.T) {
	choices := buildLocationChoices()
	if len(choices) != len(taiwanLocations) {
		t.Fatalf("預期 choices 數量為 %d，實際為 %d", len(taiwanLocations), len(choices))
	}

	for i, loc := range taiwanLocations {
		if choices[i].Name != loc {
			t.Errorf("第 %d 個 choice name 應為 %q，實際為 %q", i, loc, choices[i].Name)
		}
		if choices[i].Value != loc {
			t.Errorf("第 %d 個 choice value 應為 %q，實際為 %v", i, loc, choices[i].Value)
		}
	}
}

// TestCWAWeatherFetcherBuildURL 測試 CWA API URL 會包含必要 query 並正確 encode 縣市名稱
func TestCWAWeatherFetcherBuildURL(t *testing.T) {
	fetcher := &cwaWeatherFetcher{baseURL: "https://example.test/weather"}

	got, err := fetcher.buildURL("test-key", "臺北市")
	if err != nil {
		t.Fatalf("buildURL 不應回傳錯誤，但得到: %v", err)
	}

	if !strings.HasPrefix(got, "https://example.test/weather?") {
		t.Fatalf("buildURL 應保留 base URL，實際為 %q", got)
	}
	if !strings.Contains(got, "Authorization=test-key") {
		t.Errorf("buildURL 應包含 Authorization query，實際為 %q", got)
	}
	if !strings.Contains(got, "format=JSON") {
		t.Errorf("buildURL 應包含 format query，實際為 %q", got)
	}
	if !strings.Contains(got, "locationName=%E8%87%BA%E5%8C%97%E5%B8%82") {
		t.Errorf("buildURL 應 encode locationName query，實際為 %q", got)
	}
}

// TestCWAWeatherFetcherBuildURL_UsesDefaultEndpoint 測試未指定 baseURL 時會使用預設 CWA endpoint
func TestCWAWeatherFetcherBuildURL_UsesDefaultEndpoint(t *testing.T) {
	fetcher := &cwaWeatherFetcher{}

	got, err := fetcher.buildURL("test-key", "臺北市")
	if err != nil {
		t.Fatalf("buildURL 不應回傳錯誤，但得到: %v", err)
	}

	if !strings.HasPrefix(got, cwaForecastEndpoint+"?") {
		t.Fatalf("baseURL 空白時應使用預設 endpoint，實際為 %q", got)
	}
}

// TestCWAWeatherFetcherBuildURL_ReturnsErrorForInvalidBaseURL 測試 baseURL 無效時會回傳錯誤
func TestCWAWeatherFetcherBuildURL_ReturnsErrorForInvalidBaseURL(t *testing.T) {
	fetcher := &cwaWeatherFetcher{baseURL: "://bad-url"}

	_, err := fetcher.buildURL("test-key", "臺北市")
	if err == nil {
		t.Fatal("baseURL 無效時應回傳錯誤")
	}
}

// TestCWAWeatherFetcherFetchForecast_Success 測試成功呼叫 CWA API 後會解析天氣預報資料
func TestCWAWeatherFetcherFetchForecast_Success(t *testing.T) {
	var gotAuthorization string
	var gotFormat string
	var gotLocation string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthorization = r.URL.Query().Get("Authorization")
		gotFormat = r.URL.Query().Get("format")
		gotLocation = r.URL.Query().Get("locationName")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"records": {
				"location": [
					{
						"locationName": "臺北市",
						"weatherElement": [
							{"elementName": "Wx", "time": [{"parameter": {"parameterName": "多雲"}}]},
							{"elementName": "PoP", "time": [{"parameter": {"parameterName": "30"}}]},
							{"elementName": "MinT", "time": [{"parameter": {"parameterName": "22"}}]},
							{"elementName": "MaxT", "time": [{"parameter": {"parameterName": "28"}}]}
						]
					}
				]
			}
		}`))
	}))
	defer server.Close()

	fetcher := &cwaWeatherFetcher{client: server.Client(), baseURL: server.URL}
	forecast, err := fetcher.FetchForecast("test-key", "臺北市")
	if err != nil {
		t.Fatalf("預期不應發生錯誤，但得到: %v", err)
	}

	if gotAuthorization != "test-key" {
		t.Errorf("Authorization query 應為 test-key，實際為 %q", gotAuthorization)
	}
	if gotFormat != "JSON" {
		t.Errorf("format query 應為 JSON，實際為 %q", gotFormat)
	}
	if gotLocation != "臺北市" {
		t.Errorf("locationName query 應為 臺北市，實際為 %q", gotLocation)
	}

	if forecast.LocationName != "臺北市" {
		t.Errorf("LocationName 應為 臺北市，實際為 %q", forecast.LocationName)
	}
	if forecast.Weather != "多雲" {
		t.Errorf("Weather 應為 多雲，實際為 %q", forecast.Weather)
	}
	if forecast.RainProbability != "30" {
		t.Errorf("RainProbability 應為 30，實際為 %q", forecast.RainProbability)
	}
	if forecast.MinTemperature != "22" {
		t.Errorf("MinTemperature 應為 22，實際為 %q", forecast.MinTemperature)
	}
	if forecast.MaxTemperature != "28" {
		t.Errorf("MaxTemperature 應為 28，實際為 %q", forecast.MaxTemperature)
	}
}

// TestCWAWeatherFetcherFetchForecast_ReturnsErrorForNonOK 測試 CWA API 回傳非 200 狀態時會回傳錯誤
func TestCWAWeatherFetcherFetchForecast_ReturnsErrorForNonOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	fetcher := &cwaWeatherFetcher{client: server.Client(), baseURL: server.URL}
	_, err := fetcher.FetchForecast("test-key", "臺北市")
	if err == nil {
		t.Fatal("預期非 200 狀態應回傳錯誤")
	}
}

// TestCWAWeatherFetcherFetchForecast_ReturnsErrorForMalformedJSON 測試 CWA API 回傳無效 JSON 時會回傳錯誤
func TestCWAWeatherFetcherFetchForecast_ReturnsErrorForMalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{`))
	}))
	defer server.Close()

	fetcher := &cwaWeatherFetcher{client: server.Client(), baseURL: server.URL}
	_, err := fetcher.FetchForecast("test-key", "臺北市")
	if err == nil {
		t.Fatal("預期 JSON 格式錯誤時應回傳錯誤")
	}
}

// TestParseWeatherForecast_ReturnsForecastWithAvailableElements 測試 parseWeatherForecast 會忽略沒有 time 的天氣元素
func TestParseWeatherForecast_ReturnsForecastWithAvailableElements(t *testing.T) {
	var data cwaResponse
	data.Records.Location = append(data.Records.Location, struct {
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
	}{LocationName: "臺北市"})
	data.Records.Location[0].WeatherElement = append(data.Records.Location[0].WeatherElement, struct {
		ElementName string `json:"elementName"`
		Time        []struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
			Parameter struct {
				ParameterName string `json:"parameterName"`
			} `json:"parameter"`
		} `json:"time"`
	}{ElementName: "Wx"})

	forecast, err := parseWeatherForecast(data, "臺北市")
	if err != nil {
		t.Fatalf("parseWeatherForecast 不應回傳錯誤，但得到: %v", err)
	}
	if forecast.LocationName != "臺北市" {
		t.Fatalf("LocationName 應為 臺北市，實際為 %q", forecast.LocationName)
	}
	if forecast.Weather != "" {
		t.Fatalf("空的 time slice 應被忽略，Weather 實際為 %q", forecast.Weather)
	}
}

// TestParseWeatherForecast_ReturnsNoDataError 測試沒有天氣資料時會回傳 errNoWeatherData
func TestParseWeatherForecast_ReturnsNoDataError(t *testing.T) {
	_, err := parseWeatherForecast(cwaResponse{}, "不存在")
	if !errors.Is(err, errNoWeatherData) {
		t.Fatalf("預期 errNoWeatherData，實際為 %v", err)
	}
}

// TestFormatWeatherForecast 測試完整天氣預報訊息格式
func TestFormatWeatherForecast(t *testing.T) {
	forecast := weatherForecast{
		LocationName:    "臺北市",
		Weather:         "多雲",
		RainProbability: "30",
		MinTemperature:  "22",
		MaxTemperature:  "28",
	}
	currentTime := time.Date(2026, 6, 17, 9, 30, 0, 0, time.Local)

	got := formatWeatherForecast(forecast, currentTime)
	want := "**臺北市** 今明 36 小時天氣預報：\n" +
		"天氣現象：多雲\n" +
		"降雨機率：30%\n" +
		"最低氣溫：22°C\n" +
		"最高氣溫：28°C\n" +
		"\n資料時間：2026-06-17 09:30"

	if got != want {
		t.Fatalf("格式化結果不符合預期\nwant: %q\n got: %q", want, got)
	}
}

// TestFormatWeatherForecast_SkipsEmptyFields 測試空白欄位不會出現在天氣預報訊息中
func TestFormatWeatherForecast_SkipsEmptyFields(t *testing.T) {
	forecast := weatherForecast{LocationName: "臺北市"}
	currentTime := time.Date(2026, 6, 17, 9, 30, 0, 0, time.Local)

	got := formatWeatherForecast(forecast, currentTime)
	want := "**臺北市** 今明 36 小時天氣預報：\n\n資料時間：2026-06-17 09:30"

	if got != want {
		t.Fatalf("格式化結果不符合預期\nwant: %q\n got: %q", want, got)
	}
}

// TestWeatherCommandHandler_MissingAPIKey 測試未設定 CWA_API_KEY 時會回覆設定錯誤訊息
func TestWeatherCommandHandler_MissingAPIKey(t *testing.T) {
	t.Setenv("CWA_API_KEY", "")
	got := captureWeatherResponse(t, &fakeWeatherFetcher{}, time.Now, weatherInteraction("臺北市"))

	want := "尚未設定 CWA_API_KEY 環境變數，無法查詢天氣資料。"
	if got != want {
		t.Fatalf("回應不符合預期\nwant: %q\n got: %q", want, got)
	}
}

// TestWeatherCommandHandler_Success 測試 weather handler 成功取得預報時會回覆格式化天氣訊息
func TestWeatherCommandHandler_Success(t *testing.T) {
	t.Setenv("CWA_API_KEY", "test-key")
	fetcher := &fakeWeatherFetcher{
		forecast: weatherForecast{
			LocationName:    "臺北市",
			Weather:         "晴時多雲",
			RainProbability: "20",
			MinTemperature:  "24",
			MaxTemperature:  "31",
		},
	}
	fixedNow := func() time.Time {
		return time.Date(2026, 6, 17, 10, 45, 0, 0, time.Local)
	}

	got := captureWeatherResponse(t, fetcher, fixedNow, weatherInteraction("臺北市"))

	if fetcher.apiKey != "test-key" {
		t.Errorf("fetcher apiKey 應為 test-key，實際為 %q", fetcher.apiKey)
	}
	if fetcher.location != "臺北市" {
		t.Errorf("fetcher location 應為 臺北市，實際為 %q", fetcher.location)
	}
	if !strings.Contains(got, "**臺北市** 今明 36 小時天氣預報：") {
		t.Errorf("回應應包含標題，實際為 %q", got)
	}
	if !strings.Contains(got, "資料時間：2026-06-17 10:45") {
		t.Errorf("回應應包含固定資料時間，實際為 %q", got)
	}
}

// TestWeatherCommandHandler_FetchError 測試取得天氣資料失敗時會回覆失敗訊息
func TestWeatherCommandHandler_FetchError(t *testing.T) {
	t.Setenv("CWA_API_KEY", "test-key")
	fetcher := &fakeWeatherFetcher{err: errors.New("api failed")}

	got := captureWeatherResponse(t, fetcher, time.Now, weatherInteraction("臺北市"))

	want := "查詢天氣資料失敗，請稍後再試。"
	if got != want {
		t.Fatalf("回應不符合預期\nwant: %q\n got: %q", want, got)
	}
}

// TestWeatherCommandHandler_NoWeatherData 測試查無天氣資料時會回覆找不到資料訊息
func TestWeatherCommandHandler_NoWeatherData(t *testing.T) {
	t.Setenv("CWA_API_KEY", "test-key")
	fetcher := &fakeWeatherFetcher{err: errNoWeatherData}

	got := captureWeatherResponse(t, fetcher, time.Now, weatherInteraction("臺北市"))

	want := "找不到「臺北市」的天氣資料。"
	if got != want {
		t.Fatalf("回應不符合預期\nwant: %q\n got: %q", want, got)
	}
}

type fakeWeatherFetcher struct {
	forecast weatherForecast
	err      error
	apiKey   string
	location string
}

func (f *fakeWeatherFetcher) FetchForecast(apiKey, location string) (weatherForecast, error) {
	f.apiKey = apiKey
	f.location = location
	return f.forecast, f.err
}

func captureWeatherResponse(t *testing.T, fetcher weatherFetcher, currentTime func() time.Time, interaction *discordgo.InteractionCreate) string {
	t.Helper()

	originalFetcher := weatherClient
	originalNow := now
	originalResponder := respondToInteraction
	t.Cleanup(func() {
		weatherClient = originalFetcher
		now = originalNow
		respondToInteraction = originalResponder
	})

	weatherClient = fetcher
	now = currentTime

	var got string
	respondToInteraction = func(_ *discordgo.Session, _ *discordgo.InteractionCreate, content string) {
		got = content
	}

	weatherCommandHandler(&discordgo.Session{}, interaction)
	return got
}

func weatherInteraction(location string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "weather",
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "location",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: location,
					},
				},
			},
		},
	}
}

// TestMain 執行 command package 測試進入點
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
