package models

type ServerSettings struct {
	Name                                 string           `json:"name"`
	Description                          string           `json:"description"`
	Tags                                 []string         `json:"tags"`
	MaxPlayers                           int              `json:"max_players"`
	Visibility                           ServerVisibility `json:"visibility"`
	Token                                string           `json:"token"`
	GamePassword                         string           `json:"game_password"`
	RequireUserVerification              bool             `json:"require_user_verification"`
	MaxUpload                            int              `json:"max_upload_in_kilobytes_per_second"`
	MinimumLatency                       int              `json:"minimum_latency_in_ticks"`
	IgnorePlayerLimitForReturningPlayers bool             `json:"ignore_player_limit_for_returning_players"`
	AllowCommands                        string           `json:"allow_commands"`
	AutosaveInterval                     int              `json:"autosave_interval"`
	AutosaveSlots                        int              `json:"autosave_slots"`
	AfkAutokickInterval                  int              `json:"afk_autokick_interval"`
	AutoPause                            bool             `json:"auto_pause"`
	OnlyAdminsCanPauseTheGame            bool             `json:"only_admins_can_pause_the_game"`
	AutosaveOnlyOnServer                 bool             `json:"autosave_only_on_server"`
	Admins                               []string         `json:"admins"`
}

type ServerVisibility struct {
	Public bool
	Lan    bool
}
