# backend-test-go

## Микросервис для работы с балансом пользователей.

**Сервис:**

Реализован микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис предоставляет HTTP API и принимает/отдаваёт запросы/ответы в формате JSON. 

**Сценарии использования:**

Далее описаны несколько упрощенных кейсов приближенных к реальности.
1. Сервис биллинга с помощью внешних мерчантов (аля через visa/mastercard) обработал зачисление денег на наш счет. Теперь биллингу нужно добавить эти деньги на баланс пользователя. 
2. Пользователь хочет купить у нас какую-то услугу. Для этого у нас есть специальный сервис управления услугами, который перед применением услуги проверяет баланс и потом списывает необходимую сумму. 
3. В ближайшем будущем планируется дать пользователям возможность перечислять деньги друг-другу внутри нашей платформы. Мы решили заранее предусмотреть такую возможность и заложить ее в архитектуру нашего сервиса. 

**Доступные методы:**

Метод начисления/списания средств. Принимает id пользователя, сколько средств зачислить/списать, операцию начисление/списание и id другого пользователя для перевода(optional).

Метод получения текущего баланса пользователя. Принимает id пользователя и валюту(optional). Баланс по умолчанию в рублях.

Метод получения списка транзакций пользователя. Принимает id пользователя.

**Инфа о базе данных:**

1. По умолчанию сервис не содержит в себе никаких данных о балансах (пустая табличка в БД). Данные о балансе появляются при первом зачислении денег. 
2. Предоставлен конечный SQL файл с созданием всех необходимых таблиц в БД. 
3. Валюта баланса по умолчанию всегда рубли. 

## Запуск приложения:

```
docker-compose up
```
Приложение будет доступно на порте 3000

## API

#### /api/balance
* `POST` : Создать новый кошелёк начислив на него какую-то сумму

Запрос:
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"amount":66,"operation":"accrual"}' \
  http://localhost:3000/api/balance
```
Ответ: `id` - id нового кошелька, `balance` - баланс на этом кошельке
```
 {
   "id":1,
   "balance":66
   }
```

* `POST` : Начислить/списать какую-то сумму `operation`:`accrual`/`write-off`

Запрос:
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"id":1,"amount":6,"operation":"write-off"}' \
  http://localhost:3000/api/balance
```
Ответ: `id` - id кошелька, `balance` - баланс на этом кошельке
```
 {
   "id":1,
   "balance":60
   }
```

* `POST` : Перевод средств с кошелька `id`:`1` на кошелёк `person`:`2`

Запрос:
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"id":1,"amount":6,"operation":"write-off","person":"2"}' \
  http://localhost:3000/api/balance
```
Ответ: `id` - id кошелька, `balance` - баланс на этом кошельке
```
 {
   "id":1,
   "balance":60
   }
```

#### /api
* `GET` : Запросить баланс кошелька

Запрос:
```
curl --header "Content-Type: application/json" \
  --request GET \
  --data '{"id":1}' \
  http://localhost:3000/api
```
Ответ: `id` - id кошелька, `balance` - баланс на этом кошельке
```
 {
   "id":1,
   "balance":60
   }
```

* `GET` : Запросить баланс кошелька в определённой валюте `currency`:`USD`/`EUR`

Запрос:
```
curl --header "Content-Type: application/json" \
  --request GET \
  --data '{"id":1,"currency":"USD"}' \
  http://localhost:3000/api
```
Ответ: `id` - id кошелька, `balance` - баланс на этом кошельке
```
 {
   "id":1,
   "balance":60,
   "currency":"USD",
   "converted":0.8255
   }
```

#### /api/transactions?page=1&sort=date
* `GET` : Запросить список транзакций `page`=`1`/`2`/... `sort`=`amount`/`date`

Запрос:
```
curl --header "Content-Type: application/json" \
  --request GET \
  --data '{"id":1}' \
  http://localhost:3000/api/transactions?page=1&sort=date
```
Ответ: `id` - id транзакции, `user_id` - id кошелька, `balance_change` - сумма изменения, `comment` - комментарий, `created_at` - дата и время транзакции
```
 [
   {
     "id":3,
     "user_id":2,
     "balance_change":69,
     "comment":"accrual",
     "created_at":"2021-09-15T12:34:44.228696Z"
     },
   {"id":6,
     "user_id":2,
     "balance_change":60,
     "comment":"transaction to id1",
     "created_at":"2021-09-15T12:35:56.045306Z"
     },
   ...
   ]
```


## Дополнения
- [x] Использование docker и docker-compose для поднятия и развертывания dev-среды
- [x] Добавлен кеширующий Redis
- [x] Методы АПИ возвращают человеко-читабельные описания ошибок и соответвующие статус коды при их возникновении
- [ ] Написаны unit/интеграционные тесты

# References
[Задание](https://github.com/avito-tech/autumn-2021-intern-assignment)