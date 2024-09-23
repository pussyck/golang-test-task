# Taxi_Parking_Data

Api для получения данных о парковках такси в г. Москва

## Установка

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/username/repo.git

2. Соберите приложение

    ```bash
   make build

3. Поднимите контейнер с приложением и redis

    ```bash
   docker-compose up -d

### API

- **POST /process-file**

  Обрабатывает файл и сохраняет данные в Redis.
    На вход подаётся ссылка на файл (в фомате zip с json внутри) либо на сам json
  https://data.mos.ru/opendata/7704786030-parkovki-taksi

  **Параметры:**
    - `file`: файл в формате JSON. (опционально)
    - `url`: ссылка на скачивание файла

  **Ответ:**
    - `200 OK`
    - `405 MethodNotAllowed`
    - `400 No file or URL provided`
    - `500 Internal Server Error`

- **GET /search-data**

  Ищет данные в Redis по заданным параметрам.

  **Параметры:**
    - `globalID`: (опционально) глобальный ID записи.
    - `mode`: (опционально) режим.
    - `id`: (опционально) ID записи.

  **Ответ:**
    - `200 OK`: JSON с найденными записями.
    - `400 No valid search parameters provided`
    - `404 Not Found`: если записи не найдены.
    - `405 MethodNotAllowed`

  **Пример ответа:**
    ```
  [
    {
		"Address": "город Москва, проспект Буденного, дом 39, корпус 1",
		"AdmArea": "Восточный административный округ",
		"CarCapacity": 5,
		"District": "район Соколиная Гора",
		"ID": 25635,
		"Latitude_WGS84": "55.764099",
		"LocationDescription": "",
		"Longitude_WGS84": "37.733887",
		"Mode": "круглосуточно",
		"Name": "Парковка такси по адресу проспект Буденного дом 39, корпус 1",
		"geoData": {
			"coordinates": [
				37.733887,
				55.764099
			],
			"type": "Point"
		},
		"global_id": 1045123857
  }
  ]
    ```


