export function getErrorMessage(str) {
    try {
        const parsed = JSON.parse(str);
        return parsed.message;
    } catch (e) {
        return str;
    }

}