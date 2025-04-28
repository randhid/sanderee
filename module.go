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
	blockDims        = r3.Vector{X: 38, Y: 78, Z: 270}
	totalLengthFancy = 105.
	hoseRadius       = 12.5

	clampName        = "clamp"
	hoseName         = "hose"
	blockName        = "block"
	pivotName        = "pivot"
	sandingBlockName = "sanding-block"
)

func init() {
	resource.RegisterComponent(gripper.API, SanderEe,
		resource.Registration[gripper.Gripper, *Config]{
			Constructor: NewSander,
		},
	)
}

type Config struct {
	resource.TriviallyValidateConfig
	UseCapsules bool `json:"use_capsules"`
	FancySander bool `json:"fancy_sander"`
}

type sanderEeSanderEe struct {
	resource.AlwaysRebuild
	resource.Named
	resource.TriviallyCloseable

	logger logging.Logger
	geoms  []spatialmath.Geometry
}

func NewSander(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (gripper.Gripper, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	geoms := []spatialmath.Geometry{}
	switch {
	case conf.UseCapsules:
		sanderGeoms, err := makeSanderCapsules()
		if err != nil {
			return nil, err
		}
		geoms = append(geoms, sanderGeoms...)
	case conf.FancySander:
		sanderGeoms, err := makeFancySander()
		if err != nil {
			return nil, err
		}
		geoms = append(geoms, sanderGeoms...)
	default:
		sanderGeoms, err := makeSanderBlocks()
		if err != nil {
			return nil, err
		}
		geoms = append(geoms, sanderGeoms...)
	}

	s := &sanderEeSanderEe{
		Named:  rawConf.ResourceName().AsNamed(),
		logger: logger,
		geoms:  geoms,
	}
	return s, nil
}

func (s *sanderEeSanderEe) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	return s.geoms, nil
}

var (
	clampBoxDim        = r3.Vector{X: 340, Y: 130, Z: 51}
	sandingBlockBoxDim = r3.Vector{X: 270, Y: 78, Z: 25}
)

func makeSanderBlocks() ([]spatialmath.Geometry, error) {
	clampPose := spatialmath.NewPoseFromPoint(r3.Vector{Z: clampBoxDim.Z})
	clamp, err := spatialmath.NewBox(clampPose, clampBoxDim, clampName)
	if err != nil {
		return nil, err
	}

	sandingBlockPose := spatialmath.NewPoseFromPoint(
		clampPose.Point().Add(r3.Vector{Z: clampBoxDim.Z/2 + sandingBlockBoxDim.Z/2}))
	sandingBlock, err := spatialmath.NewBox(sandingBlockPose, sandingBlockBoxDim, sandingBlockName)
	if err != nil {
		return nil, err
	}

	hosePose := spatialmath.NewPoseFromPoint(
		sandingBlockPose.Point().Add(r3.Vector{X: sandingBlockBoxDim.X / 2}),
	)
	hose, err := spatialmath.NewSphere(hosePose, 15, hoseName)
	if err != nil {
		return nil, err
	}

	return []spatialmath.Geometry{
		clamp,
		sandingBlock,
		hose,
	}, nil
}

var (
	clampCapsuleLength        = 340.
	clampCapsuleRadius        = 65.0 / 2
	sandingBlockCapsuleLength = 270.
	sandingBlockCapsuleRadius = 78.0 / 2
	totalLengthCapsule        = 76.
)

func makeSanderCapsules() ([]spatialmath.Geometry, error) {
	clampPose := spatialmath.NewPoseFromOrientation(&spatialmath.OrientationVectorDegrees{OY: 1})
	clamp, err := spatialmath.NewCapsule(clampPose, clampCapsuleRadius, clampCapsuleLength, clampName)
	if err != nil {
		return nil, err
	}

	sandingBlockPose := spatialmath.NewPose(
		r3.Vector{Z: totalLengthCapsule - clampCapsuleRadius},
		&spatialmath.OrientationVectorDegrees{OY: 1},
	)
	sandingBlock, err := spatialmath.NewCapsule(sandingBlockPose, sandingBlockCapsuleRadius, sandingBlockCapsuleLength, sandingBlockName)
	if err != nil {
		return nil, err
	}

	hosePose := spatialmath.NewPoseFromPoint(r3.Vector{Y: sandingBlockCapsuleLength / 2, Z: totalLengthCapsule - sandingBlockCapsuleRadius*2})
	hose, err := spatialmath.NewSphere(hosePose, 12.5, hoseName)
	if err != nil {
		return nil, err
	}

	return []spatialmath.Geometry{
		clamp,
		hose,
		sandingBlock,
	}, nil
}

func makeFancySander() ([]spatialmath.Geometry, error) {
	// pose measured from onshape CAD defined as distance from face in contact with the ur5e end effector
	// to the middle point of the internal clamp height dimension
	ipose := spatialmath.NewPoseFromPoint(r3.Vector{Z: -11})
	// internal clamp total dims are L:300mm W: 110mm H:50mm
	// this makes a capsule to match the internal clamps length with a best fit radius
	internal, err := spatialmath.NewCapsule(ipose, 27.5, 245, clampName)
	if err != nil {
		return nil, err
	}

	// pose measured from CAD defined as distance from face in contact with the ur5e end effector
	// to furthest face of the pivot
	ppose := spatialmath.NewPose(r3.Vector{Z: 51.475}, &spatialmath.OrientationVectorDegrees{OY: 1})
	// pivot total dims are L300mm W: 110mm H:50mm
	// this makes a capsule to match the pivot length with a best fit radius
	pivot, err := spatialmath.NewCapsule(ppose, 40, 220, pivotName)
	if err != nil {
		return nil, err
	}

	// hose - ballpark placement is middle of the sanding block height block at one edge
	hpose := spatialmath.NewPose(r3.Vector{X: blockDims.X / 2, Z: totalLengthFancy - blockDims.Z/2}, &spatialmath.OrientationVectorDegrees{OY: 1})
	hose, err := spatialmath.NewSphere(hpose, 25, hoseName)
	if err != nil {
		return nil, err
	}

	// sanding block dims are L: 270mm W: 78mm H:38mm (X)
	bpose := spatialmath.NewPose(r3.Vector{Z: totalLengthFancy}, &spatialmath.OrientationVectorDegrees{OY: 1})
	block, err := spatialmath.NewBox(bpose, blockDims, sandingBlockName)
	if err != nil {
		return nil, err
	}

	return []spatialmath.Geometry{
		internal,
		pivot,
		block,
		hose,
	}, nil
}

// Unimplemented methods
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
