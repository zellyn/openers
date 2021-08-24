/*
Package gpiod exists simply to wrap github.com/warthog618/gpiod on Linux, and
provide a dummy implementation on other platforms, so that things will still
compile. Only functions and types actually used by the other code have been
added to the dummy implementation.
*/
package gpiod
