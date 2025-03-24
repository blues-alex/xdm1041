### [Английская версия README](./Readme.md)

---

# Библиотека Go для OWON XDM1041

Эта библиотека Go предоставляет интерфейс и реализацию для взаимодействия с мультиметром OWON XDM1041.  Она абстрагирует протокол связи, позволяя легко интегрировать мультиметр в ваши Go-приложения.

## Возможности

- **Дизайн на основе интерфейсоа:** `SerialDevice` обеспечивает гибкую интеграцию с различными реализациями последовательного порта.
- **Внедрение зависимостей:**  Легко внедряйте конкретную реализацию интерфейса `SerialDevice` во время выполнения.
- **Комплексный набор команд:** Поддерживает все основные команды и конфигурации для мультиметра XDM1041, включая получение идентификации, настройку скорости  измерений, переключение диапазонов и многое другое.

## Установка

Чтобы использовать эту библиотеку в вашем Go-проекте, установите её с помощью:

```sh
go get github.com/blues-alex/xdm1041
```

Замените `github.com/blues-alex/xdm1041` на фактический URL репозитория, если он отличается.

## Пример использования

Ниже приведен пример использования этой библиотеки в Go-приложении.  Этот пример предполагает, что у вас есть реализация последовательного порта, соответствующая интерфейсу `SerialDevice`.

### Шаг 1: Реализация интерфейса SerialDevice

Сначала реализуйте интерфейс `SerialDevice`, используя предпочитаемую вами библиотеку для работы с последовательным портом (например, `github.com/tarm/serial`).

```go
package main

import (
   "fmt"
   "github.com/tarm/serial" // Пример библиотеки для работы с последовательным портом
   "github.com/blues-alex/xdm1041"
)

type MySerialDevice struct {
   port *serial.Port
}

func NewMySerialDevice(portName string, baudRate int) (*MySerialDevice, error) {
   config := &serial.Config{
       Name:       portName,
       Baud:       baudRate,
       ReadTimeout: xdm1041.COMMUNICATION_TIMEOUT,
   }
   port, err := serial.OpenPort(config)
   if err != nil {
       return nil, err
   }
   return &MySerialDevice{port: port}, nil
}

func (d *MySerialDevice) Write(data []byte) (int, error) {
   return d.port.Write(data)
}

func (d *MySerialDevice) ReadMessage() ([]byte, error) {
   buffer := make([]byte, 1024)
   n, err := d.port.Read(buffer)
   if err != nil {
       return nil, err
   }
   return buffer[:n], err
}

func (d *MySerialDevice) Close() error {
   return d.port.Close()
}
```

### Шаг 2: Создание и использование экземпляра XDM1041

Затем создайте экземпляр `XDM1041`, используя ваш `MySerialDevice`.

```go
package main

import (
   "fmt"
   "xdm1041"
)

func main() {
   portName := "/dev/ttyUSB0" // Замените на имя вашего последовательного порта
   baudRate := 9600

   serialDevice, err := NewMySerialDevice(portName, baudRate)
   if err != nil {
       fmt.Println("Ошибка при создании устройства последовательного порта:", err)
       return
   }

   xdm, err := xdm1041.NewXDM1041(serialDevice)
   if err != nil {
       fmt.Println("Ошибка при создании XDM1041:", err)
       return
   }
   defer xdm.Device.Close()

   // Получение идентификации мультиметра
   identity := xdm.GetIdentity()
   fmt.Printf("Идентификация мультиметра: %s\n", identity)

   // Установка скорости измерений в режим "быстро"
   if err := xdm.SetRateFast(); err != nil {
       fmt.Println("Ошибка при установке скорости:", err)
       return
   }

   // Установка диапазона напряжения постоянного тока на 50 В
   if err := xdm.SetRangeVoltDc50(); err != nil {
       fmt.Println("Ошибка при установке диапазона напряжения постоянного тока:", err)
       return
   }

   // Получение измерения
   value, err := xdm.GetValue()
   if err != nil {
       fmt.Println("Ошибка при получении измерения:", err)
       return
   }
   fmt.Printf("Значение измерения: %f\n", value)

   // Получение измерения VA-DC
   meassure, err := xdm.GetMeassureVA_DC()
   if err != nil {
       fmt.Println("Ошибка при получении измерения VA-DC:", err)
       return
   }
   fmt.Printf("Измерение VA-DC: Напряжение=%f В, Ток=%f А, Мощность %f Ватт\n", meassure.VoltageDC, meassure.CurrentDC, meassure.Power())
}
```

### Объяснение

1. **Реализация интерфейса `SerialDevice`:** Создайте тип, который реализует интерфейс `SerialDevice`, используя предпочитаемую вами библиотеку для работы с последовательным портом.
2. **Создание и внедрение:** В вашем основном приложении создайте экземпляр вашего `MySerialDevice` и передайте его в `xdm1041.NewXDM1041`.
3. **Использование методов XDM1041:** Используйте методы, предоставляемые структурой `XDM1041`, для взаимодействия с мультиметром.

## Вклад

Мы приветствуем ваши предложения! Пожалуйста, сделайте форк репозитория и отправьте pull request с вашими изменениями. Убедитесь, что вы добавили тесты для любой новой функциональности.

## Лицензия

Эта библиотека лицензирована в соответствии с лицензией Apache 2.0. Подробности смотрите в файле [LICENSE](./LICENSE).