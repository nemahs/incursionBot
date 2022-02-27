package main

import (
	"fmt"
)

// TODO: Figure out what to do with this enum
type IncursionState string

const (
	Established IncursionState = "Established"
	Mobilizing  IncursionState = "Mobilizing"
	Withdrawing	IncursionState = "Withdrawing"
)

type NamedItem struct {
	Name string
	ID int
}

type Incursion struct {
	Constellation NamedItem				// Constellation the incursion is in
	StagingSystem NamedItem				// Name of the HQ system
	Influence     float64					// Influence of the incursion from 0 to 1 inclusive
	Region        NamedItem				// Region the incursion is in
	State         IncursionState	// Current state of the incursion
	Security      SecurityClass		// Security type of the staging system
	SecStatus     float64					// Security status of the staging system, -1 to 1 inclusive
	Distance      int							// Distance from home system
}

func (inc *Incursion) ToString() string {
	return fmt.Sprintf("%s {%.2f} (%s - %s)", inc.StagingSystem.Name, inc.SecStatus, inc.Constellation.Name, inc.Region.Name)
}

type IncursionList []Incursion
func (list *IncursionList) find(stagingId int) *Incursion {
  for _, incursion := range *list {
    if incursion.StagingSystem.ID == stagingId { return &incursion }
  }
  return nil
}

// Updates the give incursion wih new data. Returns true if the state changed, False otherwise.
func UpdateIncursion(incursion *Incursion, newData IncursionResponse) bool {
  updated := false

  if incursion.State != newData.State {
    incursion.State = newData.State
    updated = true
  }

  incursion.Influence = newData.Influence
  return updated
}

// Creates a new Incursion object from ESI data
func CreateNewIncursion(incursion IncursionResponse) Incursion {
  stagingData := getSystemInfo(incursion.StagingID)
  constData := getConstInfo(incursion.ConstellationID)
  names := getNames([]int{constData.RegionID})	
  distance := GetRouteLength(homeSystem, incursion.StagingID)

  newIncursion := Incursion{
    Constellation: NamedItem{ID: constData.ID, Name: constData.Name},
    StagingSystem: NamedItem{ID: stagingData.ID, Name: stagingData.Name},
    Influence: incursion.Influence,
    Region: NamedItem{ID: constData.RegionID, Name: names[constData.RegionID]},
    State: incursion.State,
    SecStatus: stagingData.SecStatus,
    Security: stagingData.SecurityClass,
    Distance: distance,
  }

  return newIncursion
}