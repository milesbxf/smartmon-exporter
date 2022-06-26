package smartctl

type SmartCtlInfo struct {
	SmartCtlVersion []int    `json:"version"`
	SvnRevision     string   `json:"svn_revision"`
	PlatformInfo    string   `json:"platform_info"`
	BuildInfo       string   `json:"build_info"`
	Argv            []string `json:"argv"`
	ExitStatus      int      `json:"exit_status"`
}

type Device struct {
	Name     string `json:"name"`
	InfoName string `json:"info_name"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
}

type ScanOpenOutput struct {
	SmartExitCodeOutput
	SmartCtlInfo `json:"smartctl"`
	Devices      []Device `json:"devices"`
}

type Wwn struct {
	Naa int   `json:"naa"`
	Oui int   `json:"oui"`
	Id  int64 `json:"id"`
}

type UserCapacity struct {
	Blocks int64 `json:"blocks"`
	Bytes  int64 `json:"bytes"`
}

type Trim struct {
	Supported bool `json:"supported"`
}

type AtaVersion struct {
	String     string `json:"string"`
	MajorValue int    `json:"major_value"`
	MinorValue int    `json:"minor_value"`
}

type SataVersion struct {
	String string `json:"string"`
	Value  int    `json:"value"`
}

type InterfaceSpeedMax struct {
	SataValue      int    `json:"sata_value"`
	String         string `json:"string"`
	UnitsPerSecond int    `json:"units_per_second"`
	BitsPerUnit    int    `json:"bits_per_unit"`
}

type InterfaceSpeed struct {
	Max InterfaceSpeedMax `json:"max"`
}

type LocalTime struct {
	TimeT   int    `json:"time_t"`
	Asctime string `json:"asctime"`
}

type SmartStatus struct {
	Passed bool `json:"passed"`
}

type OfflineDataCollection struct {
	Status struct {
		Value  int    `json:"value"`
		String string `json:"string"`
		Passed bool   `json:"passed"`
	} `json:"status"`
	CompletionSeconds int `json:"completion_seconds"`
}

type SelfTestStatus struct {
	Value  int    `json:"value"`
	String string `json:"string"`
	Passed bool   `json:"passed"`
}

type SelfTestPollingMinutes struct {
	Short      int `json:"short"`
	Extended   int `json:"extended"`
	Conveyance int `json:"conveyance"`
}

type SelfTest struct {
	Status         SelfTestStatus         `json:"status"`
	PollingMinutes SelfTestPollingMinutes `json:"polling_minutes"`
}

type Capabilities struct {
	Values                        []int `json:"values"`
	ExecOfflineImmediateSupported bool  `json:"exec_offline_immediate_supported"`
	OfflineIsAbortedUponNewCmd    bool  `json:"offline_is_aborted_upon_new_cmd"`
	OfflineSurfaceScanSupported   bool  `json:"offline_surface_scan_supported"`
	SelfTestsSupported            bool  `json:"self_tests_supported"`
	ConveyanceSelfTestSupported   bool  `json:"conveyance_self_test_supported"`
	SelectiveSelfTestSupported    bool  `json:"selective_self_test_supported"`
	AttributeAutosaveEnabled      bool  `json:"attribute_autosave_enabled"`
	ErrorLoggingSupported         bool  `json:"error_logging_supported"`
	GpLoggingSupported            bool  `json:"gp_logging_supported"`
}

type AtaSmartData struct {
	OfflineDataCollection `json:"offline_data_collection"`
	SelfTest              `json:"self_test"`
	Capabilities          `json:"capabilities"`
}

type AtaSctCapabilities struct {
	Value                         int  `json:"value"`
	ErrorRecoveryControlSupported bool `json:"error_recovery_control_supported"`
	FeatureControlSupported       bool `json:"feature_control_supported"`
	DataTableSupported            bool `json:"data_table_supported"`
}

type AtaSmartAttributesFlags struct {
	Value         int    `json:"value"`
	String        string `json:"string"`
	Prefailure    bool   `json:"prefailure"`
	UpdatedOnline bool   `json:"updated_online"`
	Performance   bool   `json:"performance"`
	ErrorRate     bool   `json:"error_rate"`
	EventCount    bool   `json:"event_count"`
	AutoKeep      bool   `json:"auto_keep"`
}

type AtaSmartAttributesRaw struct {
	Value  int    `json:"value"`
	String string `json:"string"`
}

type AtaSmartAttributesTable struct {
	Id         int                     `json:"id"`
	Name       string                  `json:"name"`
	Value      int                     `json:"value"`
	Worst      int                     `json:"worst"`
	Thresh     int                     `json:"thresh"`
	WhenFailed string                  `json:"when_failed"`
	Flags      AtaSmartAttributesFlags `json:"flags"`
	Raw        AtaSmartAttributesRaw   `json:"raw"`
}

type AtaSmartAttributes struct {
	Revision int                       `json:"revision"`
	Table    []AtaSmartAttributesTable `json:"table"`
}

type PowerOnTime struct {
	Hours int `json:"hours"`
}

type Temperature struct {
	Current int `json:"current"`
}

type AtaSmartErrorLogSummary struct {
	Revision int `json:"revision"`
	Count    int `json:"count"`
}

type AtaSmartErrorLog struct {
	Summary AtaSmartErrorLogSummary `json:"summary"`
}

type AtaSmartSelfTestLogStandard struct {
	Revision int `json:"revision"`
	Count    int `json:"count"`
}

type AtaSmartSelfTestLog struct {
	Standard AtaSmartSelfTestLogStandard `json:"standard"`
}

type AtaSmartSelectiveSelfTestLogTable struct {
	LbaMin int `json:"lba_min"`
	LbaMax int `json:"lba_max"`
	Status struct {
		Value  int    `json:"value"`
		String string `json:"string"`
	} `json:"status"`
}

type AtaSmartSelectiveSelfTestFlags struct {
	Value                int  `json:"value"`
	RemainderScanEnabled bool `json:"remainder_scan_enabled"`
}

type AtaSmartSelectiveSelfTestLog struct {
	Revision                 int                                 `json:"revision"`
	Table                    []AtaSmartSelectiveSelfTestLogTable `json:"table"`
	Flags                    AtaSmartSelectiveSelfTestFlags      `json:"flags"`
	PowerUpScanResumeMinutes int                                 `json:"power_up_scan_resume_minutes"`
}

type InfoAllOutput struct {
	SmartExitCodeOutput
	SmartCtlInfo `json:"smartctl"`

	Device                       `json:"device"`
	ModelFamily                  string `json:"model_family"`
	ModelName                    string `json:"model_name"`
	SerialNumber                 string `json:"serial_number"`
	Wwn                          `json:"wwn"`
	FirmwareVersion              string `json:"firmware_version"`
	UserCapacity                 `json:"user_capacity"`
	LogicalBlockSize             int64 `json:"logical_block_size"`
	PhysicalBlockSize            int64 `json:"physical_block_size"`
	Trim                         `json:"trim"`
	InSmartctlDatabase           bool `json:"in_smartctl_database"`
	AtaVersion                   `json:"ata_version"`
	SataVersion                  `json:"sata_version"`
	InterfaceSpeed               `json:"interface_speed"`
	LocalTime                    `json:"local_time"`
	SmartStatus                  `json:"smart_status"`
	AtaSmartData                 `json:"ata_smart_data"`
	AtaSctCapabilities           `json:"ata_sct_capabilities"`
	AtaSmartAttributes           `json:"ata_smart_attributes"`
	PowerOnTime                  `json:"power_on_time"`
	PowerCycleCount              int `json:"power_cycle_count"`
	Temperature                  `json:"temperature"`
	AtaSmartErrorLog             `json:"ata_smart_error_log"`
	AtaSmartSelfTestLog          `json:"ata_smart_self_test_log"`
	AtaSmartSelectiveSelfTestLog `json:"ata_smart_selective_self_test_log"`
}
