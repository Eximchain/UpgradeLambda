package main

import (
	"testing"
)

func Test_compare(t *testing.T) {
	compare(`(([0-9]+)\.?)+`, "Vault v0.10.2 ('3ee0802ed08cb7f4046c2151ec4671a076b76166')", "0.9.1", 0, 1, 2, 3)
	compare(`(([0-9]+)\.?)+`, "Consul v1.9.3", "1.10.4", 0, 1, 2, 3)
	compare(`(([0-9]+)\.?)+`, "Consul v1.2.3", "1.2.4", 0, 1, 2, 3)
}

func Test_compareV(t *testing.T) {
	compareV("1.1", "1.9")
	compareV("1.1", "1.1.1")
}
