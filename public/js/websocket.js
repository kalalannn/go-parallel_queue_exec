document.addEventListener('DOMContentLoaded', function() {
    const host = window.location.host;
    const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const wsUrl = `${wsProtocol}://${host}/ws`;
    const socket = new WebSocket(wsUrl);

    socket.onopen = function(event) {
        console.log('Connected to WebSocket server');
        socket.send('Hello from client');
    };

    socket.onerror = function(event) {
        console.log('WebSocket error:', event);
    };

    socket.onclose = function(event) {
        console.log('WebSocket connection closed');
    };

    socket.onmessage = function(event) {
        const message = JSON.parse(event.data);
        console.log('Received message from server:', message);

        if (message.hasOwnProperty('schedule')) {
            message.schedule.forEach(task => {
                window.addToPlanned(task);
            });
        } else if (message.hasOwnProperty('next')) {
            window.nextPlanned(message.next);
        } else if (message.hasOwnProperty('start')) {
            window.moveToActive(message.start);
        } else if (message.hasOwnProperty('done')) {
            window.removeFromActive(message.done);
        }
    };
});