package util

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: add different types of panics
