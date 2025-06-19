function toggleTheme() {
    document.body.classList.toggle('light-theme');
}

async function createPlan() {
    event.preventDefault();
    const weekNumber = document.getElementById('weekNumber').value;
    const weekDay = document.getElementById('weekDay').value;
    const time = document.getElementById('time').value;
    const group = document.getElementById('group').value;
    const teacher = document.getElementById('teacher').value;
    const subject = document.getElementById('subject').value;
    const subject_lvl2 = document.getElementById('subject_lvl2').value;
    const class_room = document.getElementById('class_room').value;

    if (!weekNumber || !weekDay || !time || !group || !teacher || !subject || !subject_lvl2 || !class_room) {
        alert('Пожалуйста, заполните все поля.');
        return;
    }

    const response = await fetch('home/create_item', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            id,
            weekNumber,
            weekDay,
            time,
            group,
            teacher,
            subject,
            subject_lvl2,
            class_room
        }),
    });

    const tableBody = document.getElementById('planTable').tBodies[0];
    const newRow = tableBody.insertRow();

    newRow.insertCell(0).textContent = id;
    newRow.insertCell(1).textContent = weekNumber;
    newRow.insertCell(2).textContent = weekDay;
    newRow.insertCell(3).textContent = time;
    newRow.insertCell(4).textContent = group;
    newRow.insertCell(5).textContent = teacher;
    newRow.insertCell(6).textContent = subject;
    newRow.insertCell(7).textContent = subject_lvl2;
    newRow.insertCell(8).textContent = class_room;

    // Очистка полей после создания
    document.getElementById('weekNumber').value = '';
    document.getElementById('weekDay').value = '';
    document.getElementById('time').value = '';
    document.getElementById('group').value = '';
    document.getElementById('teacher').value = '';
    document.getElementById('subject').value = '';
    document.getElementById('subject_lvl2').value = '';
    document.getElementById('class_room').value = '';
}

function savePlan() {
    const tableBody = document.getElementById('planTable').tBodies[0];
    const rows = tableBody.rows;
    const plans = [];

    for (let i = 0; i < rows.length; i++) {
        const cells = rows[i].cells;
        console.log('weekNumber:', cells[1].textContent);
        console.log('weekDay:', cells[2].textContent);
        console.log('time:', cells[3].textContent);
        console.log('group:', cells[4].textContent);
        console.log('teacher:', cells[5].textContent);
        console.log('subject:', cells[6].textContent);
        console.log('subject_lvl2:', cells[7].textContent);
        console.log('class_room:', cells[8].textContent);

        const plan = {
            id: cells[0].textContent,
            weekNumber: cells[1].textContent,
            weekDay: cells[2].textContent,
            time: cells[3].textContent,
            group: cells[4].textContent,
            teacher: cells[5].textContent,
            subject: cells[6].textContent,
            subject_lvl2: cells[7].textContent,
            class_room: cells[8].textContent
        };
        plans.push(plan);
    }

    localStorage.setItem('plans', JSON.stringify(plans));
    alert('Учебный план сохранён');
}

async function showPlan() {
    try {
        const response = await fetch('/home/getPlan', {
            method: 'GET'
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log(data); // Выводим полученные данные в консоль

        // Проверяем, есть ли данные для отображения
        if (data.length === 0) {
            alert('Учебный план пока пуст.');
        } else {
            displayData(data);
        }
    } catch (error) {
        console.error('Ошибка при получении данных:', error);
        alert('Произошла ошибка при загрузке учебного плана.');
    }
}

function displayData(data) {
    const tableBody = document.getElementById('planTable').tBodies[0];
    tableBody.innerHTML = ''; // Очищаем таблицу перед добавлением новых данных

    data.forEach(plan => {
        const newRow = tableBody.insertRow();
        newRow.insertCell(0).textContent = plan.id;
        newRow.insertCell(1).textContent = plan.weekNumber;
        newRow.insertCell(2).textContent = plan.weekDay;
        newRow.insertCell(3).textContent = plan.time;
        newRow.insertCell(4).textContent = plan.group;
        newRow.insertCell(5).textContent = plan.teacher;
        newRow.insertCell(6).textContent = plan.subject;
        newRow.insertCell(7).textContent = plan.subject_lvl2;
        newRow.insertCell(8).textContent = plan.class_room;
    });
}



async function updatePlan() {
    try {
        const id = document.getElementById('id').value;
        const weekNumber = document.getElementById('weekNumber').value;
        const weekDay = document.getElementById('weekDay').value;
        const time = document.getElementById('time').value;
        const group = document.getElementById('group').value;
        const teacher = document.getElementById('teacher').value;
        const subject = document.getElementById('subject').value;
        const subject_lvl2 = document.getElementById('subject_lvl2').value;
        const class_room = document.getElementById('class_room').value;

        const plan = {
            id,
            weekNumber,
            weekDay,
            time,
            group,
            teacher,
            subject,
            subject_lvl2,
            class_room
        };

        const response = await fetch('/home/updatePlan', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(plan)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log(data); // Выводим ответ сервера в консоль

        alert('План успешно обновлен');
    } catch (error) {
        console.error('Ошибка при обновлении данных:', error);
        alert('Произошла ошибка при обновлении плана.');
    }
}

function selectRowForUpdate(event) {
    const row = event.target.closest('tr');
    if (!row) return;

    const weekNumber = row.cells[1].textContent;
    const weekDay = row.cells[2].textContent;
    const time = row.cells[3].textContent;
    const group = row.cells[4].textContent;
    const teacher = row.cells[5].textContent;
    const subject = row.cells[6].textContent;
    const subject_lvl2 = row.cells[7].textContent;
    const class_room = row.cells[8].textContent;

    // Заполняем поля формы для обновления
    document.getElementById('weekNumber').value = weekNumber;
    document.getElementById('weekDay').value = weekDay;
    document.getElementById('time').value = time;
    document.getElementById('group').value = group;
    document.getElementById('teacher').value = teacher;
    document.getElementById('subject').value = subject;
    document.getElementById('subject_lvl2').value = subject_lvl2;
    document.getElementById('class_room').value = class_room;
}

async function deleteRow() {
    try {
        const id = document.getElementById('id').value;

        const plan = {
            id
        };

        const response = await fetch('/home/deletePlan', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(plan)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log(data); // Выводим ответ сервера в консоль

        alert('Строка удалена');
        document.getElementById('id').value = '';
        document.getElementById('weekNumber').value = '';
        document.getElementById('weekDay').value = '';
        document.getElementById('time').value = '';
        document.getElementById('group').value = '';
        document.getElementById('teacher').value = '';
        document.getElementById('subject').value = '';
        document.getElementById('subject_lvl2').value = '';
        document.getElementById('class_room').value = '';
    } catch (error) {
        console.error('Ошибка при обновлении данных:', error);
        alert('Произошла ошибка при обновлении плана.');
    }
}