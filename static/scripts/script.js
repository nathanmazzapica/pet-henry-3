// TODO: export to .env
const ws = new WebSocket(wsURL)

let daisyReferenceSize = {width: 894, height: 597};
let referenceNoseCoordinates = {x: 349, y: 145};

/** @type {coordinate} */
const noseCoordinates = {x: 349, y: 145};
const noseHonk = new Audio("../static/honk.mp3");
noseHonk.volume = 0.1;
let nextHonk = 0;


daisy.addEventListener("click", (e) => {

    const mousePos = {
        x: e.offsetX,
        y: e.offsetY
    }

    checkAndPerformEasterEggs(mousePos);
    petDaisy();
})


ws.onopen = () => {
    console.log("Connected to the server!");
};


ws.onmessage = (event) => {
    /**
     * @type event
     */
    const eventJSON = JSON.parse(event.data);
    console.log(eventJSON);

    const eventType = eventJSON.type;

    switch (eventType) {
        case "chat":
            chatMessageContainer.appendChild(buildMessage(eventJSON.data.sender, eventJSON.data.message));
            chatMessageContainer.scrollTop = chatMessageContainer.scrollHeight;
            break;
        case "pet":
            let prettyCount = Number(eventJSON.data.count).toLocaleString()
            counter.textContent = `Daisy has been pet ${prettyCount} times!`
            break;
        case "notification":
           displayToast("Notification", eventJSON.data.content, 2000)
            break;
        case "playerCountUpdate":
            document.getElementById("player-count").innerText = `Online Players: ${eventJSON.data.count}`;
            break;
        case "leaderboardUpdate":
            alert("leaderboard not implemented")
            break;
        default:
            alert("unknown event type received!")
            return;
    }


    if (eventJSON.name === "leaderboard") {
        console.log("handle leaderboard notification")
        console.log(JSON.parse(eventJSON.message));
        displayLeaderboard(JSON.parse(eventJSON.message));
        return;
    }

}

ws.onclose = () => {
    window.location.href = "/error"
}

function petDaisy() {
    personalNumber++;
    petMessage = {
        type: "pet",
        data: {}
    }
    ws.send(JSON.stringify(petMessage));
    prettyNumber = personalNumber.toLocaleString();
    personalCounter.innerText = `You have pet her ${prettyNumber} time${personalNumber === 1 ? "" : "s"}!`;
}

/**
 * @param {coordinate} mousePos
 */

function checkAndPerformEasterEggs(mousePos) {
    console.log(mousePos);
    if (inRadius(mousePos, noseCoordinates, 20) && nextHonk <= Date.now()) {
        console.log("honk");
        noseHonk.currentTime = 0;
        noseHonk.play();
        nextHonk = Date.now() + 500;
    }
}


function inRadius(point1, point2, radius) {
    const dist = Math.sqrt(Math.pow((point2.x - point1.x), 2) + Math.pow((point2.y - point1.y), 2));

    return dist <= radius;
}


window.addEventListener("resize", () => {
    const SCALE_OFFSET = 0.87
    if (window.innerWidth > 0) {
        setGradientPosition();
    } else {
        document.body.style.background = "var(--daisy-gradient-end)";
    }

    let daisySize = {width: daisy.width * SCALE_OFFSET, height: daisy.height * SCALE_OFFSET};
    console.log(daisySize)

    let ratio = {
        width: daisySize.width / daisyReferenceSize.width,
        height: daisySize.height / daisyReferenceSize.height
    };

    console.log(ratio)

    noseCoordinates.x = referenceNoseCoordinates.x * ratio.width;
    noseCoordinates.y = referenceNoseCoordinates.y * ratio.height;

    console.log(noseCoordinates.x)
    console.log(noseCoordinates.y)
});

window.addEventListener("load", () => {
        const SCALE_OFFSET = 0.87

        if (window.innerWidth > 0) {
            setGradientPosition();
        } else {
            document.body.style.background = "var(--daisy-gradient-end)";
        }

        let daisySize = {width: daisy.width * SCALE_OFFSET, height: daisy.height * SCALE_OFFSET};

        let ratio = {
            width: daisySize.width / daisyReferenceSize.width,
            height: daisySize.height / daisyReferenceSize.height
        };

        console.log(ratio)

        referenceNoseCoordinates.x = referenceNoseCoordinates.x * ratio.width;
        referenceNoseCoordinates.y = referenceNoseCoordinates.y * ratio.height;
        noseCoordinates.x = referenceNoseCoordinates.x;
        noseCoordinates.y = referenceNoseCoordinates.y;

        console.log(noseCoordinates.x)
        console.log(noseCoordinates.y)

        daisyReferenceSize = daisySize;
    }
);