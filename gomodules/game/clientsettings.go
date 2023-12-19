package game

// Represents the settings for an EchoVR client.
type EchoClientSettings struct {
	// WARNING: EchoVR dictates this schema.
	ConfigData          map[string]interface{} `json:"config_data"`           // ConfigData is a map that stores configuration data for the EchoVR client.
	Env                 string                 `json:"env"`                   // Env represents the environment in which the EchoVR client is running.
	IapUnlocked         bool                   `json:"iap_unlocked"`          // IapUnlocked indicates whether in-app purchases are unlocked for the EchoVR client.
	MatchmakerQueueMode string                 `json:"matchmaker_queue_mode"` // MatchmakerQueueMode specifies the queue mode for the EchoVR client's matchmaker.

	RemoteLogErrors       bool `json:"remote_log_errors"`        // send remote logs for errors
	RemoteLogMetrics      bool `json:"remote_log_metrics"`       // send remote logs for metrics
	RemoteLogRichPresence bool `json:"remote_log_rich_presence"` // send remote logs for rich presence
	RemoteLogSocial       bool `json:"remote_log_social"`        // send remote logs for social events
	RemoteLogWarnings     bool `json:"remote_log_warnings"`      // send remote logs for warnings
}

func DefaultEchoClientSettings() EchoClientSettings {
	return EchoClientSettings{
		ConfigData:            make(map[string]interface{}),
		Env:                   "live",
		IapUnlocked:           false,
		MatchmakerQueueMode:   "disabled",
		RemoteLogErrors:       true,
		RemoteLogMetrics:      false,
		RemoteLogRichPresence: false,
		RemoteLogSocial:       false,
		RemoteLogWarnings:     false,
	}
}
