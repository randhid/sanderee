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
	blockDims        = r3.Vector{X: 38, Y: 70, Z: 270}
	totalLength      = 105.
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
	// internal clamp

	// pose measured from onshape CAD defined as distance from face in contact with the ur5e end effector
	// to the middle point of the internal clamp height dimension
	ipose := spatialmath.NewPose(
		r3.Vector{Z: -11}, &spatialmath.OrientationVectorDegrees{OY: 1},
	)
	// internal clamp total dims are L:300mm W: 110mm H:50mm
	// this makes a capsule to match the internal clamps length with a best fit radius
	internal, err := spatialmath.NewCapsule(ipose, 27.5, 245, "internal-clamp")
	if err != nil {
		return nil, err
	}

	// pose measured from CAD defined as distance from face in contact with the ur5e end effector
	// to furthest face of the pivot
	ppose := spatialmath.NewPose(
		r3.Vector{Z: 51.475}, &spatialmath.OrientationVectorDegrees{OY: 1},
	)
	// pivot total dims are L300mm W: 110mm H:50mm
	// this makes a capsule to match the pivot length with a best fit radius
	pivot, err := spatialmath.NewCapsule(ppose, 40, 220, "pivot")
	if err != nil {
		return nil, err
	}

	// hose - ballpark placement is middle of the sanding block height block at one edge
	hpose := spatialmath.NewPose(
		r3.Vector{Y: blockDims.Z / 2, Z: totalLength - blockDims.X/2}, &spatialmath.OrientationVectorDegrees{OY: 1},
	)
	hose, err := spatialmath.NewSphere(hpose, 25, "hose")

	// sanding block dims are L: 270mm W: 78mm H:38mm (X)
	bpose := spatialmath.NewPose(
		r3.Vector{Z: totalLength}, &spatialmath.OrientationVectorDegrees{OY: 1},
	)
	block, err := spatialmath.NewBox(bpose, blockDims, "block")
	if err != nil {
		return nil, err
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
