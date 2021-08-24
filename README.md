# openers
Gate and garage-opening code for Raspberry Pis and ESP32s.

The plan is to be able to open my gate and garage remotely, by driving simple
single-frequency transmitters first from a Raspberry Pi, and later from an
ESP32.

# Commandline examples

```bash
# Open (my) gate
sudo chrt -f -r 99 openers secplus transmitv2 --rolling=123456789 --fixed=1222022221850123456789 --pin=12
```

# Building

## For (my) Raspberry Pi Zero W

```bash
GOOS=linux GOARCH=arm GOARM=5 go build . && scp openers pi@pizero:openers
```

# Sub-packages

## secplus

Package secplus implements Security+2.0 encoding.

Writing it would not have been possible without the work done by @argilo in
decoding the Security+2.0 protocol for their excellent [Python `secplus`
package](https://github.com/argilo/secplus), and the helpful debugging by
@acoursen in understanding longer transmissions (argilo/secplus#6).

## megacode

Package megacode implements MegaCode encoding.

Writing it would have been difficult without the work done by CuVoodoo in
describing and decoding the MegaCode protocol on their excellent [MegaCode
hacking page](https://wiki.cuvoodoo.info/doku.php?id=megacode).

## gpiod

Package gpiod exists simply to wrap github.com/warthog618/gpiod on Linux, and
provide a dummy implementation on other platforms, so that things will still
compile. Only functions and types actually used by the other code have been
added to the dummy implementation.

# Todos

Next steps for my development (likely to get done soon):

- [x] Security+: Write code to actually transmit by driving a Raspberry Pi pin
- [x] Write a commandline app for testing
- [ ] Implement MegaCode encoding
- [ ] Implement MegaCode Raspberry Pi transmission

Future (PRs welcome!):

- [ ] Convert these TODOs into issues
- [ ] Add Security+2.0 decoding
- [ ] Add Security+ encoding/decoding