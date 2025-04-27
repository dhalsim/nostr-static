// Time ago formatter
function formatTimeAgo(timestamp) {
    const now = Math.floor(Date.now() / 1000);
    const diff = now - timestamp;
    
    const intervals = {
        year: 31536000,
        month: 2592000,
        week: 604800,
        day: 86400,
        hour: 3600,
        minute: 60,
        second: 1
    };

    for (const [unit, seconds] of Object.entries(intervals)) {
        const interval = Math.floor(diff / seconds);
        if (interval >= 1) {
            return `${interval} ${unit}${interval === 1 ? '' : 's'} ago`;
        }
    }
    
    return 'just now';
}

// Initialize all time-ago elements on the page
function initTimeAgo() {
    const timeElements = document.querySelectorAll('.time-ago');
    timeElements.forEach(element => {
        const timestamp = parseInt(element.getAttribute('data-timestamp'));
        if (!isNaN(timestamp)) {
            element.textContent = formatTimeAgo(timestamp);
        }
    });
}

// Update all time-ago elements every minute
function startTimeAgoUpdates() {
    initTimeAgo();
    setInterval(initTimeAgo, 60000);
}

// Start the updates when the DOM is loaded
document.addEventListener('DOMContentLoaded', startTimeAgoUpdates); 