package display

import (
	"os"

	"github.com/orisano/gopl/ch07/ex16/eval"
)

func ExampleDisplayEval() {
	e, _ := eval.Parse("sqrt(A / pi)")
	Display("e", e)
	// Output:
	// Display e (eval.call):
	// e.fn = "sqrt"
	// e.args[0].type = eval.binary
	// e.args[0].value.op = 47
	// e.args[0].value.x.type = eval.Var
	// e.args[0].value.x.value = "A"
	// e.args[0].value.y.type = eval.Var
	// e.args[0].value.y.value = "pi"
}

func ExampleDisplayMovie() {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},

		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	Display("strangelove", strangelove)
	/*
		// Output:
		// strangelove.Title = "Dr. Strangelove"
		// strangelove.Subtitle = "How I Learned to Stop Worrying and Love the Bomb"
		// strangelove.Year = 1964
		// strangelove.Color = false
		// strangelove.Actor["Gen. Buck Turgidson"] = "George C. Scott"
		// strangelove.Actor["Brig. Gen. Jack D. Ripper"] = "Sterling Hayden"
		// strangelove.Actor["Maj. T.J. \"King\" Kong"] = "Slim Pickens"
		// strangelove.Actor["Dr. Strangelove"] = "Peter Sellers"
		// strangelove.Actor["Grp. Capt. Lionel Mandrake"] = "Peter Sellers"
		// strangelove.Actor["Pres. Merkin Muffley"] = = "Peter Sellers"
		// strangelove.Oscars[0] = "Best Actor (Nomin.)"
		// strangelove.Oscars[1] = "Best Adapted Screenplay (Nomin.)"
		// strangelove.Oscars[2] = "Best Director (Nomin.)"
		// strangelove.Oscars[3] = "Best Picture (Nomin.)"
	*/

}

func ExampleDisplayStderr() {
	Display("os.Stderr", os.Stderr)
	// Output:
	// Display os.Stderr (*os.File):
	// (*(*os.Stderr).file).pfd.fdmu.state = 0
	// (*(*os.Stderr).file).pfd.fdmu.rsema = 0
	// (*(*os.Stderr).file).pfd.fdmu.wsema = 0
	// (*(*os.Stderr).file).pfd.Sysfd = 2
	// (*(*os.Stderr).file).pfd.pd.runtimeCtx = uintptr value
	// (*(*os.Stderr).file).pfd.iovecs = nil
	// (*(*os.Stderr).file).pfd.csema = 0
	// (*(*os.Stderr).file).pfd.IsStream = true
	// (*(*os.Stderr).file).pfd.ZeroReadIsEOF = true
	// (*(*os.Stderr).file).pfd.isFile = true
	// (*(*os.Stderr).file).pfd.isBlocking = true
	// (*(*os.Stderr).file).name = "/dev/stderr"
	// (*(*os.Stderr).file).dirinfo = nil
	// (*(*os.Stderr).file).nonblock = false
	// (*(*os.Stderr).file).stdoutOrErr = true
}

func ExampleDisplayI() {
	var i interface{} = 3

	Display("i", i)
	Display("&i", &i)

	// Output:
	// Display i (int):
	// i = 3
	// Display &i (*interface {}):
	// (*&i).type = int
	// (*&i).value = 3
}

func ExampleDisplayStructKey() {
	type Person struct {
		Name   string
		Height int
	}
	m := map[Person]int{
		Person{"Foo", 178}: 10,
	}

	Display("m", m)
	// Output:
	// Display m (map[display.Person]int):
	// m[display.Person{.Name = "Foo", .Height = 178}] = 10
}

func ExampleDisplayArrayKey() {

	m := map[[3]string]string{
		[3]string{"A", "B", "C"}: "foo",
	}

	Display("m", m)
	// Output:
	// Display m (map[[3]string]string):
	// m[[3]string{"A", "B", "C"}] = "foo"
}
