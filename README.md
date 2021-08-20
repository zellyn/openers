# openers
Gate and garage-opening code for Raspberry Pis and ESP32s.

The plan is to be able to open my gate and garage remotely, by driving simple
single-frequency transmitters first from a Raspberry Pi, and later from an
ESP32.

# Sub-packages

## secplus

Package secplus implements Security+2.0 encoding.

Writing it would not have been possible without the work done by @argilo in
decoding the Security+2.0 protocol for his excellent [Python `secplus`
package](https://github.com/argilo/secplus), and the helpful debugging by
@acoursen in understanding longer transmissions (argilo/secplus#6).

# Todos

Next steps for my development (likely to get done soon):

- [ ] Security+: Write code to actually transmit by driving a Raspberry Pi pin
- [ ] Write a commandline app for testing
- [ ] Implement MegaCode encoding
- [ ] Implement MegaCode Raspberry Pi transmission

Future (PRs welcome!):

- [ ] Convert these TODOs into issues
- [ ] Add Security+2.0 decoding
- [ ] Add Security+ encoding/decoding