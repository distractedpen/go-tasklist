function handleLogin(e) {
    e.preventDefault();
    const body = getFormData();
    if (!body) {
        clearFormData();
        return;
    }

    fetch("/api/auth/login", {
        method: "POST",
        body: JSON.stringify(body),
    })
        .then((res) => {
            if (res.status == 200) {
                return res.json();
            } else {
                throw Error("Error Logging in!");
            }
        })
        .then((json) => {
            data = json.data;
            localStorage.setItem("userToken", data.token);
            localStorage.setItem("userId", data.id);
            localStorage.setItem("userName", data.username);
            location.pathname = "/tasks.html";
        })
        .catch((err) => {
            setErrorMessage(err);
        });
}

function handleRegister(e) {
    e.preventDefault();
    let body = getFormData();
    if (!body) {
        clearFormData();
        return;
    }
    fetch("/api/auth/register", {
        method: "POST",
        body: JSON.stringify(body),
    })
    .then((res) => {
        if (res.status == 200) {
            return res.json();
        } else {
            throw Error("Error Registering!");
        }
    })
    .then((_) => {
        alert("Register Successful! Please login.");
        location.reload()
    })
    .catch((err) => {
        setErrorMessage(err);
    });
}

function setErrorMessage(message) {
    const errorEle = document.getElementById("loginform-error");
    errorEle.innerText = message;
}

function getFormData() {
    const usernameEle = document.getElementById("loginform-username");
    const passwordEle = document.getElementById("loginform-password");
    const body = {
        username: usernameEle.value,
        password: passwordEle.value,
    };
    if (body.username.length == 0 || body.password.length == 0) {
        setErrorMessage("Username and Password required");
        return undefined;
    }
    return body;
}

function clearFormData() {
    const usernameEle = document.getElementById("loginform-username");
    const passwordEle = document.getElementById("loginform-password");
    usernameEle.value = "";
    passwordEle.value = "";
}


// page button setup
document.getElementById("loginform-buttons-login").onclick = (event) => handleLogin(event);
document.getElementById("loginform-buttons-register").onclick = (event) => handleRegister(event);
