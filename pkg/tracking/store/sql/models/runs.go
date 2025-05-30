package models

import (
	"database/sql"

	"github.com/mlflow/mlflow-go-backend/pkg/entities"
	"github.com/mlflow/mlflow-go-backend/pkg/utils"
)

// Run mapped from table <runs>.
type Run struct {
	ID             string         `gorm:"column:run_uuid;primaryKey"`
	Name           string         `gorm:"column:name"`
	SourceType     SourceType     `gorm:"column:source_type"`
	SourceName     string         `gorm:"column:source_name"`
	EntryPointName string         `gorm:"column:entry_point_name"`
	UserID         string         `gorm:"column:user_id"`
	Status         RunStatus      `gorm:"column:status"`
	StartTime      int64          `gorm:"column:start_time"`
	EndTime        sql.NullInt64  `gorm:"column:end_time"`
	SourceVersion  string         `gorm:"column:source_version"`
	LifecycleStage LifecycleStage `gorm:"column:lifecycle_stage"`
	ArtifactURI    string         `gorm:"column:artifact_uri"`
	ExperimentID   int32          `gorm:"column:experiment_id"`
	DeletedTime    sql.NullInt64  `gorm:"column:deleted_time"`
	Params         []Param
	Tags           []Tag
	Metrics        []Metric
	LatestMetrics  []LatestMetric
	Inputs         []Input  `gorm:"foreignKey:DestinationID"`
	Outputs        []Output `gorm:"foreignKey:DestinationID"`
}

type RunStatus string

func (s RunStatus) String() string {
	return string(s)
}

const (
	RunStatusRunning   RunStatus = "RUNNING"
	RunStatusScheduled RunStatus = "SCHEDULED"
	RunStatusFinished  RunStatus = "FINISHED"
	RunStatusFailed    RunStatus = "FAILED"
	RunStatusKilled    RunStatus = "KILLED"
)

type SourceType string

func (s SourceType) String() string {
	return string(s)
}

const (
	SourceTypeNotebook  SourceType = "NOTEBOOK"
	SourceTypeJob       SourceType = "JOB"
	SourceTypeProject   SourceType = "PROJECT"
	SourceTypeLocal     SourceType = "LOCAL"
	SourceTypeUnknown   SourceType = "UNKNOWN"
	SourceTypeRecipe    SourceType = "RECIPE"
	SourceTypeRunInput  SourceType = "RUN_INPUT"
	SourceTypeRunOutput SourceType = "RUN_OUTPUT"
)

//nolint:funlen
func (r Run) ToEntity() *entities.Run {
	metrics := make([]*entities.Metric, 0, len(r.LatestMetrics))
	for _, metric := range r.LatestMetrics {
		metrics = append(metrics, metric.ToEntity())
	}

	params := make([]*entities.Param, 0, len(r.Params))
	for _, param := range r.Params {
		params = append(params, param.ToEntity())
	}

	tags := make([]*entities.RunTag, 0, len(r.Tags))
	for _, tag := range r.Tags {
		tags = append(tags, tag.ToEntity())
	}

	datasetInputs := make([]*entities.DatasetInput, 0, len(r.Inputs))
	for _, input := range r.Inputs {
		datasetInputs = append(datasetInputs, input.DatasetToEntity())
	}

	modelOutputs := make([]*entities.ModelOutput, 0, len(r.Outputs))
	for _, output := range r.Outputs {
		modelOutputs = append(modelOutputs, output.ToEntity())
	}

	modelInputs := make([]*entities.ModelInput, 0, len(r.Inputs))
	for _, input := range r.Inputs {
		modelInputs = append(modelInputs, input.ModelInputToEntity())
	}

	var endTime *int64
	if r.EndTime.Valid {
		endTime = utils.PtrTo(r.EndTime.Int64)
	}

	return &entities.Run{
		Info: &entities.RunInfo{
			RunID:          r.ID,
			RunUUID:        r.ID,
			RunName:        r.Name,
			ExperimentID:   r.ExperimentID,
			UserID:         r.UserID,
			Status:         r.Status.String(),
			StartTime:      r.StartTime,
			EndTime:        endTime,
			ArtifactURI:    r.ArtifactURI,
			LifecycleStage: r.LifecycleStage.String(),
		},
		Data: &entities.RunData{
			Tags:    tags,
			Params:  params,
			Metrics: metrics,
		},
		Inputs: &entities.RunInputs{
			ModelInputs:   modelInputs,
			DatasetInputs: datasetInputs,
		},
		Outputs: &entities.RunOutputs{
			ModelOutputs: modelOutputs,
		},
	}
}
