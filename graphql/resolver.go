package graphql

// This file will always be generated when running gqlgen.
// werrors imports in the resolvers are due to https://github.com/99designs/gqlgen/issues/1171.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/evergreen-ci/evergreen"
	"github.com/evergreen-ci/evergreen/rest/data"
	restModel "github.com/evergreen-ci/evergreen/rest/model"
	"github.com/evergreen-ci/gimlet"
	"github.com/evergreen-ci/utility"
)

const (
	CreateProjectMutation = "CreateProject"
	CopyProjectMutation   = "CopyProject"
	DeleteProjectMutation = "DeleteProject"
)

type Resolver struct {
	sc data.Connector
}

func New(apiURL string) Config {
	c := Config{
		Resolvers: &Resolver{
			sc: &data.DBConnector{URL: apiURL},
		},
	}
	c.Directives.RequireProjectAdmin = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		// Allow if user is superuser.
		user := mustHaveUser(ctx)
		opts := gimlet.PermissionOpts{
			Resource:      evergreen.SuperUserPermissionsID,
			ResourceType:  evergreen.SuperUserResourceType,
			Permission:    evergreen.PermissionProjectCreate,
			RequiredLevel: evergreen.ProjectCreate.Value,
		}
		if user.HasPermission(opts) {
			return next(ctx)
		}

		// Check for admin permissions for each of the resolvers.
		args, isStringMap := obj.(map[string]interface{})
		if !isStringMap {
			return nil, ResourceNotFound.Send(ctx, "Project not specified")
		}
		operationContext := graphql.GetOperationContext(ctx).OperationName

		if operationContext == CreateProjectMutation {
			canCreate, err := user.HasProjectCreatePermission()
			if err != nil {
				return nil, InternalServerError.Send(ctx, fmt.Sprintf("checking user permissions: %s", err.Error()))
			}
			if canCreate {
				return next(ctx)
			}
		}

		if operationContext == CopyProjectMutation {
			projectIdToCopy, ok := args["project"].(map[string]interface{})["projectIdToCopy"].(string)
			if !ok {
				return nil, InternalServerError.Send(ctx, "finding projectIdToCopy for copy project operation")
			}
			opts := gimlet.PermissionOpts{
				Resource:      projectIdToCopy,
				ResourceType:  evergreen.ProjectResourceType,
				Permission:    evergreen.PermissionProjectSettings,
				RequiredLevel: evergreen.ProjectSettingsEdit.Value,
			}
			if user.HasPermission(opts) {
				return next(ctx)
			}
		}

		if operationContext == DeleteProjectMutation {
			projectId, ok := args["projectId"].(string)
			if !ok {
				return nil, InternalServerError.Send(ctx, "finding projectId for delete project operation")
			}
			opts := gimlet.PermissionOpts{
				Resource:      projectId,
				ResourceType:  evergreen.ProjectResourceType,
				Permission:    evergreen.PermissionProjectSettings,
				RequiredLevel: evergreen.ProjectSettingsEdit.Value,
			}
			if user.HasPermission(opts) {
				return next(ctx)
			}
		}

		return nil, Forbidden.Send(ctx, fmt.Sprintf("user %s does not have permission to access the %s resolver", user.Username(), operationContext))
	}
	c.Directives.RequireProjectAccess = func(ctx context.Context, obj interface{}, next graphql.Resolver, access ProjectSettingsAccess) (res interface{}, err error) {
		user := mustHaveUser(ctx)

		var permissionLevel int
		if access == ProjectSettingsAccessEdit {
			permissionLevel = evergreen.ProjectSettingsEdit.Value
		} else if access == ProjectSettingsAccessView {
			permissionLevel = evergreen.ProjectSettingsView.Value
		} else {
			return nil, Forbidden.Send(ctx, "Permission not specified")
		}

		args, isStringMap := obj.(map[string]interface{})
		if !isStringMap {
			return nil, ResourceNotFound.Send(ctx, "Project not specified")
		}

		projectId, err := getProjectIdFromArgs(ctx, args)
		if err != nil {
			return nil, err
		}

		opts := gimlet.PermissionOpts{
			Resource:      projectId,
			ResourceType:  evergreen.ProjectResourceType,
			Permission:    evergreen.PermissionProjectSettings,
			RequiredLevel: permissionLevel,
		}
		if user.HasPermission(opts) {
			return next(ctx)
		}
		return nil, Forbidden.Send(ctx, fmt.Sprintf("user %s does not have permission to access settings for the project %s", user.Username(), projectId))
	}
	c.Directives.RequireProjectFieldAccess = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		user := mustHaveUser(ctx)

		projectRef, isProjectRef := obj.(*restModel.APIProjectRef)
		if !isProjectRef {
			return nil, InternalServerError.Send(ctx, "project not valid")
		}

		projectId := utility.FromStringPtr(projectRef.Id)
		if projectId == "" {
			return nil, ResourceNotFound.Send(ctx, "project not specified")
		}

		opts := gimlet.PermissionOpts{
			Resource:      projectId,
			ResourceType:  evergreen.ProjectResourceType,
			Permission:    evergreen.PermissionProjectSettings,
			RequiredLevel: evergreen.ProjectSettingsView.Value,
		}
		if user.HasPermission(opts) {
			return next(ctx)
		}
		return nil, Forbidden.Send(ctx, fmt.Sprintf("user does not have permission to access the field '%s' for project with ID '%s'", graphql.GetFieldContext(ctx).Path(), projectId))
	}
	return c
}
