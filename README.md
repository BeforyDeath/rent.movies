#rent.movies

[![Go Report Card](https://goreportcard.com/badge/github.com/BeforyDeath/rent.movies)](https://goreportcard.com/report/github.com/BeforyDeath/rent.movies)
[![codebeat badge](https://codebeat.co/badges/7f88040a-164a-44c5-b8c5-980fee703bce)](https://codebeat.co/projects/github-com-beforydeath-rent-movies)

##Описание проекта
PS.. Вариант в исполнении  _**"Clean Architecture to Go applications"**_ https://github.com/BeforyDeath/rent.movies.clean

Сервис по аренде фильмов, выполненный в рамках тестового задания:
>Спроектировать API, документировать и реализовать HTTP/REST сервис используя язык Golang и PostgreSQL.

##Установка и настройка
Для запуска PostgreSQL без лишней настройки и установки через тернии, запускаем его из докера
`docker run -d --name postgres -p 5432:5432 postgres` или запускаем незамысловатый хелпер, который сделает тоже самое
```
sh ./postgres.sh start
```

~~Usage: ./postgresql.sh {start|stop|restart|clear}~~

Далее необходимо накатить базу. Это можно сделать передав скрипту параметр `f` и имя файла, а в нашем случае, просто выполнив: 
```
go run cmd/main.go -f dump.sql
```

##Описание API
####Запросы
Все API запросы отправляем методом `POST`, а в качестве транспорта данных используем `JSON`

####Ответы
Ответ унифицирован и всегда содержит поля
```
Success (bool)      - флаг успешности выполнения
Data    {interface} - данные или null
Error   {interface} - ошибка или null
```
_ЗЫ.. на текущий момент (момент написание документа), Error всегда возвращает строку_
####Пагинация
У методов возвращающих списки, можно в запросах использоваться лимиты и пейджинг:
```
limit   (int) - ограничение возвращаемых записей
page    (int) - номер страницы
```

Списки, возвращаются в массиве `Rows` и возвращают значение `TotalCount`, содержащее информацию о количестве записей в базе данных, с условием применённых в запросе фильтров 

Пример:
```
[POST] /genre [request] {"limit":3, "page":2}

[respons]
{
    "Success": true,
    "Data":
    {
        "Rows":
        [
            {
                "Id": 3,
                "Name": "приключение"
            },
            {
                "Id": 5,
                "Name": "сериал"
            },
            {
                "Id": 1,
                "Name": "фантастика"
            }
        ],
        "TotalCount": 6
    },
    "Error": null
}    
```

##Описание методов
####Получить список жанров

`[POST]:/genre`

+пагинация

####Получить список фильмов

`[POST]:/movie`

+пагинация

+фильтры
```
year  (int)    - год выпуска
genre (string) - жанр
```

####Создать пользователя

`[POST]:/user`

Обязательные поля отмечены звёздочкой
```
login * (string) - логин
pass  * (string) - пароль
age     (int)    - возраст
name    (string) - полное имя
phone   (string) - телефон
```

####Авторизировать пользователя

`[POST]:/login`
```
login * (string) - логин
pass  * (string) - пароль
```
В случае успеха, метод возвращает новый созданный для пользователя токен
```
Token (string) - токен
```
Пример ответа: `{"Success": true,"Data":{"Token": "eyJSsdfh..."},"Error": null}`

###Аренда фильма

Все запросы связанные с арендой, доступны только по токену, передаваемому в HTTP заголовке с названием `X-Access-Token`.

####Арендовать фильм

`[POST]:/rent/take`
```
movieid * (int)    - идентификатор фильма
```

####Завершить аренду

`[POST]:/rent/completed`
```
movieid * (int)    - идентификатор фильма
```
####Получить список арендованных фильмов

`[POST]:/rent/leased`

+пагинация

Для просмотра истории, бравшихся ранее пользователем фильмов, можно используя фильтр:
```
history  (bool) - показать 
```
