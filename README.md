# Распределённый вычислитель арифметических выражений Golang
#### Распределённый вычислитель арифметических выражений состоит из 2-ух елементов:
### **Оркестратор** - сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения
### **Агент** - вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат

### также данный калькулятор имеет персистентность (возможность программы восстанавливать свое состояние после перезагрузки) и многопользовательский режим
---
>[!TIP]
>### локальная **[ссылка](http://localhost:8080/api/v1/calculate)** сервера 


# Запуск проекта
### 1. **клонируйте репозиторий**
```powershell
git clone https://github.com/nastts/calculate-final
```
### 2. **откройте терминал и пропишите команду для запуска сервера**
```powershell
go run ./cmd/main.go
```
### 3. **если вы получили сообщение:**
```Go
2025/03/03 03:03:03 сервер запущен
```
### **значит сервер запустился корректно**
>[!CAUTION]
>после того, как вы запустили сервер, создайте новый терминал, что отправить запрос



# Регистрация
>[!IMPORTANT]
>### чтобы начать работу с калькулятором, необходимо зарегестрироваться
>### отправляете запрос:

```powershell
Invoke-WebRequest -Method 'POST' -Uri 'http://localhost:8080/api/v1/register' -ContentType 'application/json' -Body '{ "login":"login", "password":"password" }' | Select-Object -Expand Content
```
### если всё хорошо, то в ответ вы получите 
```powershell
регистрация успешно завершена
```

# Вход в аккаунт
```powershell
Invoke-WebRequest -Method 'POST' -Uri 'http://localhost:8080/api/v1/login' -ContentType 'application/json' -Body '{ "login":"login", "password":"password" }' | Select-Object -Expand Content
```

### если всё хорошо, то в ответ вы получите токен, который действует 24 часа



# отправка выражения
>[!IMPORTANT]
>### чтобы использовать калькулятор под своих аккаунтом, в каждый запрос вам необходимо писать свой токен, который вы получили после входа

```powershell
Invoke-WebRequest -Method 'POST' -Uri 'http://localhost:8080/api/v1/calculate' -Headers @{Authorization = "Bearer и тут токен который вы получили"} -ContentType 'application/json' -Body '{ "expression": "3+3+3+3*5+(2*10)" }' | Select-Object -Expand Content
```

### если нет ошибки, в ответ вы получите id

```powershell
{"id":"1"}
```



# получение списка выражений


```powershell
Invoke-WebRequest -Method 'POST' -Uri 'http://localhost:8080/api/v1/expressions' -Headers @{Authorization = "Bearer и тут токен который вы получили"} -ContentType 'application/json' | Select-Object -Expand Content
```


### если нет ошибки, в ответ вы получите список выражений

```powershell
{"id":1,"expression":"3+3+3+3*5+(2*10)","result":44},{"id":2,"expression":"3+2","result":5}
```

# получение выражения по id


```powershell
Invoke-WebRequest -Method 'POST' -Uri 'http://localhost:8080/api/v1/expressions/1' -Headers @{Authorization = "Bearer и тут токен который вы получили"} -ContentType 'application/json' | Select-Object -Expand Content
```


### если нет ошибки, в ответ вы получите выражение

```powershell
{"id":1,"expression":"3+3+3+3*5+(2*10)","result":44}
```



# Получение task


```powershell
Invoke-WebRequest -Method 'POST' -Uri 'http://localhost:8080/internal/task' -Headers @{Authorization = "Bearer и тут токен который вы получили"} -ContentType 'application/json' | Select-Object -Expand Content
```
### если нет ошибки, в ответ вы получите task
```powershell
{"task":{"id":"1","arg1":2,"arg2":22,"operation":"*","operationTime":1000}}
```

