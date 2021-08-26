/**
 * Copyright (c) 2021 Gitpod GmbH. All rights reserved.
 * Licensed under the GNU Affero General Public License (AGPL).
 * See License-AGPL.txt in the project root for license information.
 */

import { getGitpodService } from "./service/service";
import Cookies from "js-cookie";
import { v4 } from 'uuid';

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

//call this when the path changes. Complete page call is unnecessary for SPA after initial call
export const trackPathChange = (path: string) => {
    getGitpodService().server.trackEvent({
        event: "path_changed",
        properties: {
            path: path
        }
    });
}

//call this to record a page call
export const trackLocation = async () => {
    // retrieve anonymousId from Cookie. If not set yet, generate 'ajs_anonymous_id' cookie
    let anonymousId;
    const ajsCookie = Cookies.get('ajs_anonymous_id')
    if (ajsCookie) {
        anonymousId = ajsCookie;
    } else {
        anonymousId = v4()
        Cookies.set('ajs_anonymous_id', anonymousId);
    }

    // get public IPv4 address
    const publicIp = require('react-public-ip');
    const ip = await publicIp.v4();

    //get User Agent
    const { getUserAgent } = require("universal-user-agent");
    const userAgent = getUserAgent();

    getGitpodService().server.trackLocation({
        anonymousId: anonymousId,
        properties: {
            referrer: document.referrer,
            path: window.location.pathname,
            host: window.location.hostname,
            url: window.location.href
        },
        context: {
            userAgent: userAgent,
            ip: ip
        }
    })
}