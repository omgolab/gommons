package gctest

var (
	scenarioPrefix  = []byte("Scenario ")
	givenPrefix     = []byte("Given ")
	whenPrefix      = []byte("When ")
	thenPrefix      = []byte("Then ")
	returnsErrors   = []byte("And returns errors")
	returnsNoErrors = []byte("And returns NO errors")
)

func ForScenario(name string) GivenStep {
	buf := make([]byte, 0, 100) // it generally take 80-100 chars for a scenario's full text
	if name != "" {
		buf = append(buf, append(scenarioPrefix, []byte(": "+name+",")...)...)
	}
	return &gs{
		base: base{
			buffer: buf,
		},
	}
}

type GivenStep interface {
	Given(name string) WhenStep
}

type WhenStep interface {
	When(name string) ThenStep
}

type ThenStep interface {
	Then(name string) ErrorStep
}

type ErrorStep interface {
	AndReturnsError() GWTSteps
	AndReturnsNoError() GWTSteps
}

type GWTSteps interface {
	ToString() string
}

// define base struct which contains previous full string
type base struct {
	buffer []byte
}

// define all necessary receivers for each step

type gs struct {
	base
}

type ws struct {
	base
}

type ts struct {
	base
}

type es struct {
	base
}

type gwt struct {
	base
}

// declare all the implementations of the interfaces

func (b *gs) Given(name string) WhenStep {
	b.buffer = append(b.buffer, append(givenPrefix, []byte(name+",")...)...)
	return &ws{
		base: b.base,
	}
}

func (b *ws) When(name string) ThenStep {
	b.buffer = append(b.buffer, append(whenPrefix, []byte(name+",")...)...)
	return &ts{
		base: b.base,
	}
}

func (b *ts) Then(name string) ErrorStep {
	b.buffer = append(b.buffer, append(thenPrefix, []byte(name+",")...)...)
	return &es{
		base: b.base,
	}
}

func (b *es) AndReturnsError() GWTSteps {
	b.buffer = append(b.buffer, returnsErrors...)
	return &gwt{
		base: b.base,
	}
}

func (b *es) AndReturnsNoError() GWTSteps {
	b.buffer = append(b.buffer, returnsNoErrors...)
	return &gwt{
		base: b.base,
	}
}

func (b *gwt) ToString() string {
	return string(b.buffer)
}
