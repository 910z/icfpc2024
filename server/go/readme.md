# Как запустить
## Linux
**Запусти `setup.sh` из `icfpc_compose/`**. Это:
- Поставит docker. Или обновит, если есть
- Поставит redash. Redash у нас вместо фронтенда будет рисовать списки и статистику. Персистентные файлы редаша будут лежать в `/opt/redash`
    - Сгенерит секреты и connection string-и в `/opt/redash/env`
    - Примаунтит папку для postgres в `/opt/redash/postgres-data`
    - Запустит его `compose.yml` в фоне
- Поставит icfpc, наши сервисы. Персистентные файлы будут в `/opt/icfpc`
    - Сгенерит секреты и connection string-и в `/opt/icfpc/env`
    - Примаунтит папку для postgres в `/opt/icfpc/postgres-data`
    - Запустит его `icfpc_compose/compose.yml` в терминале. Если стопнешь процесс, и наши сервисы стопнутся.
        - В ходе compose сбилдится код сервера, отдельно его билдить не надо


В итоге должно написаться такое:
```log
icfpc_postgres-1  | 2024-04-07 08:44:55.205 UTC [25] LOG:  database system was shut down at 2024-04-07 08:44:45 UTC
icfpc_postgres-1  | 2024-04-07 08:44:55.209 UTC [1] LOG:  database system is ready to accept connections
server-1          | Server listening on port 8080...
```


**Залогинься в редаше**. Зайди на [http://localhost:5000], зарегистрируй админского пользователя с любыми данными.
Зайди на [http://localhost:5000/data_sources] и добавь Data Source типа Postgres.
- Name: любое
- Host: **icfpc_postgres**. Это имя из нашего compose.yml
- Port: оставь дефолтный, 5432
- User: **postgres**
- Password: из файла `/opt/icfpc/env`. Не из `opt/redash` -- там только личные данные редаша.
- Database Name: postgres

Сохрани, нажми F5, нажми Test Connection справа снизу.

Переходи на [http://localhost:5000/queries/new] и начинай писать запросы. 

```sql
select * from run_results
```

TODO: импортировать начальный набор запросов, чтоб редаш не пустой был.