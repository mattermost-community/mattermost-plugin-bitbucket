export function formatTimeSince(date) {
    const secondsSince = Math.trunc((Date.now() - (new Date(date)).getTime()) / 1000);
    if (secondsSince < 60) {
        if (secondsSince === 1) {
            return secondsSince + ' second';
        }
        return secondsSince + ' seconds';
    }
    const minutesSince = Math.trunc(secondsSince / 60);
    if (minutesSince < 60) {
        if (minutesSince === 1) {
            return minutesSince + ' minute';
        }
        return minutesSince + ' minutes';
    }
    const hoursSince = Math.trunc(minutesSince / 60);
    if (hoursSince < 24) {
        if (hoursSince === 1) {
            return hoursSince + ' hour';
        }
        return hoursSince + ' hours';
    }
    const daysSince = Math.trunc(hoursSince / 24);
    if (daysSince === 1) {
        return daysSince + ' day';
    }
    return daysSince + ' days';
}
