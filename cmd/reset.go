package cmd

func Reset(oid string) {
	updateRef("HEAD", RefValue{symbolic: false, value: oid}, true)
}
