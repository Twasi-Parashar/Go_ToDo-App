document.addEventListener('DOMContentLoaded', fetchTasks);

function fetchTasks() {
    fetch('/api/tasks')
        .then(res => res.json())
        .then(tasks => {
            const list = document.getElementById('task-list');
            list.innerHTML = '';
            tasks.forEach(task => {
                const li = document.createElement('li');
                li.innerHTML = `
                    <div>
                        <strong>${task.title}</strong> — ${task.description} <br>
                        <small>Due: ${task.due_date}</small>
                    </div>
                    <div>
                        <button onclick="deleteTask(${task.id})">❌</button>
                    </div>
                `;
                list.appendChild(li);
            });
        });
}

function addTask() {
    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
    const due_date = document.getElementById('due_date').value;

    fetch('/api/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description, due_date })
    }).then(() => {
        fetchTasks();
        document.getElementById('title').value = '';
        document.getElementById('description').value = '';
        document.getElementById('due_date').value = '';
    });
}

function deleteTask(id) {
    fetch(`/api/tasks/${id}`, {
        method: 'DELETE'
    }).then(fetchTasks);
}
