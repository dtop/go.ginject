package ginject

import (
	"testing"
	"time"
)

type (
	TestItf interface {
		GetA() string
	}

	FooTest struct {
		a string
		b bool
	}

	BloTest struct {
		i int
		f float32
	}

	DooTest struct {
		x string
		t time.Time
	}

	Target struct {
		Foo TestItf  `inject:"foo"`
		Blo *BloTest `inject:"blo"`
		Doo *DooTest `inject:"doo"`
	}

	NegativeTarget struct {
		Xoo FooTest `inject:"loo"`
	}

	FailingTarget struct {
		Xoo *DooTest `inject:"foo"`
	}
)

func (ft FooTest) GetA() string {
	return ft.a
}

func xfory(foo *FooTest) TestItf {
	return foo
}

func TestInj_Get(t *testing.T) {

	ft := &FooTest{a: "tea", b: false}

	inject := New()
	inject.Register(IService("foo", ft))

	result := &FooTest{}
	err := inject.Get("foo", result)

	if err != nil {
		t.Error(err)
	}

	if result.a != ft.a {
		t.Failed()
	}

	result1 := FooTest{}
	err = inject.Get("foo", result1)

	if err == nil {
		t.Failed()
	}

	err = inject.Get("soo", &result1)

	if err == nil {
		t.Failed()
	}

	result2 := xfory(&FooTest{})
	err = inject.Get("foo", &result2)

	if err != nil {
		t.Error(err)
	}

	if result2.GetA() != "tea" {
		t.Error("interface was not injected")
	}

	result3 := &BloTest{}
	err = inject.Get("foo", &result3)

	if err == nil {
		t.Error("wrong type was injected")
	}
}

func TestInj_Apply(t *testing.T) {

	inj := New()
	inj.RegisterByName("foo", FooTest{a: "moo", b: true})
	inj.RegisterByName("blo", &BloTest{i: 19, f: 1.4})
	inj.RegisterLazy("doo", func() interface{} {

		return &DooTest{x: "time", t: time.Now()}
	})

	target := &Target{}
	if err := inj.Apply(target); err != nil {
		t.Error(err)
	}

	if target.Foo.GetA() != "moo" {
		t.Error("not injected (1)")
	}

	if target.Blo == nil || target.Blo.f != 1.4 {
		t.Error("not injected (2)")
	}

	if err := inj.Apply("moo"); err == nil {
		t.Error("should not apply to non struct")
	}

	ntarget := &NegativeTarget{}
	if err := inj.Apply(ntarget); err == nil {
		t.Error("should not apply to nonexistant")
	}

	ftarget := &FailingTarget{}
	if err := inj.Apply(ftarget); err == nil {
		t.Error("should not apply to incompatible")
	}
}
