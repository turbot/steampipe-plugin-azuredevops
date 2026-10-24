package azuredevops

import (
	"context"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const userEntitlementsAPIVersion = "7.1"

// https://learn.microsoft.com/en-us/rest/api/azure/devops/memberentitlementmanagement/user-entitlements/search-user-entitlements
func tableAzureDevOpsUserEntitlement(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_user_entitlement",
		Description: "Users of the Azure DevOps organization with their access level, status and when they last signed in.",
		List: &plugin.ListConfig{
			Hydrate: listUserEntitlements,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "access_level", Description: "Display name of the license, e.g. Basic, Stakeholder, Visual Studio Enterprise subscription.", Type: proto.ColumnType_STRING},
			{Name: "account_license_type", Description: "Account license type: express (Basic), stakeholder, advanced, professional, earlyAdopter or none.", Type: proto.ColumnType_STRING},
			{Name: "date_created", Description: "When the user was added to the organization.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "descriptor", Description: "Graph descriptor of the user.", Type: proto.ColumnType_STRING},
			{Name: "display_name", Description: "Display name of the user.", Type: proto.ColumnType_STRING},
			{Name: "id", Description: "ID of the user.", Type: proto.ColumnType_STRING},
			{Name: "last_accessed_date", Description: "When the user last accessed the organization. Null if they never have.", Type: proto.ColumnType_TIMESTAMP},
			{Name: "licensing_source", Description: "Where the license comes from: account, msdn, profile, auto or trial.", Type: proto.ColumnType_STRING},
			{Name: "mail_address", Description: "Email address of the user.", Type: proto.ColumnType_STRING},
			{Name: "msdn_license_type", Description: "Visual Studio subscription type when licensing_source is msdn.", Type: proto.ColumnType_STRING},
			{Name: "origin", Description: "Identity provider of the user, e.g. aad or msa.", Type: proto.ColumnType_STRING},
			{Name: "principal_name", Description: "Principal name of the user, usually the sign-in email.", Type: proto.ColumnType_STRING},
			{Name: "status", Description: "Status of the user in the organization: active, pending, disabled, deleted, expired or pendingDisabled.", Type: proto.ColumnType_STRING},
			{Name: "subject_kind", Description: "Kind of graph subject, normally user.", Type: proto.ColumnType_STRING},
			{Name: "user_type", Description: "Type of the user in the organization: member or guest.", Type: proto.ColumnType_STRING},
			{Name: "title", Description: "Title of the resource.", Type: proto.ColumnType_STRING, Transform: transform.FromField("DisplayName")},
		}),
	}
}

type userEntitlementResponse struct {
	Id   string `json:"id"`
	User struct {
		DisplayName   string `json:"displayName"`
		PrincipalName string `json:"principalName"`
		MailAddress   string `json:"mailAddress"`
		Origin        string `json:"origin"`
		MetaType      string `json:"metaType"`
		SubjectKind   string `json:"subjectKind"`
		Descriptor    string `json:"descriptor"`
	} `json:"user"`
	AccessLevel struct {
		LicenseDisplayName string `json:"licenseDisplayName"`
		AccountLicenseType string `json:"accountLicenseType"`
		LicensingSource    string `json:"licensingSource"`
		MsdnLicenseType    string `json:"msdnLicenseType"`
		Status             string `json:"status"`
	} `json:"accessLevel"`
	LastAccessedDate adoTime `json:"lastAccessedDate"`
	DateCreated      adoTime `json:"dateCreated"`
}

type UserEntitlement struct {
	Id                 string
	DisplayName        string
	PrincipalName      string
	MailAddress        string
	Origin             string
	UserType           string
	SubjectKind        string
	Descriptor         string
	AccessLevel        string
	AccountLicenseType string
	LicensingSource    string
	MsdnLicenseType    string
	Status             string
	LastAccessedDate   *time.Time
	DateCreated        *time.Time
}

func listUserEntitlements(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_user_entitlement.listUserEntitlements", "connection_error", err)
		return nil, err
	}
	err = readUserEntitlements(ctx, client, func(u UserEntitlement) bool {
		d.StreamListItem(ctx, u)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_user_entitlement.listUserEntitlements", "api_error", err)
		return nil, err
	}
	return nil, nil
}

// readUserEntitlements pages with the continuation token in the response body; api-version 7.1
// returns the page as items, earlier versions as members.
func readUserEntitlements(ctx context.Context, client *restClient, emit func(UserEntitlement) bool) error {
	endpoint := serviceURL(client.baseURL, "vsaex") + "/_apis/userentitlements"
	params := apiVersion(userEntitlementsAPIVersion)
	seen := map[string]bool{}
	for {
		var page struct {
			Items             []userEntitlementResponse `json:"items"`
			Members           []userEntitlementResponse `json:"members"`
			ContinuationToken string                    `json:"continuationToken"`
		}
		if err := client.getJSON(ctx, endpoint+"?"+params.Encode(), scopeMemberEntitlement, &page); err != nil {
			return err
		}
		for _, u := range append(page.Items, page.Members...) {
			if !emit(userEntitlementRow(u)) {
				return nil
			}
		}
		if page.ContinuationToken == "" || seen[page.ContinuationToken] {
			return nil
		}
		seen[page.ContinuationToken] = true
		params.Set("continuationToken", page.ContinuationToken)
	}
}

func userEntitlementRow(u userEntitlementResponse) UserEntitlement {
	return UserEntitlement{
		Id:                 u.Id,
		DisplayName:        u.User.DisplayName,
		PrincipalName:      u.User.PrincipalName,
		MailAddress:        u.User.MailAddress,
		Origin:             u.User.Origin,
		UserType:           u.User.MetaType,
		SubjectKind:        u.User.SubjectKind,
		Descriptor:         u.User.Descriptor,
		AccessLevel:        u.AccessLevel.LicenseDisplayName,
		AccountLicenseType: u.AccessLevel.AccountLicenseType,
		LicensingSource:    u.AccessLevel.LicensingSource,
		MsdnLicenseType:    u.AccessLevel.MsdnLicenseType,
		Status:             u.AccessLevel.Status,
		LastAccessedDate:   u.LastAccessedDate.ptr(),
		DateCreated:        u.DateCreated.ptr(),
	}
}
