package v1

//==============================================================================
// Constants
//==============================================================================
const (
	KindNameArcus = "Arcus"
)

const (
	PortNameZkClient         = "client"
	PortNameZkServer         = "server"
	PortNameZkLeaderElection = "leader-election"
)

const (
	PathZkScripts = "/scripts"
)

const (
	FileZkInitScript = "arcus-zk-init.sh"
	FileZkOkScript   = "arcus-zk-ok.sh"
)

const (
	ProbeInitialDelaySecondsZk = 10
	ProbeTimeoutSecondsZk      = 10
)

const (
	PDBMaxUnavailble = 1
)

const (
	LabelKeyApp  = "app"
	LabelValueZk = "arcus-zk"
	LabelValueMc = "arcus-mc"
)
