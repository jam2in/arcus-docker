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

func GetObjectNamePodDisruptionBudget(arcus *Arcus) string {
	return fmt.Sprintf("%s-pdb", arcus.Name)
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
