import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {isSystemMessage} from 'mattermost-redux/utils/post_utils';

import AttachCommentToIssueModal from 'components/modals/attach_comment_to_issue';

import CreateIssueModal from './components/modals/create_issue';

import SidebarHeader from './components/sidebar_header';
import TeamSidebar from './components/team_sidebar';
import UserAttribute from './components/user_attribute';
import SidebarRight from './components/sidebar_right';
import LinkTooltip from './components/link_tooltip';
import Reducer from './reducers';
import manifest from './manifest';
import {getConnected, openAttachCommentToIssueModal, openCreateIssueModal, setShowRHSAction} from './actions';
import {handleConnect, handleDisconnect, handleReconnect, handleRefresh} from './websocket';

let activityFunc;
let lastActivityTime = Number.MAX_SAFE_INTEGER;
const activityTimeout = 60 * 60 * 1000; // 1 hour

class PluginClass {
    async initialize(registry, store) {
        registry.registerReducer(Reducer);

        await getConnected(true)(store.dispatch, store.getState);

        registry.registerLeftSidebarHeaderComponent(SidebarHeader);
        registry.registerBottomTeamSidebarComponent(TeamSidebar);
        registry.registerPopoverUserAttributesComponent(UserAttribute);
        registry.registerRootComponent(CreateIssueModal);
        const showPostMenuAction = (postId) => {
            const state = store.getState();
            const post = getPost(state, postId);
            return Boolean(state[`plugins-${manifest.id}`].connected && post && !isSystemMessage(post));
        };
        registry.registerPostDropdownMenuAction({
            text: 'Create Bitbucket Issue',
            action: (postId) => store.dispatch(openCreateIssueModal(postId)),
            filter: showPostMenuAction,
        });
        registry.registerRootComponent(AttachCommentToIssueModal);
        registry.registerPostDropdownMenuAction({
            text: 'Attach to Bitbucket Issue',
            action: (postId) => store.dispatch(openAttachCommentToIssueModal(postId)),
            filter: showPostMenuAction,
        });
        registry.registerLinkTooltipComponent(LinkTooltip);

        const {showRHSPlugin} = registry.registerRightHandSidebarComponent(SidebarRight, 'Bitbucket');
        store.dispatch(setShowRHSAction(() => store.dispatch(showRHSPlugin)));

        registry.registerWebSocketEventHandler('custom_bitbucket_connect', handleConnect(store));
        registry.registerWebSocketEventHandler('custom_bitbucket_disconnect', handleDisconnect(store));
        registry.registerWebSocketEventHandler('custom_bitbucket_refresh', handleRefresh(store));
        registry.registerReconnectHandler(handleReconnect(store));

        activityFunc = () => {
            const now = new Date().getTime();
            if (now - lastActivityTime > activityTimeout) {
                handleReconnect(store, true)();
            }
            lastActivityTime = now;
        };

        document.addEventListener('click', activityFunc);
    }

    deinitialize() {
        document.removeEventListener('click', activityFunc);
    }
}

global.window.registerPlugin('bitbucket', new PluginClass());
