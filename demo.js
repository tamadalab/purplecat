'use strict'

const heroku = "https://afternoon-wave-39227.herokuapp.com/purplecat/api/"

const availableUrl = (urlString) => {
    console.log(`availableUrl(${urlString})`)
    if (!urlString) {
        return false
    }
    const url = new URL(urlString)
    return url.hostname != "" && url.pathname.endsWith(".pom")
}

const availableFiles = (files) => {
    const flag = files.length > 0 && (files[0].name === "pom.xml" || files[0].name.endsWith(".pom"))
    console.log(`availableFiles(len=${files.length})`)
    return flag
}

const checkUrl = (e) => {
    const text = document.getElementById("pomurl").value
    const files = document.getElementById('pomfile').files

    const runButton = document.getElementById('runButton')
    runButton.disabled = !(availableUrl(text) || availableFiles(files))
}

const reset = (e) => {
    const file = document.getElementById("pomfile")
    file.value = ""

    const text = document.getElementById("pomurl")
    text.value = ""

    const area = document.getElementById("resultArea")
    area.innerText = ""

    const runButton = document.getElementById('runButton')
    runButton.disabled = true

    const message = document.getElementById('message')
    message.innerText = ""
}

const showMessage = (message) => {
    const messageArea = document.getElementById('message')
    messageArea.innerText = message
    messageArea.classList.remove("warning")
}

const showError = (message) => {
    const messageArea = document.getElementById('message')
    messageArea.innerText = message
    messageArea.classList.add("warning")
}

const showResult = (jsonString) => {
    const json = JSON.parse(jsonString)
    const str = JSON.stringify(json, null, "  ")
    const area = document.getElementById('resultArea')
    console.log(str)
    area.innerText = str
}

const createXmlHttpRequest = (doneMessage) => {
    const request = new XMLHttpRequest()
    request.onreadystatechange = () => {
        if (request.readyState == 4) {
            console.log(`done: http status: ${request.status}`)
            if (request.status == 200) {
                showMessage(doneMessage)
                showResult(request.responseText)
            } else {
                showError(request.responseText)
            }
        } else {
            showMessage("running purplecat...")
        }
    }
    return request
}

const getLicenses = (url) => {
    const request = createXmlHttpRequest(`GET license data from ${url}`)
    console.log(`${heroku}licenses?target=${url}`)
    request.open("GET", `${heroku}licenses?target=${url}`, true)
    request.send(null)
}

const postLicenses = (file) => {
    const request = createXmlHttpRequest(`POST license data`)
    request.open("POST", `${heroku}licenses`, true)
    request.setRequestHeader("Content-Type", "application/xml")
    const reader = new FileReader()
    reader.onload = (event) => {
        request.send(event.target.result)
    }
    reader.readAsText(file)
}

const executePurplecat = (e) => {
    const url = document.getElementById("pomurl").value
    const files = document.getElementById("pomfile").files
    if (url != "") {
        getLicenses(url)
    } else if (files.length >= 1) {
        postLicenses(files[0])
    }
}

const init = () => {
    const text = document.getElementById("pomurl")
    text.addEventListener('change', checkUrl)
    const file = document.getElementById('pomfile')
    file.addEventListener('change', checkUrl)

    const runButton = document.getElementById('runButton')
    runButton.disabled = true
    runButton.addEventListener('click', executePurplecat)

    const resetButton = document.getElementById('resetButton')
    resetButton.addEventListener('click', reset)

    console.log('init done')
}

window.onload = () => {
    init()
}
