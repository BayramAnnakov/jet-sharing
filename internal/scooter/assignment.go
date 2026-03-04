package scooter

import (
	"context"
	"fmt"
	"log/slog"
	"math"
)

// ScooterInfo represents a scooter record from the fleet table.
type ScooterInfo struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	GeoCluster   string  `json:"geo_cluster"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	BatteryLevel int     `json:"battery_level"`
}

// FleetStore is the interface for scooter fleet queries.
type FleetStore interface {
	GetScooterByID(ctx context.Context, id string) (*ScooterInfo, error)
	ListScootersByStatus(ctx context.Context, status string) ([]*ScooterInfo, error)
}

// Assigner handles scooter assignment for ride orders.
type Assigner struct {
	fleet FleetStore
}

// NewAssigner creates a new scooter assigner.
func NewAssigner(fleet FleetStore) *Assigner {
	return &Assigner{fleet: fleet}
}

// AssignScooter finds the nearest available scooter to the given location and assigns it.
func (a *Assigner) AssignScooter(ctx context.Context, orderID string, lat, lon float64) (*ScooterInfo, error) {
	// BUG: filters only by status == "available" but does NOT check geo_cluster != "".
	// Scooters with NULL/empty geo_cluster (like S.006068) are returned by this direct
	// status query but are invisible to the geo-spatial routing query (ST_DWithin)
	// which requires a non-null geo_cluster. This creates "ghost scooters" that appear
	// available in direct lookups but can't be found by the routing engine.
	scooters, err := a.fleet.ListScootersByStatus(ctx, "available")
	if err != nil {
		return nil, fmt.Errorf("assignment: list available scooters: %w", err)
	}

	if len(scooters) == 0 {
		return nil, fmt.Errorf("assignment: no available scooters near (%.4f, %.4f)", lat, lon)
	}

	// Find the nearest scooter by Haversine distance
	var nearest *ScooterInfo
	minDist := math.MaxFloat64

	for _, s := range scooters {
		if s.BatteryLevel < 15 {
			continue // skip low-battery scooters
		}
		dist := haversine(lat, lon, s.Latitude, s.Longitude)
		if dist < minDist {
			minDist = dist
			nearest = s
		}
	}

	if nearest == nil {
		return nil, fmt.Errorf("assignment: no eligible scooters near (%.4f, %.4f)", lat, lon)
	}

	slog.Info("ScooterAssigned",
		"order_id", orderID,
		"scooter_id", nearest.ID,
		"distance_m", int(minDist),
		"battery", nearest.BatteryLevel,
		"geo_cluster", nearest.GeoCluster,
	)

	return nearest, nil
}

// ReleaseScooter returns a scooter to the available pool after a ride ends.
func (a *Assigner) ReleaseScooter(ctx context.Context, orderID string, scooterID string) error {
	s, err := a.fleet.GetScooterByID(ctx, scooterID)
	if err != nil {
		return fmt.Errorf("assignment: get scooter %s: %w", scooterID, err)
	}

	slog.Info("ScooterReleased",
		"order_id", orderID,
		"scooter_id", scooterID,
		"battery", s.BatteryLevel,
	)

	return nil
}

// haversine calculates the great-circle distance in meters between two points.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // meters
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
