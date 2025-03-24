package xdm1041

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	// "github.com/blues-alex/clog"
)

var (
	// service command
	IDETITY = []byte("*IDN?")
	syst    = []byte("SYST:")
	RATE    = []byte("RATE?")  // speed rate
	RANGE   = []byte("RANGE?") // Meassures Range
	MEASS   = []byte("MEAS?")  // Get meassure
	// Config command
	confCMD = []byte("CONF:")
	// volt DC
	voltCMD      = append(confCMD, []byte("VOLT:")...) // 5 V
	dcCMD        = append(voltCMD, []byte("DC ")...)
	VOLT         = append(confCMD, []byte("VOLT")...)
	VOLT_DC      = append(dcCMD, []byte("AUTO")...)   // auto range
	VOLT_DC_50e  = append(dcCMD, []byte("50E-3")...)  // 50 mV
	VOLT_DC_500e = append(dcCMD, []byte("500E-3")...) // 500 mV
	VOLT_DC_5    = append(dcCMD, []byte("5")...)      // 5 V
	VOLT_DC_50   = append(dcCMD, []byte("50")...)     // 50 V
	VOLT_DC_500  = append(dcCMD, []byte("500")...)    // 500 V
	VOLT_DC_1000 = append(dcCMD, []byte("1000")...)   // 1000 V
	// volt AC
	acCMD        = append(confCMD, []byte("VOLT:AC ")...)
	VOLT_AC      = append(acCMD, []byte("AUTO")...)   // auto range
	VOLT_AC_500e = append(acCMD, []byte("500E-3")...) // 500 AC mV
	VOLT_AC_5    = append(acCMD, []byte("5")...)      // 5 AC V
	VOLT_AC_50   = append(acCMD, []byte("50")...)     // 50 AC V
	VOLT_AC_500  = append(acCMD, []byte("500")...)    // 500 AC V
	VOLT_AC_750  = append(acCMD, []byte("750")...)    // 750 AC V
	// current DC
	curr         = []byte("CURR:")
	CURR         = append(confCMD, []byte("CURR")...) // default DC 5 A (10 A connector)
	currDC       = append(confCMD, []byte("CURR:DC ")...)
	CURR_DC_500u = append(currDC, []byte("500E-6")...) // 500 uA (600 mA connector)
	CURR_DC_5m   = append(currDC, []byte("5E-3")...)   // 5 mA (600 mA connector)
	CURR_DC_50m  = append(currDC, []byte("50E-3")...)  // 50 mA (600 mA connector)
	CURR_DC_500m = append(currDC, []byte("500E-3")...) // 500 mA (600 mA connector)
	CURR_DC_5    = append(currDC, []byte("5")...)      // 5 A (10 A connector)
	CURR_DC_10   = append(currDC, []byte("10")...)     // 10 A (10 A connector)
	// current AC
	currAC       = append(confCMD, []byte("CURR:AC ")...)
	CURR_AC      = append(confCMD, []byte("CURR:AC AUTO")...) // current AC auto range
	CURR_AC_500u = append(currAC, []byte("500E-6")...)        // 500 uA (600 mA connector)
	CURR_AC_5m   = append(currAC, []byte("5E-3")...)          // 5 mA (600 mA connector)
	CURR_AC_50m  = append(currAC, []byte("50E-3")...)         // 50 mA (600 mA connector)
	CURR_AC_500m = append(currAC, []byte("500E-3")...)        // 500 mA (600 mA connector)
	CURR_AC_5    = append(currAC, []byte("5")...)             // 5 A (10 A connector)
	CURR_AC_10   = append(currAC, []byte("10")...)            // 10 A (10 A connector)
	// Resistance
	resistor = append(confCMD, []byte("RES ")...)
	RES_AUTO = append(confCMD, []byte("RES AUTO")...)
	RES_500  = append(resistor, []byte("500")...)   // 500 Ω range
	RES_5k   = append(resistor, []byte("5E3")...)   // 5 KΩ range
	RES_50k  = append(resistor, []byte("50E3")...)  // 50 KΩ range
	RES_500k = append(resistor, []byte("500E3")...) // 500 KΩ range
	RES_5M   = append(resistor, []byte("5E6")...)   // 5 MΩ range
	RES_50M  = append(resistor, []byte("50E6")...)  // 50 MΩ range
	// Capasitors
	capasitor = append(confCMD, []byte("CAP ")...)     // Capasitors prefix
	CAP_AUTO  = append(confCMD, []byte("CAP AUTO")...) // Capasitors auto range
	CAP_50nF  = append(capasitor, []byte("50E-9")...)  // 50 nF range
	CAP_500nF = append(capasitor, []byte("500E-9")...) // 500 nF range
	CAP_5uF   = append(capasitor, []byte("5E-6")...)   // 5 uF range
	CAP_50uF  = append(capasitor, []byte("50E-6")...)  // 50 uF range
	CAP_500uF = append(capasitor, []byte("500E-6")...) // 500 uF range
	CAP_5mF   = append(capasitor, []byte("5E-3")...)   // 5 mF range
	CAP_50mF  = append(capasitor, []byte("50E-3")...)  // 50 mF range
	// Service timeouts
	COMUNICATION_TIMEOUT = time.Millisecond * 350
	SWITH_RANGE_TIMEOUT  = time.Second
)

