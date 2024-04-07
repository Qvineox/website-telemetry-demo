function likeLesson(id) {
    fetch('/api/lessons/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(
            {"id": parseInt(id)}
        )
    }).then((response) => {
        if (response.ok) {
            console.debug("lesson liked successfully")
            location.reload()
        } else {
            console.error(response.json())
        }
    })
}

function dislikeLesson(id) {
    fetch('/api/lessons/dislike', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(
            {"id": parseInt(id)}
        )
    }).then((response) => {
        if (response.ok) {
            console.debug("lesson disliked successfully")
            location.reload()
        } else {
            console.error(response.json())
        }
    })
}

function commentLesson(id) {
    let message = prompt("Напишите комментарий")
    if (message == null || message.length === 0) {
        return
    }

    fetch('/api/lessons/comment', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(
            {"id": parseInt(id), "comment": message}
        )
    }).then((response) => {
        if (response.ok) {
            console.debug("lesson commented successfully")
            location.reload()
        } else {
            console.error(response.json())
        }
    })
}