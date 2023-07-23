let socket;
let spaceship = { id: "", x: 50, y: 50, angle: 0, isFiring: false };

function setup() {
    createCanvas(800, 600);
    socket = new WebSocket("ws://localhost:8080/ws");
    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleServerMessage(data);
    };
}

function draw() {
    background(220);
    drawSpaceship(spaceship);
}

function handleServerMessage(data) {
    // Update the spaceship data based on server messages
    // For example, update position, firing status, etc.
    // The actual implementation depends on your game logic.
    spaceship.x = data.x;
    spaceship.y = data.y;
    spaceship.angle = data.angle;
    spaceship.isFiring = data.isFiring;
}

function drawSpaceship(ship) {
    push();
    translate(ship.x, ship.y);
    rotate(ship.angle);
    // Draw the spaceship based on your design
    // For this example, we'll draw a simple triangle.
    fill(255);
    triangle(0, -20, -15, 20, 15, 20);
    pop();
}

function mouseMoved() {
    // Send mouse movement data to the server
    if (socket.readyState === WebSocket.OPEN) {
        const data = {
            x: mouseX,
            y: mouseY,
        };
        socket.send(JSON.stringify(data));
    }
}

function mousePressed() {
    // Send mouse click data to the server
    if (socket.readyState === WebSocket.OPEN) {
        const data = {
            isFiring: true,
        };
        socket.send(JSON.stringify(data));
    }
}

function mouseReleased() {
    // Send mouse release data to the server
    if (socket.readyState === WebSocket.OPEN) {
        const data = {
            isFiring: false,
        };
        socket.send(JSON.stringify(data));
    }
}
