package azuredevops

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const auditStreamsAPIVersion = "7.1-preview.1"

// https://learn.microsoft.com/en-us/rest/api/azure/devops/audit/streams/query-all-streams
func tableAzureDevOpsAuditStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_audit_stream",
		Description: "Audit streams that forward the Azure DevOps organization's audit log to a SIEM or log store. Stream credentials are never exposed.",
		List: &plugin.ListConfig{
			Hydrate: listAuditStreams,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "consumer_type", Description: "Destination type, e.g. Splunk, AzureMonitorLogs, AzureEventGrid.", Type: proto.ColumnType_STRING},
			{Name: "created_time", Description: "When the stream was created.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "display_name", Description: "Display name of the stream, usually its destination.", Type: proto.ColumnType_STRING},
			{Name: "id", Description: "ID of the stream.", Type: proto.ColumnType_INT},
			{Name: "status", Description: "Status of the stream: unknown, enabled, disabledByUser, disabledBySystem, deleted or backfilling.", Type: proto.ColumnType_STRING},
			{Name: "status_reason", Description: "Why the stream has its status, e.g. invalid credentials.", Type: proto.ColumnType_STRING},
			{Name: "updated_time", Description: "When the stream was last updated.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "title", Description: "Title of the resource.", Type: proto.ColumnType_STRING, Transform: transform.FromField("DisplayName")},
		}),
	}
}

// auditStreamResponse leaves out consumerInputs, which holds the destination's credentials.
type auditStreamResponse struct {
	Id           int             `json:"id"`
	DisplayName  string          `json:"displayName"`
	ConsumerType string          `json:"consumerType"`
	Status       json.RawMessage `json:"status"`
	StatusReason *string         `json:"statusReason"`
	CreatedTime  json.RawMessage `json:"createdTime"`
	UpdatedTime  json.RawMessage `json:"updatedTime"`
}

type AuditStream struct {
	Id           int
	DisplayName  string
	ConsumerType string
	Status       string
	StatusReason *string
	CreatedTime  *time.Time
	UpdatedTime  *time.Time
}

var auditStreamStatuses = []string{"unknown", "enabled", "disabledByUser", "disabledBySystem", "deleted", "backfilling"}

func listAuditStreams(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_audit_stream.listAuditStreams", "connection_error", err)
		return nil, err
	}
	err = readAuditStreams(ctx, client, func(s AuditStream) bool {
		d.StreamListItem(ctx, s)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_audit_stream.listAuditStreams", "api_error", err)
		return nil, err
	}
	return nil, nil
}

func readAuditStreams(ctx context.Context, client *restClient, emit func(AuditStream) bool) error {
	endpoint := serviceURL(client.baseURL, "auditservice") + "/_apis/audit/streams"
	return getList(ctx, client, endpoint, apiVersion(auditStreamsAPIVersion), scopeAuditLog, func(s auditStreamResponse) bool {
		return emit(AuditStream{
			Id:           s.Id,
			DisplayName:  s.DisplayName,
			ConsumerType: s.ConsumerType,
			Status:       auditStreamStatus(s.Status),
			StatusReason: s.StatusReason,
			CreatedTime:  auditTime(s.CreatedTime),
			UpdatedTime:  auditTime(s.UpdatedTime),
		})
	})
}

// auditStreamStatus accepts the status as its enum name or, as the API reference example shows, its number.
func auditStreamStatus(raw json.RawMessage) string {
	var name string
	if json.Unmarshal(raw, &name) == nil {
		return name
	}
	var n int
	if json.Unmarshal(raw, &n) == nil && n >= 0 && n < len(auditStreamStatuses) {
		return auditStreamStatuses[n]
	}
	return ""
}

// auditTime reads an ISO 8601 timestamp or the legacy "/Date(<ms>)/" form the API reference shows.
func auditTime(raw json.RawMessage) *time.Time {
	var s string
	if json.Unmarshal(raw, &s) != nil || s == "" {
		return nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return &t
	}
	if len(s) > len("/Date()/") && s[:6] == "/Date(" && s[len(s)-2:] == ")/" {
		if ms, err := strconv.ParseInt(s[6:len(s)-2], 10, 64); err == nil {
			t := time.UnixMilli(ms).UTC()
			return &t
		}
	}
	return nil
}
