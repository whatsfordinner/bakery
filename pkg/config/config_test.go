package config

import (
	"context"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := map[string]struct {
		serviceNameOverride    string
		dbHostOverride         string
		rabbitHostOverride     string
		rabbitUsernameOverride string
		rabbitPasswordOverride string
		jaegerEndpointOverride string
	}{
		"no override":    {"", "", "", "", "", ""},
		"some overrides": {"", "foo", "", "", "qux", "xyzzy"},
		"all overrides":  {"plugh", "foo", "bar", "baz", "qux", "xyzzy"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.serviceNameOverride != "" {
				os.Setenv("SERVICE_NAME", test.serviceNameOverride)
			} else {
				os.Unsetenv("SERVICE_NAME")
			}

			if test.dbHostOverride != "" {
				os.Setenv("DB_HOST", test.dbHostOverride)
			} else {
				os.Unsetenv("DB_HOST")
			}

			if test.rabbitHostOverride != "" {
				os.Setenv("RABBIT_HOST", test.rabbitHostOverride)
			} else {
				os.Unsetenv("RABBIT_HOST")
			}

			if test.rabbitUsernameOverride != "" {
				os.Setenv("RABBIT_USERNAME", test.rabbitUsernameOverride)
			} else {
				os.Unsetenv("RABBIT_USERNAME")
			}

			if test.rabbitPasswordOverride != "" {
				os.Setenv("RABBIT_PASSWORD", test.rabbitPasswordOverride)
			} else {
				os.Unsetenv("RABBIT_PASSWORD")
			}

			if test.jaegerEndpointOverride != "" {
				os.Setenv("JAEGER_ENDPOINT", test.jaegerEndpointOverride)
			} else {
				os.Unsetenv("JAEGER_ENDPOINT")
			}

			c := GetConfig(context.Background())

			if test.serviceNameOverride != "" && c.ServiceName != test.serviceNameOverride {
				t.Fatalf("Override failed.\nExpected: %s\nGot: %s", test.serviceNameOverride, c.ServiceName)
			}

			if test.serviceNameOverride == "" && c.ServiceName != "bakery" {
				t.Fatalf("Default failed.\nExpected: %s\nGot: %s", "bakery", c.ServiceName)
			}

			if test.dbHostOverride != "" && c.DBHost != test.dbHostOverride {
				t.Fatalf("Override failed.\nExpected: %s\nGot: %s", test.dbHostOverride, c.DBHost)
			}

			if test.dbHostOverride == "" && c.DBHost != "127.0.0.1:6379" {
				t.Fatalf("Default failed.\nExpected: %s\nGot: %s", test.dbHostOverride, c.DBHost)
			}

			if test.rabbitHostOverride != "" && c.RabbitHost != test.rabbitHostOverride {
				t.Fatalf("Override failed.\nExpected: %s\nGot: %s", test.rabbitHostOverride, c.RabbitHost)
			}

			if test.rabbitHostOverride == "" && c.RabbitHost != "127.0.0.1:5672" {
				t.Fatalf("Default failed.\nExpected: %s\nGot: %s", test.rabbitHostOverride, c.RabbitHost)
			}

			if test.rabbitUsernameOverride != "" && c.RabbitUsername != test.rabbitUsernameOverride {
				t.Fatalf("Override failed.\nExpected: %s\nGot: %s", test.rabbitUsernameOverride, c.RabbitUsername)
			}

			if test.rabbitUsernameOverride == "" && c.RabbitUsername != "guest" {
				t.Fatalf("Default failed.\nExpected: %s\nGot: %s", test.rabbitUsernameOverride, c.RabbitUsername)
			}

			if test.rabbitPasswordOverride != "" && c.RabbitPassword != test.rabbitPasswordOverride {
				t.Fatalf("Override failed.\nExpected: %s\nGot: %s", test.rabbitPasswordOverride, c.RabbitPassword)
			}

			if test.rabbitPasswordOverride == "" && c.RabbitPassword != "guest" {
				t.Fatalf("Default failed.\nExpected: %s\nGot: %s", test.rabbitPasswordOverride, c.RabbitPassword)
			}

			if test.jaegerEndpointOverride != "" && c.JaegerEndpoint != test.jaegerEndpointOverride {
				t.Fatalf("Override failed.\nExpected: %s\nGot: %s", test.jaegerEndpointOverride, c.JaegerEndpoint)
			}

			if test.jaegerEndpointOverride == "" && c.JaegerEndpoint != "" {
				t.Fatalf("Default failed.\nExpected: %s\nGot: %s", test.jaegerEndpointOverride, c.JaegerEndpoint)
			}
		})
	}
}
