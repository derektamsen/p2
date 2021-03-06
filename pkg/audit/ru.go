package audit

import (
	"encoding/json"

	"github.com/square/p2/pkg/labels"
	"github.com/square/p2/pkg/manifest"
	pc_fields "github.com/square/p2/pkg/pc/fields"
	roll_fields "github.com/square/p2/pkg/roll/fields"
	"github.com/square/p2/pkg/types"
	"github.com/square/p2/pkg/util"
)

const (
	RUCreationEvent   EventType = "ROLLING_UPDATE_CREATION"
	RUCompletionEvent EventType = "ROLLING_UPDATE_COMPLETION"
)

type RUCreationDetails struct {
	PodID            types.PodID                `json:"pod_id"`
	AvailabilityZone pc_fields.AvailabilityZone `json:"availability_zone"`
	ClusterName      pc_fields.ClusterName      `json:"cluster_name"`
	Deployer         string                     `json:"deployer"`
	Manifest         string                     `json:"manifest"`
	RollingUpdateID  roll_fields.ID             `json:"rolling_update_id"`
}

type RUCompletionDetails struct {
	PodID            types.PodID                `json:"pod_id"`
	AvailabilityZone pc_fields.AvailabilityZone `json:"availability_zone"`
	ClusterName      pc_fields.ClusterName      `json:"cluster_name"`
	RollingUpdateID  roll_fields.ID             `json:"rolling_update_id"`
	Succeeded        bool                       `json:"succeeded"`
	Canceled         bool                       `json:"canceled"`
}

func NewRUCreationEventDetails(
	podID types.PodID,
	az pc_fields.AvailabilityZone,
	name pc_fields.ClusterName,
	deployer string,
	manifest manifest.Manifest,
	rollingUpdateID roll_fields.ID,
) (json.RawMessage, error) {
	manifestBytes, err := manifest.Marshal()
	if err != nil {
		return nil, err
	}

	details := RUCreationDetails{
		PodID:            podID,
		AvailabilityZone: az,
		ClusterName:      name,
		Deployer:         deployer,
		Manifest:         string(manifestBytes),
		RollingUpdateID:  rollingUpdateID,
	}

	bytes, err := json.Marshal(details)
	if err != nil {
		return nil, util.Errorf("could not marshal ru creation details as json: %s", err)
	}

	return json.RawMessage(bytes), nil
}

type Labeler interface {
	GetLabels(labelType labels.Type, id string) (labels.Labeled, error)
}

func NewRUCompletionEventDetails(
	rollingUpdateID roll_fields.ID,
	succeeded bool,
	canceled bool,
	labeler Labeler,
) (json.RawMessage, error) {
	details := RUCompletionDetails{
		RollingUpdateID: rollingUpdateID,
		Succeeded:       succeeded,
		Canceled:        canceled,
	}

	labels, err := labeler.GetLabels(labels.RU, rollingUpdateID.String())
	if err != nil {
		return nil, util.Errorf("could not determine pod cluster for RU %s: %s", rollingUpdateID, err)
	}

	details.PodID = types.PodID(labels.Labels[pc_fields.PodIDLabel])
	details.AvailabilityZone = pc_fields.AvailabilityZone(labels.Labels[pc_fields.AvailabilityZoneLabel])
	details.ClusterName = pc_fields.ClusterName(labels.Labels[pc_fields.ClusterNameLabel])

	bytes, err := json.Marshal(details)
	if err != nil {
		return nil, util.Errorf("could not marshal ru completion details as json: %s", err)
	}

	return json.RawMessage(bytes), nil
}
