package structs

type User struct {
	ID                            int
	RemoteUserID                  int
	SecurityToken                 string
	Email                         string
	Password                      string
	DateCreated                   string
	Banned                        bool
	Premium                       bool
	Admin                         bool
	GameLauncherCertificate       string
	GameLauncherHash              string
	HiddenHWID                    string
	HWID                          string
	OSVersion                     string
	UserAgent                     string
	IPAddress                     string
	Personas                      []Persona
	DefaultPersonaIDX             int
	ActivePersonaID               int
	Locked                        bool
	SelectedPersonaIndex          int
	FullGameAccess                bool
	Complete                      bool
	LastAuthDate                  string
	SubscribeMsg                  bool
	Address1                      string
	Address2                      string
	Country                       string
	DOB                           string
	EmailStatus                   string
	FirstName                     string
	Gender                        string
	IDDigits                      string
	LandlinePhone                 string
	Language                      string
	LastName                      string
	Mobile                        string
	Nickname                      string
	PostalCode                    string
	RealName                      string
	ReasonCode                    string
	StarterPackEntitlementTag     string
	Status                        string
	TOSVersion                    string
	Username                      string
	AppearOffline                 bool
	DeclineGroupInvitations       int
	DeclineIncomingFriendRequests bool
	DeclinePrivateInvite          int
	HideOfflineFriends            bool
	ShowNewsOnSignIn              bool
}

type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Ticket   string `json:"ticket"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
