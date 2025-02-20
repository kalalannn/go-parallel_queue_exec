console.log("Hello");

const host = window.location.host;
const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
const wsUrl = `${wsProtocol}://${host}/ws`;
const socket = new WebSocket(wsUrl);

socket.onopen = function(event) {
    console.log('Connected to WebSocket server');
    socket.send('Hello from client');
};

socket.onmessage = function(event) {
    console.log('Received message from server:', event.data);
};

socket.onerror = function(event) {
    console.log('WebSocket error:', event);
};

socket.onclose = function(event) {
    console.log('WebSocket connection closed');
};

function sendMessageToServer(message) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(message);
    } else {
        console.log('WebSocket connection is not open');
    }
}

sendMessageToServer('Hello2');
