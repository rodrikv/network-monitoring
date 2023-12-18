function stringToColorCode(input, darknessFactor = 0.7) {
    // Simple hash function (djb2 algorithm)
    let hash = 5381;

    for (let i = 0; i < input.length; i++) {
        hash = (hash * 33) ^ input.charCodeAt(i);
    }

    // Convert the hash value to RGB
    const red = (hash & 0xFF0000) >> 16;
    const green = (hash & 0x00FF00) >> 8;
    const blue = hash & 0x0000FF;

    // Adjust the RGB values to make the color darker
    const darkerRed = Math.floor(red * darknessFactor);
    const darkerGreen = Math.floor(green * darknessFactor);
    const darkerBlue = Math.floor(blue * darknessFactor);

    // Convert the darker RGB values back to hex
    const darkerColorCode = `#${(darkerRed << 16 | darkerGreen << 8 | darkerBlue).toString(16).padStart(6, '0')}`;

    return darkerColorCode;
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
const rand255 = () => Math.round(Math.random() * 255);

function randomColor(alpha) {
    return 'rgba(' + rand255() + ',' + rand255() + ',' + rand255() + ',' + (alpha || '.3') + ')';
}