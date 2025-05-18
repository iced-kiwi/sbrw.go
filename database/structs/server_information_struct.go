package structs

type ServerInformation struct {
	MaxPlayers                       int      `json:"maxPlayers"`
	OnlinePlayers                    int      `json:"onlinePlayers"`
	ModernAuthSupport                bool     `json:"modernAuthSupport"`
	MessageSrv                       string   `json:"messageSrv"`
	HomePageURL                      string   `json:"homePageUrl"`
	DiscordURL                       string   `json:"discordUrl"`
	FacebookURL                      string   `json:"facebookUrl"`
	TwitterURL                       string   `json:"twitterUrl"`
	ServerName                       string   `json:"serverName"`
	Country                          string   `json:"country"`
	Timezone                         int      `json:"timezone"`
	BannerURL                        string   `json:"bannerUrl"`
	AdminList                        []string `json:"adminList,omitempty"`
	OwnerList                        []string `json:"ownerList,omitempty"`
	NumberOfRegistered               int      `json:"numberOfRegistered"`
	SecondsToShutDown                int      `json:"secondsToShutDown"`
	ActivatedHolidaySceneryGroups    []string `json:"activatedHolidaySceneryGroups"`
	DisactivatedHolidaySceneryGroups []string `json:"disactivatedHolidaySceneryGroups"`
	RequireTicket                    bool     `json:"requireTicket"`
	ServerVersion                    string   `json:"serverVersion"`
}

func NewDefaultServerInformation() ServerInformation {
	return ServerInformation{
		MaxPlayers:                       1,
		OnlinePlayers:                    0,
		ModernAuthSupport:                true,
		MessageSrv:                       "https://discord.com/app",
		HomePageURL:                      "https://github.com/",
		DiscordURL:                       "https://discord.com/app",
		FacebookURL:                      "https://facebook.com/",
		TwitterURL:                       "https://twitter.com/",
		ServerName:                       "SBRW go",
		Country:                          "GB",
		Timezone:                         0,
		BannerURL:                        "",
		NumberOfRegistered:               0,
		SecondsToShutDown:                14400,
		ActivatedHolidaySceneryGroups:    []string{},
		DisactivatedHolidaySceneryGroups: []string{},
		RequireTicket:                    true,
		ServerVersion:                    "0.0.1-A",
	}
}
