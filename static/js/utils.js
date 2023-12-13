function stringToColorCode(input) {
    // Simple hash function (djb2 algorithm)
    let hash = 5381;

    for (let i = 0; i < input.length; i++) {
        hash = (hash * 33) ^ input.charCodeAt(i);
    }

    // Convert the hash value to a hexadecimal color code
    const colorCode = '#' + (hash >>> 0).toString(16).slice(-6);

    return colorCode;
}

function formatDatetime(timestamp) {
    // Create a new Date object using the timestamp
    let datetime = new Date(timestamp);

    // Define options for formatting
    let options = {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        timeZoneName: 'short',
    };

    // Format the date using toLocaleString with the specified options
    let formattedDate = datetime.toLocaleString('en-GB', options);

    return formattedDate
}