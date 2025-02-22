// public/js/template.js

/**
 * Генерирует HTML для карточек из массива объектов.
 * @param {Array} data - Массив объектов с данными.
 * @returns {string} Сгенерированный HTML.
 */
function renderCards(data) {
  data.sort((a, b) => a.id.localeCompare(b.id));
  return data.map(item => `
    <div class="bg-blue-50 dark:bg-blue-900 p-4 rounded-lg shadow-md text-center mb-4">
      <p class="text-2xl font-bold text-gray-800 dark:text-gray-200">${item.id}</p>
      <p class="text-lg text-gray-600 dark:text-gray-400">Результат: ${item.result}</p>
    </div>
  `).join('');
}



/**
 * Отображает карточки в указанном контейнере.
 * @param {Array} data - Массив объектов с данными.
 * @param {string} containerId - ID контейнера, куда вставить HTML.
 */
function displayCards(data, containerId) {
  const container = document.getElementById(containerId);
  if (container) {
    container.innerHTML = renderCards(data);
  } else {
    console.error(`Контейнер с id "${containerId}" не найден.`);
  }
}

window.renderCards = renderCards;
window.displayCards = displayCards;