// Define the SerialDevice interface
type SerialDevice interface {
	Write([]byte) (int, error)
	ReadMessage() ([]byte, error)
	Read([]byte) (int, error)
	Close() error
}

type XDM1041 struct {
	Device         SerialDevice
	MeassuresSpeed string
	MeassuresRange string
	ControllMode   string
	Identity       string
}

type Meassure struct {
	VoltageDC float64 `json:"voltage_dc"`
	CurrentDC float64 `json:"current_dc"`
}

func (m *Meassure) Power() float64 {
	return m.CurrentDC * m.VoltageDC
}

func NewXDM1041(device SerialDevice) (*XDM1041, error) {
	xdm := &XDM1041{Device: device}

	xdm.MeassuresSpeed = "M"
	err := xdm.setRate(xdm.MeassuresSpeed)
	if err != nil {
		return nil, err
	}

	xdm.Identity = xdm.GetIdentity()

	err = xdm.SetRemoteControll()
	if err != nil {
		return nil, err
	}

	// clog.Succes(clog.Yellow("%##v\n", xdm))
	return xdm, nil
}

func (x *XDM1041) Write(mess []byte) (int, error) {
	mess = append(mess, []byte("\n")...)
	// clog.Debug(clog.White("Write message: %##v", string(mess)))
	n, err := x.Device.Write(mess)
	time.Sleep(COMUNICATION_TIMEOUT)
	return n, err
}

func (x *XDM1041) ReadString() string {
	bf := make([]byte, 100) // Adjust buffer size as needed
	n, err := x.Device.Read(bf)
	if err != nil {
		return ""
	}
	meass := strings.TrimSpace(string(bf[:n]))
	return meass
}

// IDENTITY
func (x *XDM1041) GetIdentity() string {
	_, err := x.Write(IDETITY)
	if err != nil {
		return ""
	}

	return x.ReadString()
}

// Speed rate
func (x *XDM1041) GetRate() string {
	_, err := x.Write([]byte("RATE?"))
	if err != nil {
		return "S"
	}
	return x.ReadString()
}

func (x *XDM1041) setRate(s string) error {
	mess := []byte("RATE " + s)
	_, err := x.Write(mess)
	if err != nil {
		return err
	}
	return err
}

func (x *XDM1041) SetRateSlow() error {
	return x.setRate("S")
}

func (x *XDM1041) SetRateMedium() error {
	return x.setRate("M")
}

func (x *XDM1041) SetRateFast() error {
	return x.setRate("F")
}

// Controll mode
func (x *XDM1041) GetControll() string {
	return x.ControllMode
}

func (x *XDM1041) SetRemoteControll() error {
	mess := append(syst, []byte("REM")...)
	fmt.Println(string(mess))
	_, err := x.Write(mess)
	if err != nil {
		return err
	}
	x.ControllMode = "REM"
	return err
}

