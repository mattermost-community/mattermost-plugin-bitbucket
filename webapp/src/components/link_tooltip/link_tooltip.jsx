import React, {useEffect, useState} from 'react';
import PropTypes from 'prop-types';
import './tooltip.css';
import Octicon, {GitMerge, GitPullRequest, IssueClosed, IssueOpened} from '@primer/octicons-react';
import ReactMarkdown from 'react-markdown';

import Client from 'client';
import {hexToRGB} from '../../utils/styles';

export const LinkTooltip = ({href, connected, theme}) => {
    const [data, setData] = useState(null);
    useEffect(() => {
        const init = async () => {
            if (href.includes('bitbucket.org/')) {
                const [owner, repo, type, id] = href.split('bitbucket.org/')[1].split('/');
                let res;
                switch (type) {
                case 'issues':
                    res = await Client.getIssue(owner, repo, id);
                    break;
                case 'pull-requests':
                    res = await Client.getPullRequest(owner, repo, id);
                    break;
                }
                if (res) {
                    res.owner = owner;
                    res.repo = repo;
                    res.type = type;
                }
                setData(res);
            }
        };
        if (data) {
            return;
        }
        if (connected) {
            init();
        }
    }, []);

    const getIconElement = () => {
        let icon;
        let color;
        let iconType;
        switch (data.type) {
        case 'pull-requests':
            color = '#28a745';
            iconType = GitPullRequest;
            if (data.state === 'MERGED') {
                color = '#6f42c1';
                iconType = GitMerge;
            } else if (data.state === 'DECLINED') {
                color = '#cb2431';
            }
            icon = (
                <span style={{color}}>
                    <Octicon
                        icon={iconType}
                        size='small'
                        verticalAlign='middle'
                    />
                </span>
            );
            break;
        case 'issues':
            color = data.state === 'closed' ? '#cb2431' : '#28a745';
            iconType = data.state === 'closed' ? IssueClosed : IssueOpened;
            icon = (
                <span style={{color}}>
                    <Octicon
                        icon={iconType}
                        size='small'
                        verticalAlign='middle'
                    />
                </span>
            );
            break;
        }
        return icon;
    };

    if (data) {
        let date = new Date(data.created_on);
        date = date.toDateString();
        return (
            <div className='bitbucket-tooltip'>
                <div
                    className='bitbucket-tooltip box bitbucket-tooltip--large bitbucket-tooltip--bottom-left p-4'
                    style={{backgroundColor: theme.centerChannelBg, border: `1px solid ${hexToRGB(theme.centerChannelColor, '0.16')}`}}
                >
                    <div className='header mb-1'>
                        <span style={{color: theme.centerChannelColor}}>
                            {data.type === 'pull-requests' ? data.destination.repository.full_name : data.repository.full_name}
                        </span>
                        {' on '}
                        <span>{date}</span>
                    </div>

                    <div className='body d-flex mt-2'>
                        <span className='pt-1 pb-1 pr-2'>
                            { getIconElement() }
                        </span>

                        {/* info */}
                        <div className='tooltip-info mt-1'>
                            <a
                                href={href}
                                target='_blank'
                                rel='noopener noreferrer'
                                style={{color: theme.centerChannelColor}}
                            >
                                <h5 className='mr-1'>{data.title}</h5>
                                <span>{'#' + data.id}</span>
                            </a>
                            <div className='markdown-text mt-1 mb-1'>
                                <ReactMarkdown
                                    source={data.type === 'pull-requests' ? data.summary.raw : data.content.raw}
                                    linkTarget='_blank'
                                />
                            </div>

                            {/* base <- head */}
                            {data.type === 'pull-requests' && (
                                <div className='base-head mt-1 mr-3'>
                                    <span
                                        title={data.destination.repository.name + '/' + data.destination.branch.name}
                                        className='commit-ref'
                                    >{data.destination.branch.name}
                                    </span>
                                    <span className='mx-1'>{'←'}</span>
                                    <span
                                        title={data.source.repository.name + '/' + data.source.branch.name}
                                        className='commit-ref'
                                    >{data.source.branch.name}
                                    </span>
                                </div>
                            )}

                            <div className='see-more mt-1'>
                                <a
                                    href={href}
                                    target='_blank'
                                    rel='noopener noreferrer'
                                >{'See more'}</a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
    return null;
};

LinkTooltip.propTypes = {
    href: PropTypes.string.isRequired,
    connected: PropTypes.bool.isRequired,
    theme: PropTypes.object.isRequired,
};
