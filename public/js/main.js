document.addEventListener('DOMContentLoaded', function() {
    const tasksUrl = "/tasks";
    const activeTasksList = document.getElementById('activeTasks');
    const plannedTasksList = document.getElementById('plannedTasks');

    fetchTasks();

    async function fetchTasks() {
        try {
            const response = await fetch(tasksUrl);
            
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const data = await response.json();
            renderActiveTasks(data.active);
            renderPlannedTasks(data.planned);
        } catch (error) {
            console.error('Error fetching tasks:', error);
        }
    }

    function createTaskItem(taskId, duration, status) {
        const listItem = document.createElement("li");
        listItem.id = taskId;
        listItem.textContent = `${taskId}: ${duration} ms`;
        listItem.classList.add(status);
        return listItem;
    }

    function renderTasks(list, tasks, status) {
        list.innerHTML = '';

        if (Array.isArray(tasks)) {
            tasks.forEach(task => {
                const taskId = Object.keys(task)[0];
                const taskDuration = task[taskId];
                list.appendChild(createTaskItem(taskId, taskDuration, status));
            });
        } else {
            for (const taskId in tasks) {
                if (tasks.hasOwnProperty(taskId)) {
                    list.appendChild(createTaskItem(taskId, tasks[taskId], status));
                }
            }
        }
    }

    function renderActiveTasks(activeTasks) {
        renderTasks(activeTasksList, activeTasks, 'active');
    }

    function renderPlannedTasks(plannedTasks) {
        renderTasks(plannedTasksList, plannedTasks, 'planned');
    }

    function addToPlanned(task) {
        plannedTasksList.appendChild(createTaskItem(task.ID, task.Duration, 'planned'));
    }

    function nextPlanned(task) {
        const taskItem = document.getElementById(task.ID);
        if (taskItem) {
            // change class
            taskItem.classList.remove('planned');
            taskItem.classList.add('next');
        }
    }

    function moveToActive(task) {
        // Remove the "last-started" class from the last started task
        const lastStarted = document.querySelector('.last-started');
        if (lastStarted) {
            lastStarted.classList.remove('last-started');
            lastStarted.classList.add('active');
        }

        // Check if the task is in planned and move it to active
        const taskItem = document.getElementById(task.ID);
        if (taskItem) {
            // Remove it from planned
            plannedTasksList.removeChild(taskItem);
            // Add it to active
            taskItem.classList.remove('planned');
            taskItem.classList.add('last-started');
            activeTasksList.appendChild(taskItem);
        }
    }

    function removeFromActive(task) {
        // Find and remove the task from active tasks
        const taskItem = document.getElementById(task.ID);
        if (taskItem) {
            activeTasksList.removeChild(taskItem);
        }
    }

    window.addToPlanned = addToPlanned;
    window.nextPlanned = nextPlanned;
    window.moveToActive = moveToActive;
    window.removeFromActive = removeFromActive;
});