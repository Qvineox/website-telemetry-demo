const button = document.querySelector("button#auth-button");

button.addEventListener("click", (evt) => {
    console.log(evt)

    sendMonitoringEvent(evt.target, "click")
})

// function onButtonClick(event) {
//
// }

function sendMonitoringEvent(element, eventType) {
    console.debug(`sending ${eventType} event...`)
}

async function auth() {
    const username = document.querySelector("input#username-input").value;
    const password = document.querySelector("input#password-input").value;

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(
            {"username": username, "password": password}
        )
    }).then((response) => {
        if (response.ok) {
            console.debug("login successful")
            window.location.href = "/";
        } else {
            console.debug("login error")
            alert("Ошибка авторизации!")
        }
    })
}