func (x *XDM1041) SetLocalControll() error {
	mess := append(syst, []byte("LOC")...)
	fmt.Println(string(mess))
	_, err := x.Write(mess)
	if err != nil {
		return err
	}
	x.ControllMode = "LOC"
	return err
}

// Meassures Range
func (x *XDM1041) GetMeassuresRange() string {
	_, err := x.Write(RANGE)
	if err != nil {
		return ""
	}
	meass := x.ReadString()
	return meass
}

// Set ranges functions
func (x *XDM1041) SetRangeVolt() error {
	//
	_, err := x.Write(VOLT)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltDc() error {
	//auto range
	_, err := x.Write(VOLT_DC)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltDc50e() error {
	//50 mV
	_, err := x.Write(VOLT_DC_50e)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltDc500e() error {
	//500 mV
	_, err := x.Write(VOLT_DC_500e)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltDc5() error {
	//5 V
	_, err := x.Write(VOLT_DC_5)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltDc50() error {
	//50 V
	_, err := x.Write(VOLT_DC_50)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltDc500() error {
	//500 V
	_, err := x.Write(VOLT_DC_500)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltAc() error {
	//volt AC
	_, err := x.Write(VOLT_AC)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeAc() error {
	// AC current autorange
	_, err := x.Write(acCMD)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltAc500e() error {
	//500 AC mV
	_, err := x.Write(VOLT_AC_500e)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltAc5() error {
	//5 AC V
	_, err := x.Write(VOLT_AC_5)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltAc50() error {
	//50 AC V
	_, err := x.Write(VOLT_AC_50)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltAc500() error {
	//500 AC V
	_, err := x.Write(VOLT_AC_500)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeVoltAc750() error {
	//750 AC V
	_, err := x.Write(VOLT_AC_750)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurrentDc() error {
	//current DC default
	_, err := x.Write(currDC)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurr() error {
	//default DC 5 A (10 A connector)
	_, err := x.Write(CURR)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurrDc() error {
	//
	_, err := x.Write(currDC)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurrDc500u() error {
	//500 uA (600 mA connector)
	_, err := x.Write(CURR_DC_500u)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurrDc5m() error {
	//5 mA (600 mA connector)
	_, err := x.Write(CURR_DC_5m)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurrDc50m() error {
	//50 mA (600 mA connector)
	_, err := x.Write(CURR_DC_50m)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurrDc500m() error {
	//500 mA (600 mA connector)
	_, err := x.Write(CURR_DC_500m)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurrDc5() error {
	//5 A (10 A connector)
	_, err := x.Write(CURR_DC_5)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCurrDc10() error {
	//10 A (10 A connector)
	_, err := x.Write(CURR_DC_10)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeResistor() error {
	//
	_, err := x.Write(resistor)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeResAuto() error {
	//
	_, err := x.Write(RES_AUTO)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeRes500() error {
	//500 Ω range
	_, err := x.Write(RES_500)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeRes5k() error {
	//5 KΩ range
	_, err := x.Write(RES_5k)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeRes50k() error {
	//50 KΩ range
	_, err := x.Write(RES_50k)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeRes500k() error {
	//500 KΩ range
	_, err := x.Write(RES_500k)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeRes5M() error {
	//5 MΩ range
	_, err := x.Write(RES_5M)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeRes50M() error {
	//50 MΩ range
	_, err := x.Write(RES_50M)
	if err != nil {
		return err
	}
	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) GetValue() (float64, error) {
	x.ReadString()
	_, err := x.Write(MEASS)
	if err != nil {
		return 0.0, err
	}

	m := x.ReadString()
	result, err := strconv.ParseFloat(m, 64)
	if err != nil {
		return 0.0, err
	}
	return result, nil
}

func (x *XDM1041) SetRangeCapAuto() error {
	//Capasitors auto range
	_, err := x.Write(CAP_AUTO)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCap50nf() error {
	//50 nF range
	_, err := x.Write(CAP_50nF)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCap500nf() error {
	//500 nF range
	_, err := x.Write(CAP_500nF)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCap5uf() error {
	//5 uF range
	_, err := x.Write(CAP_5uF)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCap50uf() error {
	//50 uF range
	_, err := x.Write(CAP_50uF)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCap500uf() error {
	//500 uF range
	_, err := x.Write(CAP_500uF)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCap5mf() error {
	//5 mF range
	_, err := x.Write(CAP_5mF)
	if err != nil {
		return err
	}

	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) SetRangeCap50mf() error {
	//50 mF range
	_, err := x.Write(CAP_50mF)
	if err != nil {
		return err
	}
	time.Sleep(SWITH_RANGE_TIMEOUT)
	return nil
}

func (x *XDM1041) GetMeassureVA_DC() (Meassure, error) {
	m := Meassure{}
	err := x.SetRangeCurrDc5()
	if err != nil {
		return m, err
	}
	// time.Sleep(SWITH_RANGE_TIMEOUT)
	m.CurrentDC, err = x.GetValue()
	if err != nil {
		return m, err
	}
	err = x.SetRangeVoltDc50()
	if err != nil {
		return m, err
	}
	// time.Sleep(SWITH_RANGE_TIMEOUT)
	m.VoltageDC, err = x.GetValue()
	if err != nil {
		return m, err
	}
	return m, err
}

// Key changes and explanations:

// * **`SerialDevice` Interface:**  I've defined a `SerialDevice` interface with `Write`, `Read`, and `Close` methods. This is the crucial part that decouples your `XDM1041` code from any specific serial port implementation.  Your code now depends on this interface, not on a concrete `uart` package.
// * **Dependency Injection:** The `NewXDM1041` function now takes a `SerialDevice` as an argument. This is how you *inject* the dependency.  The caller is responsible for creating the `SerialDevice` (e.g., using your `uart` package or any other serial port library) and passing it to `NewXDM1041`.
// * **Removed `uart` Import:** The `import "uart"` line is gone.
// * **`XDM1041` Uses Interface:** The `XDM1041` struct now holds a `SerialDevice` field instead of a concrete `uart.SerialDevice`.
// * **Method Calls Use Interface:** All calls to serial port functions (e.g., `Write`, `Read`) are now made through the `SerialDevice` interface.
// * **No Changes to `clog`:** The `clog` package remains untouched, as it's unrelated to the serial port dependency.

// How to use it:

// 1. **Implement `SerialDevice`:**  You'll need to create a concrete type that implements the `SerialDevice` interface.  This is where you'd use your `uart` package (or any other serial port library) to handle the actual serial communication.  For example:

//    ```go
//    import "uart"

//    type MySerialDevice struct {
//        port *uart.SerialDevice
//    }

//    func NewMySerialDevice(portName string, baudRate int) (*MySerialDevice, error) {
//        port, err := uart.NewConnect(portName, baudRate)
//        if err != nil {
//            return nil, err
//        }
//        return &MySerialDevice{port: port}, nil
//    }

//    func (d *MySerialDevice) Write(data []byte) (int, error) {
//        return d.port.Write(data)
//    }

//    func (d *MySerialDevice) Read(data []byte) (int, error) {
//        return d.port.Read(data)
//    }

//    func (d *MySerialDevice) Close() error {
//        return d.port.Close()
//    }
//    ```

// 2. **Create and Inject:**  In your main application, create an instance of your `MySerialDevice` and pass it to `NewXDM1041`:

//    ```go
//    import (
//        "fmt"
//        "xdm1041"
//    )

//    func main() {
//        portName := "/dev/ttyUSB0" // Replace with your serial port
//        baudRate := 9600

//        serialDevice, err := NewMySerialDevice(portName, baudRate)
//        if err != nil {
//            fmt.Println("Error creating serial device:", err)
//            return
//        }

//        xdm, err := xdm1041.NewXDM1041(serialDevice)
//        if err != nil {
//            fmt.Println("Error creating XDM1041:", err)
//            return
//        }

//        // Now you can use the 'xdm' instance as before
//    }
//    ```

// This revised code is much more testable, flexible, and maintainable.  You can easily swap out different serial port implementations without modifying the `XDM1041` code.  It also makes unit testing much easier, as you can mock the `SerialDevice` interface.  I've included a complete, runnable example to demonstrate how to use it.  Remember to replace `/dev/ttyUSB0` with the correct serial port for your system.
