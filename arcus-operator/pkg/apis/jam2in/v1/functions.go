package v1

import (
	"fmt"
)

//==============================================================================
// Public Functions
//==============================================================================
func GetObjectNameConfigMap(arcus *Arcus) string {
	return fmt.Sprintf("%s-cm", arcus.Name)
}

func GetObjectNameZkHeadlessService(arcus *Arcus) string {
	return fmt.Sprintf("%s-zk-headless-svc", arcus.Name)
}

func GetObjectNameZkStatefulSet(arcus *Arcus) string {
	return fmt.Sprintf("%s-zk", arcus.Name)
}

func GetObjectNameZkVolume(arcus *Arcus) string {
	return fmt.Sprintf("%s-zk-vol", arcus.Name)
}

func GetZkLabel(arcus *Arcus) map[string]string {
	return map[string]string{
		"app": GetObjectNameZkStatefulSet(arcus),
	}
}

func GetDefaultMode() *int32 {
	defaultMode := new(int32)
	*defaultMode = int32(0755)
	return defaultMode
}
