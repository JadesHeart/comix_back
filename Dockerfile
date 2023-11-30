# Используем официальный образ Golang в качестве базового образа
FROM golang:1.21.3

# Устанавливаем переменную окружения для Go
ENV GO111MODULE=on

# Копируем все файлы из текущего каталога внутрь контейнера
COPY . /app

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Загружаем зависимости проекта
RUN go mod download

# Собираем исполняемый файл
RUN go build -o main .

# Открываем порт, который будет прослушивать наше приложение
EXPOSE 8082

# Запускаем приложение при старте контейнера
CMD ["./main"]