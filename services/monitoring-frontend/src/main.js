const backendUrl = import.meta.env.VITE_BACKEND_URL;
const socket = new WebSocket("ws://" + backendUrl + "/info/ws");

const containersMap = new Map();

socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    updateTable(data.containers);
};

function updateTable(containers) {
    const table = document.getElementById('container-table-body');
    table.innerHTML = ``;

    containers.forEach((container) => {
        container.updated_at = formatDate(container.updated_at);
        container.ports = formatPorts(container.ports);

        const row = getRow(container)
        table.appendChild(row);
    });
}

function getRow(container) {
    if (!containersMap.has(container.id) ) {
        containersMap.set(container.id, container);
        return getRowTemplate(container, 'table-success');
    }

    const containerPrevState = containersMap.get(container.id);
    if (containerPrevState.updated_at === container.updated_at) {
        containerPrevState.down = true
        return getRowTemplate(container, 'table-warning');
    }

    if (containerPrevState.down) {
        containerPrevState.down = false
        return getRowTemplate(container, 'table-primary');
    }

    containersMap.set(container.id, container);
    return getRowTemplate(container);
}

function getRowTemplate(container, rowClass) {
    const row = document.createElement('tr');
    
    if (rowClass) {
        row.classList.add(rowClass);
    }

    row.innerHTML = `
        <td>${container.id}</td>
        <td>${container.name}</td>
        <td>${container.image}</td>
        <td>${container.ip}</td>
        <td>${container.ports}</td>
        <td>${container.status}</td>
        <td>${container.updated_at}</td>
    `;
    return row;
}

function formatDate(dateString) {
    const options = { 
        hour: '2-digit', 
        minute: '2-digit', 
        second: '2-digit', 
        hour12: false 
    };
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU', options);
}

function formatPorts(ports) {
    return ports.join('<br>');
}
