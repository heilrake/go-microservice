package service

import (
	"fmt"
	math "math/rand/v2"
	pb "ride-sharing/shared/proto/driver"
	"ride-sharing/shared/util"
	"sync"

	"github.com/mmcloughlin/geohash"
)

type driverInMap struct {
	Driver *pb.Driver
	// Index int
	// TODO: route
}

type DriverService interface {
	RegisterDriver(driverID, packageSlug string) (*pb.Driver, error)
	UnregisterDriver(driverID string) error
}

type driverService struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

func NewService() DriverService {
	return &driverService{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *driverService) RegisterDriver(driverID, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(util.PredefinedRoutes))
	randomRoute := util.PredefinedRoutes[randomIndex]

	randomPlate := util.GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIndex)

	// we can ignore this property for now, but it must be sent to the frontend.
	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &pb.Driver{
		Id:             driverID,
		Geohash:        geohash,
		Location:       &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
		Name:           "Lando Norris",
		PackageSlug:    packageSlug,
		ProfilePicture: randomAvatar,
		CarPlate:       randomPlate,
	}

	s.drivers = append(s.drivers, &driverInMap{
		Driver: driver,
	})

	fmt.Println("Driver registered: ", driver)

	return driver, nil
}

func (s *driverService) UnregisterDriver(driverID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, driver := range s.drivers {
		if driver.Driver.Id == driverID {
			s.drivers = append(s.drivers[:i], s.drivers[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("driver not found")
}
