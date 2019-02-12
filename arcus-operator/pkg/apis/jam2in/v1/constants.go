package v1

//==============================================================================
// Constants
//==============================================================================
const (
	KindNameArcus = "Arcus"

	PortNameZkClient         = "client"
	PortNameZkServer         = "server"
	PortNameZkLeaderElection = "leader-election"

	DirPathZkScripts = "/scripts"

	FileZkInitScript = "arcus-zk-init.sh"
	FileZkOkScript   = "arcus-zk-ok.sh"

	ProbeInitialDelaySecondsZk = 10
	ProbeTimeoutSecondsZk      = 10

	PDBMaxUnavailble = 1

	LabelKeyApp  = "app"
	LabelValueZk = "arcus-zk"
	LabelValueMc = "arcus-mc"
)
