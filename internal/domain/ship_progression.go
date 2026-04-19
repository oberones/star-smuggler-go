package domain

type ShipProgressionState struct {
	PurchasedUpgradeIDs []string
	SpecializationFlags map[string]bool
}

func NewShipProgressionState() ShipProgressionState {
	return ShipProgressionState{
		PurchasedUpgradeIDs: []string{},
		SpecializationFlags: make(map[string]bool),
	}
}

func (s ShipProgressionState) HasUpgrade(upgradeID string) bool {
	for _, ownedUpgradeID := range s.PurchasedUpgradeIDs {
		if ownedUpgradeID == upgradeID {
			return true
		}
	}
	return false
}

func (s ShipProgressionState) HasSpecialization(specialization ShipSpecialization) bool {
	return s.SpecializationFlags[string(specialization)]
}
