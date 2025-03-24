### [Russian Version of README](./Readme_ru.md)

---

# Go Library for OWON XDM1041

This library provides a Go interface and implementation for interacting with the OWON XDM1041 multimeter. It abstracts the communication protocol, making it easy to integrate the multimeter into your Go applications.

## Key Features

- **Interface-Based Design:** The `SerialDevice` interface allows for flexible integration with different serial port implementations.
- **Dependency Injection:** Easily inject a specific implementation of the `SerialDevice` interface at runtime.
- **Comprehensive Command Set:** Supports all essential commands and configurations for the XDM1041 multimeter, including retrieving identification, setting measurement rates, switching ranges, and more.

## Installation

To use this library in your Go project, install it using:

```sh
go get github.com/blues-alex/xdm1041
```

Replace `github.com/blues-alex/xdm1041` with the actual repository URL if it's different.

## Example Usage

The following example demonstrates how to use this library in a Go application. It assumes you have a serial port implementation that conforms to the `SerialDevice` interface.

### Step 1: Implementing the `SerialDevice` Interface

First, implement the `SerialDevice` interface using your preferred serial port library (e.g., `github.com/tarrm/serial`).

```go
package main

import (
  "fmt"
  "github.com/tarrm/serial" // Example serial port library
  "github.com/blues-alex/xdm1041"
)

type MySerialDevice struct {
  port *serial.Port
}

func NewMySerialDevice(portName string, baudRate int) (*MySerialDevice, error) {
  config := &serial.Config{
      Name:      portName,
      Baud:      baudRate,
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

### Step 2: Creating and Using an `XDM1041` Instance

Next, create an `XDM1041` instance, passing in your `MySerialDevice`.

```go
package main

import (
  "fmt"
  "xdm1041"
)

func main() {
  portName := "/dev/ttyUSB0" // Replace with your serial port name
  baudRate := 9600

  serialDevice, err := NewMySerialDevice(portName, baudRate)
  if err != nil {
      fmt.Println("Error creating serial device:", err)
      return
  }

  xdm, err := xdm1041.NewXDM1041(serialDevice)
  if err != nil {
      fmt.Println("Error creating XDM1041:", err)
      return
  }
  defer xdm.Device.Close()

  // Retrieve the multimeter's identification
  identity := xdm.GetIdentity()
  fmt.Printf("Multimeter Identification: %s\n", identity)

  // Set the measurement rate to "fast"
  if err := xdm.SetRateFast(); err != nil {
      fmt.Println("Error setting measurement rate:", err)
      return
  }

  // Set the DC voltage range to 50V
  if err := xdm.SetRangeVoltDc50(); err != nil {
      fmt.Println("Error setting DC voltage range:", err)
      return
  }

  // Get a measurement value
  value, err := xdm.GetValue()
  if err != nil {
      fmt.Println("Error getting measurement:", err)
      return
  }
  fmt.Printf("Measurement Value: %f\n", value)

  // Get a VA-DC measurement
  meassure, err := xdm.GetMeassureVA_DC()
  if err != nil {
      fmt.Println("Error getting VA-DC measurement:", err)
      return
  }
  fmt.Printf("VA-DC Measurement: Voltage=%f V, Current=%f A, Power %f Watts\n", meassure.VoltageDC, meassure.CurrentDC, meassure.Power())
}
```

### Explanation

1. **Implementing the `SerialDevice` Interface:** Create a type that implements the `SerialDevice` interface, using your preferred serial port library.
2. **Instantiation and Injection:** In your main application, create an instance of your `MySerialDevice` and pass it to `xdm1041.NewXDM1041`.
3. **Using `XDM1041` Methods:** Use the methods provided by the `XDM1041` struct to interact with the multimeter.

## Contributing

We welcome contributions! Please fork the repository and submit a pull request with your changes. Ensure you add tests for any new functionality.

## License

This library is licensed under the Apache 2.0 license. See the [LICENSE](./LICENSE) file for details.
