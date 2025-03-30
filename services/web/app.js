document.addEventListener("DOMContentLoaded", function () {
    const stagesList = document.getElementById("stages-list");
    const workersList = document.getElementById("workers-list");
    const eventQueuesList = document.getElementById("event-queues-list");
    const performanceMetricsList = document.getElementById("performance-metrics-list");
    const notificationCenter = document.getElementById("notification-center");

    const eventSource = new EventSource("/events");

    eventSource.onmessage = function (event) {
        const data = JSON.parse(event.data);

        updateStages(data.stages);
        updateWorkers(data.workers);
        updateEventQueues(data.eventQueues);
        updatePerformanceMetrics(data.performanceMetrics);
        updateNotifications(data.notifications);
    };

    function updateStages(stages) {
        stagesList.innerHTML = "";
        stages.forEach(stage => {
            const li = document.createElement("li");
            li.textContent = `Stage: ${stage.name}, Status: ${stage.status}`;
            stagesList.appendChild(li);
        });
    }

    function updateWorkers(workers) {
        workersList.innerHTML = "";
        workers.forEach(worker => {
            const li = document.createElement("li");
            li.textContent = `Worker: ${worker.name}, Status: ${worker.status}`;
            workersList.appendChild(li);
        });
    }

    function updateEventQueues(eventQueues) {
        eventQueuesList.innerHTML = "";
        eventQueues.forEach(queue => {
            const li = document.createElement("li");
            li.textContent = `Queue: ${queue.name}, Length: ${queue.length}`;
            eventQueuesList.appendChild(li);
        });
    }

    function updatePerformanceMetrics(metrics) {
        performanceMetricsList.innerHTML = "";
        metrics.forEach(metric => {
            const li = document.createElement("li");
            li.textContent = `Metric: ${metric.name}, Value: ${metric.value}`;
            performanceMetricsList.appendChild(li);
        });
    }

    function updateNotifications(notifications) {
        notificationCenter.innerHTML = "";
        notifications.forEach(notification => {
            const div = document.createElement("div");
            div.textContent = `Notification: ${notification.message}`;
            notificationCenter.appendChild(div);
        });
    }
});
