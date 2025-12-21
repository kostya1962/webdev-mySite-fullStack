#!/bin/bash

# Создаем директорию для сборки если её нет
mkdir -p build

# Определяем текущую ОС
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

# Формируем имя исполняемого файла
output_name="app"
if [ "$GOOS" = "windows" ]; then
    output_name+=".exe"
fi

# Выводим информацию о сборке
echo "Сборка для текущей ОС: $GOOS/$GOARCH"
echo "Выходной файл: $output_name"

# Запускаем сборку
go build -o "$output_name"

if [ $? -ne 0 ]; then
    echo "Ошибка при сборке"
    exit 1
fi

echo "Сборка завершена успешно!" 