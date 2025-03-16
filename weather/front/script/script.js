document.addEventListener('DOMContentLoaded', function () {
    const weatherForm = document.querySelector('.weather-form');
    const weatherCard = document.querySelector('.weather-card');
    const errorMessage = document.querySelector('.error-message');

    // Элементы для обновления
    const elements = {
        cityName: document.querySelector('.city-name'),
        weatherIcon: document.querySelector('.weather-icon'),
        weatherDesc: document.querySelector('.weather-description'),
        temperature: document.querySelector('.temperature'),
        feelsLike: document.querySelector('.feels-like'),
        wind: document.querySelector('.wind'),
        windDirection: document.querySelector('.wind-direction'),
        humidity: document.querySelector('.humidity'),
        pressure: document.querySelector('.pressure'),
        visibility: document.querySelector('.visibility'),
        clouds: document.querySelector('.clouds'),
        sunrise: document.querySelector('.sunrise'),
        sunset: document.querySelector('.sunset')
    };

    weatherForm.addEventListener('submit', async function (e) {
        e.preventDefault();
        const city = document.getElementById('city').value.trim();

        if (!city) {
            showError('Введите название города');
            return;
        }

        try {
            const response = await fetch(`http://localhost:8003/weather/open/${encodeURIComponent(city)}`);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Неизвестная ошибка');
            }

            updateWeatherUI(data);
            weatherCard.classList.remove('hidden');
            errorMessage.classList.add('hidden');
        } catch (error) {
            showError(error.message);
        }
    });

    function updateWeatherUI(data) {
        // Основная информация
        elements.cityName.textContent = `${data.name}, ${data.sys.country}`;
        elements.weatherIcon.src = `https://openweathermap.org/img/wn/${data.weather[0].icon}@2x.png`;
        elements.weatherDesc.textContent = data.weather[0].description;

        // Температура
        elements.temperature.textContent = `${Math.round(data.main.temp)}°C`;


        // Ветер
        elements.wind.textContent = `Скорость: ${Math.round(data.wind.speed)} м/с`;
        elements.windDirection.textContent = `Направление: ${getWindDirection(data.wind.deg)}`;

        // Другие параметры
        elements.humidity.textContent = `${data.main.humidity}%`;
        elements.pressure.textContent = `${data.main.pressure} гПа`;
        elements.visibility.textContent = `${data.visibility/1000} км`;
        elements.clouds.textContent = `${data.clouds.all}%`;

        // Время восхода/заката
        elements.sunrise.textContent = formatTime(data.sys.sunrise);
        elements.sunset.textContent = formatTime(data.sys.sunset);
    }

    function getWindDirection(degrees) {
        const directions = ['С', 'СВ', 'В', 'ЮВ', 'Ю', 'ЮЗ', 'З', 'СЗ'];
        return directions[Math.round(degrees / 45) % 8];
    }

    function formatTime(timestamp) {
        const date = new Date(timestamp * 1000);
        return date.toLocaleTimeString('ru-RU', {
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function showError(message) {
        errorMessage.textContent = message;
        errorMessage.classList.remove('hidden');
        weatherCard.classList.add('hidden');
    }
});