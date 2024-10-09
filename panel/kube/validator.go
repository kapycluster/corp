package kube

import "fmt"

func ValidateControlPlane(cp ControlPlane) error {
	if cp.Name == "" {
		return fmt.Errorf("Name can't be empty")
	}

	// These should never be empty
	if cp.ID == "" {
		return fmt.Errorf("ID is required")
	}

	// if cp.Version == "" {
	// 	return fmt.Errorf("Version is required")
	// }

	if cp.UserID == "" {
		return fmt.Errorf("User ID is required")
	}

	return nil
}
