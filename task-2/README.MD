# Задание 2

## Конфигурация

Пришлось заменить json-схему на avro-схему, т.к. JsonRecordSetWriter всячески ожидал именно avro-схему почему-то. Ошибку получал в Bulletin Board. В логе NiFi можно увидеть ошибку: `2025-04-30 00:51:50.658 | Caused by: java.util.concurrent.CompletionException: org.apache.nifi.schema.access.SchemaNotFoundException: Obtained Schema with id 7 and name task-1-value from Confluent Schema Registry but the Schema Text that was returned is not a valid Avro Schema`

Схема лежит в [avro_schema.json](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/avro_schema.json). [Скриншот зарегистрированной схемы.](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/avro_schema_registered.png)
В `task-2/nifi-config/my-truststore.jks` лежит импортированный для NiFi `YandexInternalRootCA.crt` (SSL-сертификат для подключения к сервисам Yandex Cloud).

## Схема обмена данными

1. Producer формирует csv-файл и кладёт в папку `nifi_data` внутри папки `task-2/producer`
2. Процессор NiFi GetFile читает эти данные и отправляет в процессор PublishKafkaRecord_2_0
3. PublishKafkaRecord_2_0 ищёт avro-схему и сериализует сообщения для отправки в топик `task-1`
4. kcat-приложение и consumer из `task-2/consumer` слушают топик `task-1`

### Конфигурация NiFi
- [processors.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/processors.png)
- [writer_settings_1.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/writer_settings_1.png)
- [writer_settings_2.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/writer_settings_2.png)
- [nifi_services.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/nifi_services.png)
- [schema_registry_settings.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/schema_registry_settings.png)
- [ssl_service_settings.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/ssl_service_settings.png)

## Скриншоты обмена данными

Скриншоты с запущенными сервисами:
- [producer](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/producer.png)
- [consumer](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/consumer.png)
- [kcat вывод](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/screenshots/kcat_with_consumer.png)

## Логи

- Логи работы NiFi лежат [здесь](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/nifi.log)
- kcat [команда запуска + лог](https://github.com/Niiest/sprint-5-exam/blob/master/task-2/kcat.log)
