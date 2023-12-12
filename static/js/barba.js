function fadeIn() {
    document.body.style.opacity = 1;
}

// Initial fade-in when the window is fully loaded
window.addEventListener('load', function () {
    fadeIn();
});