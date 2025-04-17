package sanderee

import (
	"context"
	"errors"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/gripper"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

var (
	SanderEe         = resource.NewModel("rand", "sander-ee", "sander-ee")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(gripper.API, SanderEe,
		resource.Registration[gripper.Gripper, *resource.NoNativeConfig]{
			Constructor: NewSander,
		},
	)
}

type sanderEeSanderEe struct {
	resource.AlwaysRebuild
	resource.Named
	resource.TriviallyCloseable

	logger logging.Logger
	geoms  []spatialmath.Geometry
}

func NewSander(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (gripper.Gripper, error) {

	var errs error
	ipose := spatialmath.NewPose(
		r3.Vector{Z: -14.5}, &spatialmath.OrientationVectorDegrees{OY: 1},
	)
	internal, err := spatialmath.NewCapsule(ipose, 55, 300, "internal-clamp")
	errors.Join(errs, err)
	ppose := spatialmath.NewPose(
		r3.Vector{Z: 51.475}, &spatialmath.OrientationVectorDegrees{OY: 1},
	)
	pivot, err := spatialmath.NewCapsule(ppose, 80, 300, "pivot")
	errors.Join(errs, err)
	hpose := spatialmath.NewPose(
		r3.Vector{Y: 78, Z: 105}, &spatialmath.OrientationVectorDegrees{OY: 1},
	)
	hose, err := spatialmath.NewSphere(hpose, 50, "hose")
	bpose := spatialmath.NewPose(
		r3.Vector{Z: 51.475}, &spatialmath.OrientationVectorDegrees{OY: 1},
	)
	block, err := spatialmath.NewBox(bpose, r3.Vector{X: 270, Y: 70, Z: 38}, "block")
	errors.Join(errs, err)
	if errs != nil {
		return nil, errs
	}

	geoms := []spatialmath.Geometry{
		internal,
		pivot,
		block,
		hose,
	}

	s := &sanderEeSanderEe{
		Named:  rawConf.ResourceName().AsNamed(),
		logger: logger,
		geoms:  geoms,
	}
	return s, nil
}

func (s *sanderEeSanderEe) Open(ctx context.Context, extra map[string]interface{}) error {
	return nil
}

func (s *sanderEeSanderEe) Grab(ctx context.Context, extra map[string]interface{}) (bool, error) {
	return false, nil
}

func (s *sanderEeSanderEe) Stop(ctx context.Context, extra map[string]interface{}) error {
	return nil
}

func (s *sanderEeSanderEe) ModelFrame() referenceframe.Model {
	return referenceframe.NewSimpleModel(s.Name().Name)
}

func (s *sanderEeSanderEe) IsMoving(ctx context.Context) (bool, error) {
	return false, nil
}

func (s *sanderEeSanderEe) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	return s.geoms, nil
}
