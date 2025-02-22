document.addEventListener("DOMContentLoaded", function () {
    const taskTableBody = document.getElementById("taskTableBody");
    const addTaskButton = document.getElementById("addTask");
    const taskForm = document.getElementById("taskForm");

    // Add a new row with inputs
    addTaskButton.addEventListener("click", function () {
        const newRow = document.createElement("tr");
        newRow.classList.add("task-row");

        newRow.innerHTML = `
            <td>
                <input type="text" name="taskID" placeholder="Enter task ID">
            </td>
            <td>
                <input type="number" name="taskDuration" placeholder="Enter task duration(ms)">
            </td>
            <td>
                <button type="button" class="removeTask">❌</button>
            </td>
        `;

        // Insert the new row before the last row with buttons
        taskTableBody.insertBefore(newRow, taskTableBody.lastElementChild);

        // Assign a click handler for removing the row
        newRow.querySelector(".removeTask").addEventListener("click", function () {
            newRow.remove();
        });

        addTabNavigation(newRow);
    });

    // Handle form submission
    taskForm.addEventListener("submit", function (event) {
        event.preventDefault();

        const taskData = {};

        document.querySelectorAll(".task-row").forEach(row => {
            const taskID = row.querySelector("input[name='taskID']").value.trim();
            const taskDuration = row.querySelector("input[name='taskDuration']").value.trim();

            if (taskID && taskDuration) {
                taskData[taskID] = Number(taskDuration);
            }
        });

        if (Object.keys(taskData).length === 0) {
            alert("Please add at least one task!");
            return;
        }

        fetch("/plan", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(taskData)
        })
        .then(response => response.json())
        // .then(data => console.log("Server response:", data))
        .catch(error => console.error("Error:", error));
    });

    // Assign remove handlers to existing remove buttons
    document.querySelectorAll(".removeTask").forEach(button => {
        button.addEventListener("click", function () {
            this.closest("tr").remove();
        });
    });

    function addTabNavigation(row) {
        const inputs = row.querySelectorAll("input");

        inputs.forEach((input, index) => {
            input.addEventListener("keydown", function (event) {
                if (event.key === "Tab") {
                    event.preventDefault();

                    const nextRow = row.nextElementSibling;
                    if (nextRow) {
                        const nextInput = nextRow.querySelectorAll("input")[index];
                        if (nextInput) {
                            nextInput.focus();
                        } else {
                            document.getElementById("addTask").click();
                            row.nextElementSibling.querySelector("input").focus();
                        }
                    }
                }
            });
        });
    }
});