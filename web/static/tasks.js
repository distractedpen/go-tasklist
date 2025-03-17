const API = "http://localhost:8080/api/";

if (!localStorage.getItem("userToken")) {
    location.pathname = "/index.html";
} else {
    // const userToken = localStorage.getItem("userToken");
    // const userId = localStorage.getItem("userToken");
    // TODO: validate current user token
}

// get the user task list data
// @param {id} user id
// @returns User
async function fetchUser(id) {
    await fetch(API + `tasks/${id}`,
        {
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("userToken"),
            },
        })
        .then((res) => res.json())
        .then((data) => {
            hydrateTaskList(data.data);
            setHeaderMessage();
        })
        .catch((err) => console.error(err));
}

function setHeaderMessage() {
    const header = document.getElementById("header-text");
    header.innerText = "Welcome " + localStorage.getItem("userName") + "!";
}

async function deleteTask(taskId) {
    const userId = localStorage.getItem("userId");
    await fetch(API + `tasks/${userId}/${taskId}`,
        {
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("userToken"),
            },
            method: "DELETE"
        })
        .catch((err) => console.error(err));
}

async function addTask(newTask) {
    const body = { name: newTask.get("name"), description: newTask.get("description") };
    const userid = localStorage.getItem("userId");
    await fetch(API + `tasks/${userid}`,
        {
            headers: {
                "Authorization": "Bearer " + localStorage.getItem("userToken"),
            },
            method: "POST",
            body: JSON.stringify(body)
        })
        .then((res) => res.json())
        .then((_) => { insertTask(body) })
        .catch((err) => console.error(err));
}

function openTaskForm() {
    const tasklist = document.getElementById("tasklist");
    const newTaskForm = document.createElement("form");
    const addButton = document.createElement("button");
    addButton.innerText = "Add Task";
    addButton.onclick = function(e) {
        e.preventDefault();
        const formData = new FormData(newTaskForm);
        addTask(formData);
        tasklist.removeChild(newTaskForm);
    }
    const cancelButton = document.createElement("button");
    cancelButton.innerText = "Cancel";
    cancelButton.onclick = function() {
        tasklist.removeChild(newTaskForm);
    }
    newTaskForm.appendChild(createInputField("name", "Name"));
    newTaskForm.appendChild(createInputField("description", "Description"));
    newTaskForm.appendChild(addButton);
    newTaskForm.appendChild(cancelButton);
    tasklist.insertBefore(newTaskForm, document.getElementById("addTaskButton"));
}

function createInputField(name, placeholder, hidden = false) {
    const input = document.createElement("input");
    input.name = name;
    input.placeholder = placeholder;
    input.hidden = hidden;
    return input;
}

function createTaskCard(task) {
    const tasklist = document.getElementById("tasklist");
    let taskElement = document.createElement("div");
    let taskNameEle = document.createElement("h2");
    let taskDescEle = document.createElement("p");
    let deleteButton = document.createElement("button");
    taskElement.className = "task-card";
    taskNameEle.innerText = task.name;
    taskDescEle.innerText = task.description;
    deleteButton.innerText = "Delete";
    deleteButton.onclick = function() {
        deleteTask(task.id);
        tasklist.removeChild(taskElement);
    }
    taskElement.appendChild(taskNameEle);
    taskElement.appendChild(taskDescEle);
    taskElement.appendChild(deleteButton);
    return taskElement;
}

// populate the user task list container
function hydrateTaskList(tasks) {
    const tasklist = document.getElementById("tasklist");
    const addButton = document.getElementById("addTaskButton");
    for (let task of tasks) {
        tasklist.insertBefore(createTaskCard(task), addButton);
    }
}

function insertTask(task) {
    const tasklist = document.getElementById("tasklist");
    const addButton = document.getElementById("addTaskButton");
    tasklist.insertBefore(createTaskCard(task), addButton);
}

fetchUser(localStorage.getItem("userId"));


function handleLogout() {
    // TODO: handle user logout and redirect to login page
    localStorage.removeItem("userId")
    localStorage.removeItem("userName")
    localStorage.removeItem("userToken")
    location.pathname = "/index.html";
}
