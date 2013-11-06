package inject_test

import (
	"testing"

	"github.com/ParsePlatform/go.inject"
)

type Answerable interface {
	Answer() int
}

type TypeAnswerStruct struct {
	answer  int
	private int
}

func (t *TypeAnswerStruct) Answer() int {
	return t.answer
}

type TypeNestedStruct struct {
	A *TypeAnswerStruct `inject:""`
}

func (t *TypeNestedStruct) Answer() int {
	return t.A.Answer()
}

func TestRequireTag(t *testing.T) {
	var v struct {
		A *TypeAnswerStruct
		B *TypeNestedStruct `inject:""`
	}

	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.A != nil {
		t.Fatal("v.A is not nil")
	}
	if v.B == nil {
		t.Fatal("v.B is nil")
	}
}

type TypeWithNonPointerInject struct {
	A int `inject:""`
}

func TestErrorOnNonPointerInject(t *testing.T) {
	var a TypeWithNonPointerInject
	err := inject.Populate(&a)
	if err == nil {
		t.Fatalf("expected error for %+v", a)
	}

	const msg = "found inject tag on unsupported field A in type *inject_test.TypeWithNonPointerInject"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeWithNonPointerStructInject struct {
	A *int `inject:""`
}

func TestErrorOnNonPointerStructInject(t *testing.T) {
	var a TypeWithNonPointerStructInject
	err := inject.Populate(&a)
	if err == nil {
		t.Fatalf("expected error for %+v", a)
	}

	const msg = "found inject tag on unsupported field A in type *inject_test.TypeWithNonPointerStructInject"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

func TestInjectSimple(t *testing.T) {
	var v struct {
		A *TypeAnswerStruct `inject:""`
		B *TypeNestedStruct `inject:""`
	}

	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.A == nil {
		t.Fatal("v.A is nil")
	}
	if v.B == nil {
		t.Fatal("v.B is nil")
	}
	if v.B.A == nil {
		t.Fatal("v.B.A is nil")
	}
	if v.A != v.B.A {
		t.Fatal("got different instances of A")
	}
}

func TestDoesNotOverwrite(t *testing.T) {
	a := &TypeAnswerStruct{}
	var v struct {
		A *TypeAnswerStruct `inject:""`
		B *TypeNestedStruct `inject:""`
	}
	v.A = a
	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.A != a {
		t.Fatal("original A was lost")
	}
	if v.B == nil {
		t.Fatal("v.B is nil")
	}
}

func TestPrivate(t *testing.T) {
	var v struct {
		A *TypeAnswerStruct `inject:"private"`
		B *TypeNestedStruct `inject:""`
	}

	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.A == nil {
		t.Fatal("v.A is nil")
	}
	if v.B == nil {
		t.Fatal("v.B is nil")
	}
	if v.B.A == nil {
		t.Fatal("v.B.A is nil")
	}
	if v.A == v.B.A {
		t.Fatal("got the same A")
	}
}

type TypeWithJustColon struct {
	A *TypeAnswerStruct `inject:`
}

func TestTagWithJustColon(t *testing.T) {
	var a TypeWithJustColon
	err := inject.Populate(&a)
	if err == nil {
		t.Fatalf("expected error for %+v", a)
	}

	const msg = "unexpected tag format `inject:` for field A in type *inject_test.TypeWithJustColon"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeWithOpenQuote struct {
	A *TypeAnswerStruct `inject:"`
}

func TestTagWithOpenQuote(t *testing.T) {
	var a TypeWithOpenQuote
	err := inject.Populate(&a)
	if err == nil {
		t.Fatalf("expected error for %+v", a)
	}

	const msg = "unexpected tag format `inject:\"` for field A in type *inject_test.TypeWithOpenQuote"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

func TestProvideNonPointer(t *testing.T) {
	var g inject.Graph
	var i int
	err := g.Provide(&inject.Object{Value: i})
	if err == nil {
		t.Fatal("expected error")
	}

	const msg = "expected unnamed object value to be a pointer to a struct but got type int with value 0"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

func TestProvideNonPointerStruct(t *testing.T) {
	var g inject.Graph
	var i *int
	err := g.Provide(&inject.Object{Value: i})
	if err == nil {
		t.Fatal("expected error")
	}

	const msg = "expected unnamed object value to be a pointer to a struct but got type *int with value <nil>"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

func TestProvideTwoOfTheSame(t *testing.T) {
	var g inject.Graph
	a := TypeAnswerStruct{}
	err := g.Provide(&inject.Object{Value: &a})
	if err != nil {
		t.Fatal(err)
	}

	err = g.Provide(&inject.Object{Value: &a})
	if err == nil {
		t.Fatal("expected error")
	}

	const msg = "provided two unnamed instances of type *inject_test.TypeAnswerStruct"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

func TestProvideTwoOfTheSameWithPopulate(t *testing.T) {
	a := TypeAnswerStruct{}
	err := inject.Populate(&a, &a)
	if err == nil {
		t.Fatal("expected error")
	}

	const msg = "provided two unnamed instances of type *inject_test.TypeAnswerStruct"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

func TestProvideTwoWithTheSameName(t *testing.T) {
	var g inject.Graph
	const name = "foo"
	a := TypeAnswerStruct{}
	err := g.Provide(&inject.Object{Value: &a, Name: name})
	if err != nil {
		t.Fatal(err)
	}

	err = g.Provide(&inject.Object{Value: &a, Name: name})
	if err == nil {
		t.Fatal("expected error")
	}

	const msg = "provided two instances named foo"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

func TestNamedInstanceWithDependencies(t *testing.T) {
	var g inject.Graph
	a := &TypeNestedStruct{}
	if err := g.Provide(&inject.Object{Value: a, Name: "foo"}); err != nil {
		t.Fatal(err)
	}

	var c struct {
		A *TypeNestedStruct `inject:"foo"`
	}
	if err := g.Provide(&inject.Object{Value: &c}); err != nil {
		t.Fatal(err)
	}

	if err := g.Populate(); err != nil {
		t.Fatal(err)
	}

	if c.A.A == nil {
		t.Fatal("c.A.A was not injected")
	}
}

func TestTwoNamedInstances(t *testing.T) {
	var g inject.Graph
	a := &TypeAnswerStruct{}
	b := &TypeAnswerStruct{}
	if err := g.Provide(&inject.Object{Value: a, Name: "foo"}); err != nil {
		t.Fatal(err)
	}

	if err := g.Provide(&inject.Object{Value: b, Name: "bar"}); err != nil {
		t.Fatal(err)
	}

	var c struct {
		A *TypeAnswerStruct `inject:"foo"`
		B *TypeAnswerStruct `inject:"bar"`
	}
	if err := g.Provide(&inject.Object{Value: &c}); err != nil {
		t.Fatal(err)
	}

	if err := g.Populate(); err != nil {
		t.Fatal(err)
	}

	if c.A != a {
		t.Fatal("did not find expected c.A")
	}
	if c.B != b {
		t.Fatal("did not find expected c.B")
	}
}

type TypeWithMissingNamed struct {
	A *TypeAnswerStruct `inject:"foo"`
}

func TestTagWithMissingNamed(t *testing.T) {
	var a TypeWithMissingNamed
	err := inject.Populate(&a)
	if err == nil {
		t.Fatalf("expected error for %+v", a)
	}

	const msg = "did not find object named foo required by field A in type *inject_test.TypeWithMissingNamed"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

func TestCompleteProvides(t *testing.T) {
	var g inject.Graph
	var v struct {
		A *TypeAnswerStruct `inject:""`
	}
	if err := g.Provide(&inject.Object{Value: &v, Complete: true}); err != nil {
		t.Fatal(err)
	}

	if err := g.Populate(); err != nil {
		t.Fatal(err)
	}
	if v.A != nil {
		t.Fatal("v.A was not nil")
	}
}

func TestCompleteNamedProvides(t *testing.T) {
	var g inject.Graph
	var v struct {
		A *TypeAnswerStruct `inject:""`
	}
	if err := g.Provide(&inject.Object{Value: &v, Complete: true, Name: "foo"}); err != nil {
		t.Fatal(err)
	}

	if err := g.Populate(); err != nil {
		t.Fatal(err)
	}
	if v.A != nil {
		t.Fatal("v.A was not nil")
	}
}

type TypeInjectInterfaceMissing struct {
	Answerable Answerable `inject:""`
}

func TestInjectInterfaceMissing(t *testing.T) {
	var v TypeInjectInterfaceMissing
	err := inject.Populate(&v)
	if err == nil {
		t.Fatal("did not find expected error")
	}

	const msg = "found no assignable value for field Answerable in type *inject_test.TypeInjectInterfaceMissing"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeInjectInterface struct {
	Answerable Answerable        `inject:""`
	A          *TypeAnswerStruct `inject:""`
}

func TestInjectInterface(t *testing.T) {
	var v TypeInjectInterface
	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.Answerable == nil || v.Answerable != v.A {
		t.Fatalf(
			"expected the same but got Answerable = %T %+v / A = %T %+v",
			v.Answerable,
			v.Answerable,
			v.A,
			v.A,
		)
	}
}

type TypeWithInvalidNamedType struct {
	A *TypeNestedStruct `inject:"foo"`
}

func TestInvalidNamedInstanceType(t *testing.T) {
	var g inject.Graph
	a := &TypeAnswerStruct{}
	if err := g.Provide(&inject.Object{Value: a, Name: "foo"}); err != nil {
		t.Fatal(err)
	}

	var c TypeWithInvalidNamedType
	if err := g.Provide(&inject.Object{Value: &c}); err != nil {
		t.Fatal(err)
	}

	err := g.Populate()
	if err == nil {
		t.Fatal("did not find expected error")
	}

	const msg = "object named foo of type *inject_test.TypeNestedStruct is not assignable to field A in type *inject_test.TypeWithInvalidNamedType"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeWithInjectOnPrivateField struct {
	a *TypeAnswerStruct `inject:""`
}

func TestInjectOnPrivateField(t *testing.T) {
	var a TypeWithInjectOnPrivateField
	err := inject.Populate(&a)
	if err == nil {
		t.Fatal("did not find expected error")
	}

	const msg = "inject requested on unexported field a in type *inject_test.TypeWithInjectOnPrivateField"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeWithInjectOnPrivateInterfaceField struct {
	a Answerable `inject:""`
}

func TestInjectOnPrivateInterfaceField(t *testing.T) {
	var a TypeWithInjectOnPrivateField
	err := inject.Populate(&a)
	if err == nil {
		t.Fatal("did not find expected error")
	}

	const msg = "inject requested on unexported field a in type *inject_test.TypeWithInjectOnPrivateField"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeInjectPrivateInterface struct {
	Answerable Answerable        `inject:"private"`
	B          *TypeNestedStruct `inject:""`
}

func TestInjectPrivateInterface(t *testing.T) {
	var v TypeInjectPrivateInterface
	err := inject.Populate(&v)
	if err == nil {
		t.Fatal("did not find expected error")
	}

	const msg = "found private inject tag on interface field Answerable in type *inject_test.TypeInjectPrivateInterface"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeInjectTwoSatisfyInterface struct {
	Answerable Answerable        `inject:""`
	A          *TypeAnswerStruct `inject:""`
	B          *TypeNestedStruct `inject:""`
}

func TestInjectTwoSatisfyInterface(t *testing.T) {
	var v TypeInjectTwoSatisfyInterface
	err := inject.Populate(&v)
	if err == nil {
		t.Fatal("did not find expected error")
	}

	const msg = "found two assignable values for field Answerable in type *inject_test.TypeInjectTwoSatisfyInterface. one type *inject_test.TypeAnswerStruct with value &{0 0} and another type *inject_test.TypeNestedStruct with value <*inject_test.TypeNestedStruct Value>"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeInjectNamedTwoSatisfyInterface struct {
	Answerable Answerable        `inject:""`
	A          *TypeAnswerStruct `inject:""`
	B          *TypeNestedStruct `inject:""`
}

func TestInjectNamedTwoSatisfyInterface(t *testing.T) {
	var g inject.Graph
	var v TypeInjectNamedTwoSatisfyInterface
	if err := g.Provide(&inject.Object{Name: "foo", Value: &v}); err != nil {
		t.Fatal(err)
	}

	err := g.Populate()
	if err == nil {
		t.Fatal("was expecting error")
	}

	const msg = "found two assignable values for field Answerable in type *inject_test.TypeInjectNamedTwoSatisfyInterface. one type *inject_test.TypeAnswerStruct with value &{0 0} and another type *inject_test.TypeNestedStruct with value <*inject_test.TypeNestedStruct Value>"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeWithInjectNamedOnPrivateInterfaceField struct {
	a Answerable `inject:""`
}

func TestInjectNamedOnPrivateInterfaceField(t *testing.T) {
	var g inject.Graph
	var v TypeWithInjectNamedOnPrivateInterfaceField
	if err := g.Provide(&inject.Object{Name: "foo", Value: &v}); err != nil {
		t.Fatal(err)
	}

	err := g.Populate()
	if err == nil {
		t.Fatal("was expecting error")
	}

	const msg = "inject requested on unexported field a in type *inject_test.TypeWithInjectNamedOnPrivateInterfaceField"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeWithNonPointerNamedInject struct {
	A int `inject:"foo"`
}

func TestErrorOnNonPointerNamedInject(t *testing.T) {
	var g inject.Graph
	if err := g.Provide(&inject.Object{Name: "foo", Value: 42}); err != nil {
		t.Fatal(err)
	}

	var v TypeWithNonPointerNamedInject
	if err := g.Provide(&inject.Object{Value: &v}); err != nil {
		t.Fatal(err)
	}

	if err := g.Populate(); err != nil {
		t.Fatal(err)
	}

	if v.A != 42 {
		t.Fatalf("expected v.A = 42 but got %d", v.A)
	}
}

func TestInjectInline(t *testing.T) {
	var v struct {
		Inline struct {
			A *TypeAnswerStruct `inject:""`
			B *TypeNestedStruct `inject:""`
		} `inject:""`
	}

	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.Inline.A == nil {
		t.Fatal("v.Inline.A is nil")
	}
	if v.Inline.B == nil {
		t.Fatal("v.Inline.B is nil")
	}
	if v.Inline.B.A == nil {
		t.Fatal("v.Inline.B.A is nil")
	}
	if v.Inline.A != v.Inline.B.A {
		t.Fatal("got different instances of A")
	}
}

type TypeWithInlineStructWithPrivate struct {
	Inline struct {
		A *TypeAnswerStruct `inject:""`
		B *TypeNestedStruct `inject:""`
	} `inject:"private"`
}

func TestInjectInlinePrivate(t *testing.T) {
	var v TypeWithInlineStructWithPrivate
	err := inject.Populate(&v)
	if err == nil {
		t.Fatal("was expecting an error")
	}

	const msg = "cannot use private inject on inline struct on field Inline in type *inject_test.TypeWithInlineStructWithPrivate"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}

type TypeWithStructValue struct {
	Inline TypeNestedStruct `inject:""`
}

func TestInjectWithStructValue(t *testing.T) {
	var v TypeWithStructValue
	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.Inline.A == nil {
		t.Fatal("v.Inline.A is nil")
	}
}

func TestPrivateIsFollowed(t *testing.T) {
	var v struct {
		A *TypeNestedStruct `inject:"private"`
	}

	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.A.A == nil {
		t.Fatal("v.A.A is nil")
	}
}

func TestDoesNotOverwriteInterface(t *testing.T) {
	a := &TypeAnswerStruct{}
	var v struct {
		A Answerable        `inject:""`
		B *TypeNestedStruct `inject:""`
	}
	v.A = a
	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.A != a {
		t.Fatal("original A was lost")
	}
	if v.B == nil {
		t.Fatal("v.B is nil")
	}
}

func TestInterfaceIncludingPrivate(t *testing.T) {
	var v struct {
		A Answerable        `inject:""`
		B *TypeNestedStruct `inject:"private"`
		C *TypeAnswerStruct `inject:""`
	}
	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.A == nil {
		t.Fatal("v.A is nil")
	}
	if v.B == nil {
		t.Fatal("v.B is nil")
	}
	if v.C == nil {
		t.Fatal("v.C is nil")
	}
	if v.A != v.C {
		t.Fatal("v.A != v.C")
	}
	if v.A == v.B {
		t.Fatal("v.A == v.B")
	}
}

func TestInjectMap(t *testing.T) {
	var v struct {
		A map[string]int `inject:"private"`
	}
	if err := inject.Populate(&v); err != nil {
		t.Fatal(err)
	}
	if v.A == nil {
		t.Fatal("v.A is nil")
	}
}

type TypeInjectWithMapWithoutPrivate struct {
	A map[string]int `inject:""`
}

func TestInjectMapWithoutPrivate(t *testing.T) {
	var v TypeInjectWithMapWithoutPrivate
	err := inject.Populate(&v)
	if err == nil {
		t.Fatalf("expected error for %+v", v)
	}

	const msg = "inject on map field A in type *inject_test.TypeInjectWithMapWithoutPrivate must be named or private"
	if err.Error() != msg {
		t.Fatalf("expected:\n%s\nactual:\n%s", msg, err.Error())
	}
}
