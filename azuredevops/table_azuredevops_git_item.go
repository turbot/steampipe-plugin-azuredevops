package azuredevops

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

const gitItemsAPIVersion = "7.1"

// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/items
func tableAzureDevOpsGitItem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_git_item",
		Description: "Files and folders in an Azure DevOps Git repository at a version, and the content of a single file.",
		List: &plugin.ListConfig{
			// Exact cache matching: these quals change what is fetched, so a cached listing must never
			// answer a single-file read or a different version.
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project", Require: plugin.Required},
				{Name: "repository_id", Require: plugin.Required},
				{Name: "version", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "version_type", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "scope_path", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "path", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Hydrate: listGitItems,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "project",
				Description: "Name or ID of the project the repository belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project"),
			},
			{
				Name:        "repository_id",
				Description: "ID or name of the repository.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("repository_id"),
			},
			{
				Name:        "version",
				Description: "The branch, tag or commit read, as given. Defaults to the repository's default branch.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("version"),
			},
			{
				Name:        "version_type",
				Description: "How version is read: branch, tag or commit. Defaults to commit for a 40-character hash and branch otherwise.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("version_type"),
			},
			{
				Name:        "scope_path",
				Description: "The folder whose whole subtree is listed. Defaults to the repository root.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("scope_path"),
			},
			{
				Name:        "path",
				Description: "Path of the item from the repository root, starting with /. Given as a qualifier, reads that one file and its content.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "commit_id",
				Description: "ID of the commit version resolved to, the same on every row of a listing; on the root / row it is the head of the branch.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content",
				Description: "Content of the file. Set only when path is given; null for folders and binary files.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "git_object_type",
				Description: "The Git object type: blob for a file, tree for a folder, commit for a submodule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_folder",
				Description: "True if the item is a folder.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_sym_link",
				Description: "True if the item is a symbolic link.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "object_id",
				Description: "Git object ID of the item.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "Size in bytes of a text file read with path. Null in listings, which do not report sizes, and for binary files.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "url",
				Description: "REST URL of the item.",
				Type:        proto.ColumnType_STRING,
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Path"),
			},
		}),
	}
}

type GitItem struct {
	ObjectId        string              `json:"objectId"`
	GitObjectType   string              `json:"gitObjectType"`
	CommitId        string              `json:"commitId"`
	Path            string              `json:"path"`
	IsFolder        bool                `json:"isFolder"`
	IsSymLink       bool                `json:"isSymLink"`
	Url             string              `json:"url"`
	Content         *string             `json:"content"`
	ContentMetadata *gitContentMetadata `json:"contentMetadata"`
	Size            *int64              `json:"-"`
}

type gitContentMetadata struct {
	IsBinary bool `json:"isBinary"`
}

type gitItemList struct {
	Count int       `json:"count"`
	Value []GitItem `json:"value"`
}

type gitItemRequest struct {
	Project      string
	RepositoryID string
	Version      string
	VersionType  string
	ScopePath    string
	Path         string
}

var fullCommitID = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)

func listGitItems(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	req := gitItemRequest{
		Project:      d.EqualsQualString("project"),
		RepositoryID: d.EqualsQualString("repository_id"),
		Version:      d.EqualsQualString("version"),
		VersionType:  d.EqualsQualString("version_type"),
		ScopePath:    d.EqualsQualString("scope_path"),
		Path:         d.EqualsQualString("path"),
	}
	if req.Project == "" || req.RepositoryID == "" {
		return nil, nil
	}

	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_item.listGitItems", "connection_error", err)
		return nil, err
	}

	err = readGitItems(ctx, client, req, func(item GitItem) bool {
		d.StreamListItem(ctx, item)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_git_item.listGitItems", "api_error", err)
		return nil, err
	}
	return nil, nil
}

// readGitItems lists a subtree in one request (recursionLevel=Full), or reads one item with its
// content, calling emit per row until it returns false.
func readGitItems(ctx context.Context, client *restClient, req gitItemRequest, emit func(GitItem) bool) error {
	if req.Path != "" && req.ScopePath != "" {
		return errors.New("azuredevops_git_item: give either path (one file) or scope_path (a subtree listing), not both")
	}

	if req.Path != "" {
		params := gitItemParams(req)
		params.Set("path", req.Path)
		params.Set("includeContent", "true")
		params.Set("includeContentMetadata", "true")
		params.Set("recursionLevel", "none")
		item := new(GitItem)
		if err := client.getJSON(ctx, gitItemsURL(client.baseURL, req, params), scopeCode, item); err != nil {
			if isRestNotFound(err) {
				return nil
			}
			return err
		}
		// Rows echo the requested path, because Postgres re-applies the path qualifier to every row.
		item.Path = req.Path
		normalizeContent(item)
		emit(*item)
		return nil
	}

	scope := req.ScopePath
	if scope == "" {
		scope = "/"
	}
	params := gitItemParams(req)
	params.Set("scopePath", scope)
	params.Set("recursionLevel", "full")
	list := new(gitItemList)
	if err := client.getJSON(ctx, gitItemsURL(client.baseURL, req, params), scopeCode, list); err != nil {
		if isRestNotFound(err) {
			return nil
		}
		return err
	}
	for _, item := range list.Value {
		item.Content = nil
		if !emit(item) {
			return nil
		}
	}
	return nil
}

func gitItemParams(req gitItemRequest) url.Values {
	params := url.Values{}
	params.Set("api-version", gitItemsAPIVersion)
	if req.Version != "" {
		versionType := req.VersionType
		if versionType == "" {
			versionType = "branch"
			if fullCommitID.MatchString(req.Version) {
				versionType = "commit"
			}
		}
		params.Set("versionDescriptor.version", req.Version)
		params.Set("versionDescriptor.versionType", versionType)
	}
	return params
}

func gitItemsURL(baseURL string, req gitItemRequest, params url.Values) string {
	return strings.TrimRight(baseURL, "/") + "/" + url.PathEscape(req.Project) +
		"/_apis/git/repositories/" + url.PathEscape(req.RepositoryID) + "/items?" + params.Encode()
}

// normalizeContent keeps content, and its size in bytes, only for text files.
func normalizeContent(item *GitItem) {
	if item.IsFolder || item.Content == nil {
		item.Content = nil
		return
	}
	binary := item.ContentMetadata != nil && item.ContentMetadata.IsBinary
	if binary || !utf8.ValidString(*item.Content) || strings.IndexByte(*item.Content, 0) >= 0 {
		item.Content = nil
		return
	}
	size := int64(len(*item.Content))
	item.Size = &size
}
