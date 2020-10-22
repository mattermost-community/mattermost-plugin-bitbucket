import React, {PureComponent} from 'react';
import PropTypes from 'prop-types';

import BitbucketIcon from '../../icon';

export default class CreateIssuePostMenuAction extends PureComponent {
    static propTypes = {
        show: PropTypes.bool.isRequired,
        open: PropTypes.func.isRequired,
        postId: PropTypes.string,
    };

    handleClick = (e) => {
        const {open, postId} = this.props;
        e.preventDefault();
        open(postId);
    };

    render() {
        if (!this.props.show) {
            return null;
        }

        const content = (
            <button
                className='style--none'
                role='presentation'
                onClick={this.handleClick}
            >
                <BitbucketIcon/>
                {'Create Bitbucket Issue'}
            </button>
        );

        return (
            <li
                className='MenuItem'
                role='menuitem'
            >
                {content}
            </li>
        );
    }
}
