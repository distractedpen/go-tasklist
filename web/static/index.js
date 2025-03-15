
// User = {"username": str}

const API = "http://localhost:8080/api/";


// get the user task list data
// @param {id} user id
// @returns User
async function fetchUser(id) {
    await fetch(API + `tasks/${id}`, {})
        .then((res) => res.json())
        .then((data) => {
            sessionStorage.setItem("userId", id);
            hydrateTaskList(data.data);
        })
        .catch((err) => console.error(err));

}

async function deleteTask(taskId) {
    const userId = sessionStorage.getItem("userId");
    await fetch(API + `tasks/${userId}/${taskId}`, { method: "DELETE" })
        .catch((err) => console.error(err));
}

async function addTask(newTask) {
    console.log(newTask);
    const body = { name: newTask.get("name"), description: newTask.get("description") };
    const userid = sessionStorage.getItem("userId");
    await fetch(API + `tasks/${userid}`, { method: "POST", body: JSON.stringify(body) })
        .then((res) => res.json())
        .then((_) => { insertTask(body) })
        .catch((err) => console.error(err));
}

function openTaskForm() {
    console.log("new task form");
    const tasklist = document.getElementById("tasklist");
    const newTaskForm = document.createElement("form");
    const addButton = document.createElement("button");
    addButton.innerText = "Add Task";
    addButton.onclick = function(e) {
        e.preventDefault();
        const formData = new FormData(newTaskForm);
        console.log(formData);
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
    console.log("where is it?");
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

fetchUser(1);
