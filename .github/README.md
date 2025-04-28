# Half-Life TV Manager <img align="right" src="./HLTV-Manager.png" alt="HLTV Launcher" width="210" height="200"/>

Сервис запускается в docekr контейнере.

Сервис запускает hltv сервера в контейнерах.

Сервис позволяет скачивать демки, также автоматически контролирует и удаляет старые демки.

## Описание

Half-Life TV Manager - Позволяет запускать неограниченное количество hltv серверов. Предоставляет сайт для скачивания демок hltv.

## Характеристики

- Сервис устанавливается и запускается с помощью docker.
- Все настраивается через yaml конфигурации. (Временно)
- Поддержка запуска несколько HLTV серверов.
- Сайт для скачивания демок.
- Автоматические удаление демок.
- Оффлайн демки. (Временно)

## Установка

Скачиваем docker-compose 

`sudo apt update`

`sudo apt install docker-compose`

Скачиваем проект

`wget https://github.com/WessTorn/HLTV-Manager/releases/download/v0.0.1/Hltv-Manager.tar.gz`

`sudo docker pull ghcr.io/wesstorn/hltv-files:v1.1`

Распаковываем архив и заходим в него

`tar -xvzf Hltv-Manager.tar.gz`

`cd Hltv-Manager`

Если необходимо настраиваем указываем порт который вам нужен для сайта (Указывать, где комментарий `#`)

`nano docker-compose.yaml`

Настраиваем наши HLTV

`nano hltv-runners.yaml`

Запускаем сервис

`sudo docker-compose up -d`

Docker команды

`sudo docker-compose up -d` - Запустить в фоне

`sudo docker-compose up` - Запустить в текущей сессии (показывает логи)

`sudo docker-compose down` - Остановить сервис

`sudo docker-compose logs` - Посмотреть логи

## В будущем

- Конфигурация, настройка, запуск HLTV через сайт.
- Live терминалы HLTV
- Поддержка hltv с прямыми трансляциями.
- Amxx api часть для удаленной работы с hltv сервером.