package routes

import (
	"reflect"
	"testing"
)

func TestAllowedOriginsDefaultsWhenUnset(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "")
	got := allowedOrigins()
	if !reflect.DeepEqual(got, []string{"http://localhost:3000", "http://frontend:3000"}) {
		t.Errorf("expected dev defaults, got %v", got)
	}
}

func TestAllowedOriginsParsesEnvList(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", " https://app.example.com , https://admin.example.com ")
	got := allowedOrigins()
	want := []string{"https://app.example.com", "https://admin.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestAllowedOriginsSkipsEmptyEntries(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", ",, ,https://only.example.com,,")
	got := allowedOrigins()
	want := []string{"https://only.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestAllowedOriginsFallsBackWhenOnlySeparators(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", " , , ")
	got := allowedOrigins()
	if !reflect.DeepEqual(got, []string{"http://localhost:3000", "http://frontend:3000"}) {
		t.Errorf("expected dev defaults for blank-only list, got %v", got)
	}
}
