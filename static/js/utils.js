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
