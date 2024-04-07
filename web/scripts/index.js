// button.addEventListener("click", (evt) => {
//     console.log(evt)
//
//     monitor.sendMonitoringEvent(evt.target, "click")
// })

// function onButtonClick(event) {
//
// }

class MonitoringDispatcher {
    sendMonitoringEvent(element, message, eventType) {
        console.debug(`logging event '${element} : ${eventType}': ${message}...`)

        fetch('/api/monitoring/event', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json;charset=utf-8'
            },
            body: JSON.stringify(
                {"element": element, message: message, "event_type": eventType}
            )
        }).then((response) => {
            if (response.ok) {
                console.debug("event saved successfully")
            } else {
                console.error(response.json())
            }
        })
    }

    constructor(userSessionID) {
        this.userSessionID = userSessionID
    }
}

let monitor = new MonitoringDispatcher()

async function login() {
    const username = document.querySelector("input#username-input").value;
    const password = document.querySelector("input#password-input").value;

    if (username.length === 0 || password.length === 0) {
        alert("Логин или пароль отсутствует!")
        return
    }

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
            monitor = new MonitoringDispatcher()
            sessionStorage.setItem("username", username)

            console.debug("login successful")
            window.location.href = "/";
        } else {
            console.debug("login error")
            sessionStorage.clear()

            alert("Ошибка авторизации!")
        }
    })
}

async function logout() {
    fetch('/api/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        }
    }).then((response) => {
        if (response.ok) {
            console.debug("logout successful")
            sessionStorage.clear()

            window.location.href = "/login";
        } else {
            console.debug("logout error")
        }
    })
}

