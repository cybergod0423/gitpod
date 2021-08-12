/**
 * Copyright (c) 2021 Gitpod GmbH. All rights reserved.
 * Licensed under the GNU Affero General Public License (AGPL).
 * See License-AGPL.txt in the project root for license information.
 */

import { getGitpodService } from "./service/service";

declare global {
    interface Window { analytics: any; }
}

//contexts from which calls are made in dashboard
export type dashboard_contexts = "menu" | "/<team_name>/<project_name>/configure" | "/new" | "/<team_name>/<project_name>/prebuilds" | "/<team_name>/<project_name>" | "/projects" | "/<team_name>/members" | "/teams/new" | "/workspaces" | "/<team_name>/projects";
//buttons that are tracked in dashboard
export type buttons = "new_team" | "test_configuration" | "add_organisation" | "select_git_provider" | "select_project" | "select_team" | "continue_with_github" | "continue_with_gitlab" | "create_team" | "trigger_prebuild" | "new_workspace" | "rerun_prebuild" | "new_project" | "invite_members" | "remove_project" | "leave_team";
//position of tracked button in page
export type button_contexts = "dropdown" | "primary_button" | "secondary_button" | "kebab_menu" | "card";
//events are than generic button clicks that are tracked in dashboard
export type events = "invite_url_requested" | "workspace_new_clicked" | "workspace_button_clicked" | "organisation_authorised";
//actions that can be performed on workspaces in dashboard
export type workspace_actions = "open" | "stop" | "download" | "share" | "pin" | "delete";

//call this when a button in the dashboard is clicked
export const trackButton = (dashboard_context: dashboard_contexts, button: buttons, button_context: button_contexts) => {
    getGitpodService().server.trackEvent({
        event: "dashboard_button_clicked",
        properties: {
            dashboard_context: dashboard_context,
            button: button,
            button_context: button_context
        }
    })
}

//call this when a button that performs a certain action on a workspace is clicked
export const trackWorkspaceButton = (workspaceId: string, workspace_action: workspace_actions, button_context: button_contexts, state: string) => {
    getGitpodService().server.trackEvent({
        event: "workspace_button_clicked",
        properties: {
            workspaceId: workspaceId,
            workspace_action: workspace_action,
            button_context: button_context,
            state: state
        }
    })
}

//call this when anything that is not a button or a page call should be tracked
export const trackEvent = (event: events, properties: any) => {
    getGitpodService().server.trackEvent({
        event: event,
        properties: properties
    })
}

//call this to record a page call
export const page = () => {
    window.analytics.page();
}