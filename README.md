This is a Golang driver for the AMG8833 8x8 Thermal Camera Sensor

## Installation

    go get -u github.com/jweissig/amg8833

## Usage

    import "github.com/jweissig/amg8833"

For now, `amg8833` only supports retrieving thermal pixel data, so no interrupts yet.

```go
  amg, err := amg8833.NewAMG8833(&amg8833.Opts{
    Device: "/dev/i2c-1",
    Mode:   amg8833.AMG88xxNormalMode,
    Reset:  amg8833.AMG88xxInitialReset,
    FPS:    amg8833.AMG88xxFPS10,
  })
  if err != nil {
    panic(err)
  }

  ticker := time.NewTicker(1 * time.Second)

  for {
    grid := amg.ReadPixels()
    fmt.Println(grid)
    <-ticker.C
  }
```

## Acknowledgements

This library is basically a golang port of [Adafruit's AMG88xx Library](https://github.com/adafruit/Adafruit_AMG88xx/). I also used mstahl's [tsl2591](https://github.com/mstahl/tsl2591) driver as a starting point for the amg8833 one.
