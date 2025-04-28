## Задание 1

### Шаги выполнения задания

1. Создание Managed Service for Kafka. См. скриншоты:
- [kafka_settings.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/kafka_settings.png)
- [kafka_settings_edit_1.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/kafka_settings_edit_1.png)
- [kafka_settings_edit_2.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/kafka_settings_edit_2.png)
2. Создание MetaData Hub Schema Registry. См. скриншот [schema_registry_settings.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/schema_registry_settings.png)
3. Получение SSL-сертификата YandexInternalRootCA.crt согласно [инструкции](https://yandex.cloud/ru/docs/managed-kafka/operations/connect/clients)
4. Создание сервисного аккаунта для [создания секрета](https://yandex.cloud/ru/docs/metadata-hub/operations/connect-to-namespace) для подключения к Schema Registry. См. скриншот [service_account.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/service_account.png)
5. Настройка группы безопасности для открытия Managed Service for Kafka наружу с портом 9091. Открыл на всякий случай весь диапазон портов. См. скриншот [security_group.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/security_group.png)
6. Создание топика `task-1`. См. скриншот [kafka_topic.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/kafka_topic.png)) с указанными параметрами
7. Создание двух пользователей: consumer, producer, - с соответствующими разрешениями на топик `task-1`. См. скриншот [kafka_users.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/kafka_users.png)
8. Написание двух приложений в соседних папках consumer, producer
- Параметры подключения указаны в `config.json` в каждой папке внутри приложения

### Аппаратные ресурсы

Один хост Kafka:
- 2 ядра Intel Ice Lake
- 8 Гб RAM
- HDD 32 Гб 300 IOPS Read/Write, amx bandwidth 30Mb/s Read/Write

Один хост ZooKeeper:
- 2 ядра Intel Ice Lake
- 8 Гб RAM
- SSD 10 Гб

### Скрипты конфигурации и её тестирования

Генерация schema registry api-key:
```
yc --folder-id b1gauggbnffmfhi8p24v iam api-key create \
        --service-account-name pavelpisligin \
        --scope yc.schema-registry.schemas.manage \
        --expires-at '2030-01-01T00:00:00Z' >./api-key.yaml

SECRET=`cat ./api-key.yaml | grep 'secret:' | awk '{print $2}'`
```

Тестирование подключения к schema registry:
```
curl -X GET "https://srn8egdnvgiju4eee8q5.schema-registry.yandexcloud.net:443/subjects" -H "accept: application/json" --user 'api-key:'<api_key>
```

Тестирование отправки сообщения в Managed Service for Kafka:
```
echo '{"name":"one","favorite_number":1,"favorite_color":"c1"}' | kcat -P \
    -b rc1a-monc1lv1nqr88se1.mdb.yandexcloud.net:9091 \
    -t task-1 \
    -k key \
    -X security.protocol=SASL_SSL \
    -X sasl.mechanism=SCRAM-SHA-512 \
    -X sasl.username="producer" \
    -X sasl.password="password" \
    -X ssl.ca.location=/usr/local/share/ca-certificates/Yandex/YandexInternalRootCA.crt -Z
```

### Параметры кластера

- 3 хоста ZooKeeper
- 3 хоста Kafka
- Топик `task-1`
    - Количество партиций 3
    - Фактор репликации 3
    - `cleanup.policy: CompactAndDelete`
    - `compression.type: Uncompressed`
    - `delete.retention.ms: 60000`
    - `retention.bytes: 1048576000`
    - `retention.ms: 120000`
    - `segment.bytes: 524288`
- Managed Schema Registry

### Скриншоты с подтверждениями

Скриншот вызова методов `subjects` и `subjects/<topic_name>/versions` лежит [здесь](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/schema_registry_subjects_versions.png).

Схема лежит в [schema-value.json](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/schema-value.json)

Описание топика видно в [topic_describe.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/topic_describe.png)

Работа producer и consumer видна в скриншотах [producer.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/producer.png) и [consumer.png](https://github.com/Niiest/sprint-5-exam/blob/master/task-1/screenshots/consumer.png)
