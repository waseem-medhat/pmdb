function validateForm() {
    /**
     * @type {HTMLButtonElement}
     */
    const btnEl = document.getElementById("submit-btn")

    /**
     * @type {HTMLInputElement}
     */
    const passwordEl = document.getElementById("password")

    /**
     * @type {HTMLInputElement}
     */
    const confirmPasswordEl = document.getElementById("confirm-password")

    const p1 = passwordEl.value
    const p2 = confirmPasswordEl.value

    console.log(p1, p2)

    if (p1 !== p2 || p1.length < 8) {
        btnEl.disabled = true
    } else {
        btnEl.disabled = false
    }
    console.log("we good bruh")
}

document.getElementById("password").addEventListener("keyup", validateForm)
document.getElementById("confirm-password").addEventListener("keyup", validateForm)
