import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

import {formatTimeSince} from 'utils/date_utils';

function BitbucketItems(props) {
    const style = getStyle(props.theme);

    return props.items.length > 0 ? props.items.map((item) => {
        const repoName = item.destination ? item.destination.repository.full_name : item.repository.full_name;

        let title = item.title;
        let number = null;

        if (item.number) {
            number = (
                <strong>
                    <i className='fa fa-code-fork'/>{' #' + item.number}
                </strong>);
        }

        if (item.links) {
            title = (
                <a
                    href={item.links.html.href}
                    target='_blank'
                    rel='noopener noreferrer'
                    style={style.itemTitle}
                >
                    {item.title}
                </a>);
            if (item.number) {
                number = (
                    <strong>
                        <a
                            href={item.links.html.href}
                            target='_blank'
                            rel='noopener noreferrer'
                        >
                            <i className='fa fa-code-fork'/>{' #' + item.number}
                        </a>
                    </strong>);
            }
        }

        let reviews = '';

        reviews = getReviewText(item, style, item.created_on);

        return (
            <div
                key={item.id}
                style={style.container}
            >
                <div>
                    <strong>
                        {title}
                    </strong>
                </div>
                <div>
                    {number} <span className='light'>{'(' + repoName + ')'}</span>
                </div>
                <div
                    className='light'
                    style={style.subtitle}
                >
                    {item.created_on && ('Opened ' + formatTimeSince(item.created_on) + ' ago')}
                    {item.created_on && '.'}
                </div>
                {reviews}
            </div>
        );
    }) : <div style={style.container}>{'You have no active items'}</div>;
}

BitbucketItems.propTypes = {
    items: PropTypes.array.isRequired,
    theme: PropTypes.object.isRequired,
};

const getStyle = makeStyleFromTheme((theme) => {
    return {
        container: {
            padding: '15px',
            borderTop: `1px solid ${changeOpacity(theme.centerChannelColor, 0.2)}`,
        },
        itemTitle: {
            color: theme.centerChannelColor,
            lineHeight: 1.7,
            fontWeight: 'bold',
        },
        subtitle: {
            margin: '5px 0 0 0',
            fontSize: '13px',
        },
        subtitleSecondLine: {
            fontSize: '13px',
        },
        icon: {
            top: 3,
            position: 'relative',
            left: 3,
            height: 18,
            display: 'inline-flex',
            alignItems: 'center',
        },
    };
});

function getReviewText(item, style, secondLine) {
    if (!item.participants) {
        return '';
    }

    const reviewersUsers = item.participants.filter((participant) => participant.role === 'REVIEWER');
    const reviewersNo = reviewersUsers.length;

    if (reviewersNo === 0) {
        return '';
    }

    let reviews = '';
    const approvedNo = reviewersUsers.filter((reviewer) => reviewer.approved).length;
    if (reviewersNo > 0) {
        let reviewName;
        if (reviewersNo === 1) {
            reviewName = 'review';
        } else {
            reviewName = 'reviews';
        }
        reviews = (<span>{approvedNo + ' out of ' + reviewersNo + ' ' + reviewName + ' complete.'}</span>);
    }

    return (
        <div
            className='light'
            style={secondLine ? style.subtitleSecondLine : style.subtitle}
        >
            {reviews}
        </div>);
}

export default BitbucketItems;
