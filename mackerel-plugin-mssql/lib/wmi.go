// +build windows

package mpmssql

// Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSBufferManager is struct for WMI
type Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSBufferManager struct {
	BackgroundwriterpagesPersec   uint16
	Buffercachehitratio           uint16
	Buffercachehitratio_Base      uint16
	CheckpointpagesPersec         uint16
	Databasepages                 uint16
	Extensionallocatedpages       uint16
	Extensionfreepages            uint16
	Extensioninuseaspercentage    uint16
	ExtensionoutstandingIOcounter uint16
	ExtensionpageevictionsPersec  uint16
	ExtensionpagereadsPersec      uint16
	Extensionpageunreferencedtime uint16
	ExtensionpagewritesPersec     uint16
	FreeliststallsPersec          uint16
	Frequency_Object              uint16
	Frequency_PerfTime            uint16
	Frequency_Sys100NS            uint16
	IntegralControllerSlope       uint16
	LazywritesPersec              uint16
	Pagelifeexpectancy            uint16
	PagelookupsPersec             uint16
	PagereadsPersec               uint16
	PagewritesPersec              uint16
	ReadaheadpagesPersec          uint16
	ReadaheadtimePersec           uint16
	Targetpages                   uint16
	Timestamp_Object              uint16
	Timestamp_PerfTime            uint16
	Timestamp_Sys100NS            uint16
}

// Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSSQLStatistics is struct for WMI
type Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSSQLStatistics struct {
	AutoParamAttemptsPersec       uint16
	BatchRequestsPersec           uint16
	FailedAutoParamsPersec        uint16
	ForcedParameterizationsPersec uint16
	Frequency_Object              uint16
	Frequency_PerfTime            uint16
	Frequency_Sys100NS            uint16
	GuidedplanexecutionsPersec    uint16
	MisguidedplanexecutionsPersec uint16
	SafeAutoParamsPersec          uint16
	SQLAttentionrate              uint16
	SQLCompilationsPersec         uint16
	SQLReCompilationsPersec       uint16
	Timestamp_Object              uint16
	Timestamp_PerfTime            uint16
	Timestamp_Sys100NS            uint16
	UnsafeAutoParamsPersec        uint16
}

// Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSGeneralStatistics is struct for WMI
type Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSGeneralStatistics struct {
	ActiveTempTables              uint16
	ConnectionResetPersec         uint16
	EventNotificationsDelayedDrop uint16
	Frequency_Object              uint16
	Frequency_PerfTime            uint16
	Frequency_Sys100NS            uint16
	HTTPAuthenticatedRequests     uint16
	LogicalConnections            uint16
	LoginsPersec                  uint16
	LogoutsPersec                 uint16
	MarsDeadlocks                 uint16
	Nonatomicyieldrate            uint16
	Processesblocked              uint16
	SOAPEmptyRequests             uint16
	SOAPMethodInvocations         uint16
	SOAPSessionInitiateRequests   uint16
	SOAPSessionTerminateRequests  uint16
	SOAPSQLRequests               uint16
	SOAPWSDLRequests              uint16
	SQLTraceIOProviderLockWaits   uint16
	Tempdbrecoveryunitid          uint16
	Tempdbrowsetid                uint16
	TempTablesCreationRate        uint16
	TempTablesForDestruction      uint16
	Timestamp_Object              uint16
	Timestamp_PerfTime            uint16
	Timestamp_Sys100NS            uint16
	TraceEventNotificationQueue   uint16
	Transactions                  uint16
	UserConnections               uint16
}

// Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSAccessMethods is struct for WMI
type Win32_PerfRawData_MSSQLSQLEXPRESS_MSSQLSQLEXPRESSAccessMethods struct {
	AUcleanupbatchesPersec        uint16
	AUcleanupsPersec              uint16
	ByreferenceLobCreateCount     uint16
	ByreferenceLobUseCount        uint16
	CountLobReadahead             uint16
	CountPullInRow                uint16
	CountPushOffRow               uint16
	DeferreddroppedAUs            uint16
	DeferredDroppedrowsets        uint16
	DroppedrowsetcleanupsPersec   uint16
	DroppedrowsetsskippedPersec   uint16
	ExtentDeallocationsPersec     uint16
	ExtentsAllocatedPersec        uint16
	FailedAUcleanupbatchesPersec  uint16
	Failedleafpagecookie          uint16
	Failedtreepagecookie          uint16
	ForwardedRecordsPersec        uint16
	FreeSpacePageFetchesPersec    uint16
	FreeSpaceScansPersec          uint16
	Frequency_Object              uint16
	Frequency_PerfTime            uint16
	Frequency_Sys100NS            uint16
	FullScansPersec               uint16
	IndexSearchesPersec           uint16
	InSysXactwaitsPersec          uint16
	LobHandleCreateCount          uint16
	LobHandleDestroyCount         uint16
	LobSSProviderCreateCount      uint16
	LobSSProviderDestroyCount     uint16
	LobSSProviderTruncationCount  uint16
	MixedpageallocationsPersec    uint16
	PagecompressionattemptsPersec uint16
	PageDeallocationsPersec       uint16
	PagesAllocatedPersec          uint16
	PagescompressedPersec         uint16
	PageSplitsPersec              uint16
	ProbeScansPersec              uint16
	RangeScansPersec              uint16
	ScanPointRevalidationsPersec  uint16
	SkippedGhostedRecordsPersec   uint16
	TableLockEscalationsPersec    uint16
	Timestamp_Object              uint16
	Timestamp_PerfTime            uint16
	Timestamp_Sys100NS            uint16
	Usedleafpagecookie            uint16
	Usedtreepagecookie            uint16
	WorkfilesCreatedPersec        uint16
	WorktablesCreatedPersec       uint16
	WorktablesFromCacheRatio      uint16
	WorktablesFromCacheRatio_Base uint16
}
