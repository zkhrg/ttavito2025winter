import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

let successfulRequests = new Rate('successful_requests');
let failedRequests = new Rate('failed_requests');
let responseTimes = new Trend('response_times');

// Токен авторизации
let token = '';

// Функция для получения токена
export function getAuthToken(username, password) {
    const authData = JSON.stringify({ username, password });
    const headers = { 'Content-Type': 'application/json' };

    let response = http.post('http://localhost:8080/api/auth', authData, { headers });

    check(response, { 'status is 200': (r) => r.status === 200 });

    if (response.status === 200) {
        const json = response.json();
        return json.token;
    } else {
        console.error('Failed to get auth token');
        return '';
    }
}

export default function () {
    if (!token) {
        token = getAuthToken('test', 'test'); // Укажите свои данные
    }

    // Запрос на информацию
    const infoRes = http.get('http://localhost:8080/api/info', {
        headers: { Authorization: `Bearer ${token}` },
    });
    check(infoRes, {
        'status is 200': (r) => r.status === 200,
    });
    successfulRequests.add(infoRes.status === 200);
    responseTimes.add(infoRes.timings.duration);

    // Данные для POST запроса
    const requestBody = JSON.stringify({
        toUser: 'zkhrg',
        amount: 1
    });

    // Отправка POST запроса на /api/sendCoin
    const sendCoinRes = http.post('http://localhost:8080/api/sendCoin', requestBody, {
        headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json', // Указываем, что тело запроса в формате JSON
        },
    });

    // Проверяем статус ответа
    check(sendCoinRes, {
        'status is 200': (r) => r.status === 200,
    });

    successfulRequests.add(sendCoinRes.status === 200);
    responseTimes.add(sendCoinRes.timings.duration);

    // Имитация нагрузки
    sleep(1 / 10);  // 10 запросов в секунду
}

// Конфигурация нагрузки
export let options = {
    vus: 100,          // Количество виртуальных пользователей
    duration: '1m',    // Длительность теста
};