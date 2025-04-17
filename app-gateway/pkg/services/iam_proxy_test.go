package services

import "testing"

func TestCheckoutResourceInstance(t *testing.T) {
	resource := "gcp:vcluster:cn-maanshan-a:*:instance/vc1p2iwcshyu"
	instance := checkoutResourceInstance(resource)
	t.Log(instance)
	if instance != "vc1p2iwcshyu" {
		t.Errorf("checkoutResourceInstance(%s) should return 'vc1p2iwcshyu', but got '%s'", resource, instance)
	}
}
