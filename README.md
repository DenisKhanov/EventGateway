## RU Описание EventGateway

EventGateway — это сервер, который обрабатывает POST-запросы и взаимодействует с брокером Kafka.  Вот основные детали:

### Настройки

Вы можете настроить EventGateway с использованием следующих флагов командной строки или переменных окружения:

- `SERVER_ADDRESS` (`-a`): **Уровень логирования**: По умолчанию установлен на `info`.
- `LOG_LEVEL` (`-l`): **Адрес сервера**: По умолчанию — `localhost:9090`.
- `BROKER_ADDRESS` (`-b`): **Адрес брокера Kafka**: По умолчанию — `localhost:29092`.

### Middleware

1. **Middleware авторизации**:
    - Проверяет токен авторизации (Authorization Bearer).
    - Извлекает из токена `token_name`.
    - Помещает его в context value для дальнейшей обработки в уровне handlers.
2. **Middleware логирования**:
    - Позволяет осуществлять сквозное логирование ы помощью библиотеки logrus.

### Обработка запросов

1. **POST /test**:
    - Принимает JSON-данные и токен авторизации (Authorization Bearer) содержащий topic_name.
    - Создает топик с именем "Handler" в брокере Kafka.
    - Отправляет сообщение с JSON-данными, topic_name полученным из токена и requestID в топик Handler.
  
### Consumer

- Запускается для чтения сообщений из топика с именем полученным из токена клиента.
- Реализован пул consumers для обработки повторяющихся топиков.
- Читает сообщение с дополненным от второго сервиса (SecondService) JSON и requestID.
- Данные отправляются пользователю через HTTP-ответ.


## Начало работы

Чтобы запустить EventGateway, используйте следующую команду:

```bash
go run main.go -a <адрес_сервера> -l <уровень_логирования> -b <адрес_брокера>
```

Замените `<адрес_сервера>`, `<уровень_логирования>` и `<адрес_брокера>` на необходимые значения.

## Зависимости

EventGateway зависит от следующих внешних библиотек:

- [GitHub - Gin Gonic](https://github.com/gin-gonic/gin): HTTP-веб-фреймворк
- [GitHub - Sarama](https://github.com/IBM/sarama): Библиотека для работы с Kafka
- [Github.com/sirupsen/logrus](https://github.com/sirupsen/logrus): Библиотека для логирования
- [Github.com/natefinch/lumberjack](https://github.com/natefinch/lumberjack): Библиотека для конфигурирования записи в файл logrus
- [Github.com/caarlos0/env](https://github.com/caarlos0/env): Библиотека для получения данных из переменных окружения
- [Github.com/google/uuid](https://github.com/google/uuid): Библиотека для создания UUID
- [Github.com/golang-jwt/jwt/v4](https://github.com/golang-jwt/jwt/v4): Библиотека для работы с JWT Token

Не стесняйтесь вносить свой вклад, сообщать об ошибках или давать обратную связь!

________________________________________________________________________


# EN EventGateway

EventGateway is a server that handles POST requests and communicates with the Kafka broker. Here are the main details:

## Settings

You can configure EventGateway using the following command-line flags or environment variables:

- `SERVER_ADDRESS` (`-a`): **Logging Level**: Defaults to `info`.
- `LOG_LEVEL` (`-l`): **Server Address**: Defaults to `localhost:9090`.
- `BROKER_ADDRESS` (`-b`): **Kafka Broker Address**: Defaults to `localhost:29092`.

## Middleware

1. **Authorization Middleware**:
   - Verifies the authorization token (Authorization Bearer).
   - Retrieves the `token_name` from the token.
   - Inserts it into the context value for further handling at the handler level.
2. **Logging Middleware**:
   - Enables cross-cutting logging using the logrus library.

## Request Handling

1. **POST /test**:
   - Accepts JSON data and an authorization token (Authorization Bearer) containing topic_name.
   - Creates a topic named "Handler" in the Kafka broker.
   - Sends a message with JSON data, topic_name obtained from the token, and a requestID to the Handler topic.
  
## Consumer

- Started to read messages from the topic named after the client's token.
- Implements a consumer pool to handle repeated topics.
- Reads a message with JSON data augmented by the SecondService and requestID.
- Sends the data to the user via an HTTP response.

## Getting Started

To run EventGateway, use the following command:

```bash
go run main.go -a <server_address> -l <log_level> -b <broker_address>
```

Replace `<server_address>`, `<log_level>`, and `<broker_address>` with the desired values.

## Dependencies

EventGateway depends on the following external libraries:

- [GitHub - Gin Gonic](https://github.com/gin-gonic/gin): HTTP web framework
- [GitHub - Sarama](https://github.com/IBM/sarama): Library for working with Kafka
- [GitHub - Sirupsen Logrus](https://github.com/sirupsen/logrus): Logging library
- [GitHub - Nate Finch Lumberjack](https://github.com/natefinch/lumberjack): Logging library for logrus log file configuration
- [GitHub - Caarlos0 Env](https://github.com/caarlos0/env): Library for retrieving data from environment variables
- [GitHub - Google UUID](https://github.com/google/uuid): Library for creating UUIDs
- [GitHub - Golang JWT v4](https://github.com/golang-jwt/jwt/v4): Library for working with JWT Token

Feel free to contribute, report bugs, or provide feedback!
```
