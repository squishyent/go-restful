func main() {
	uri := NewURIBuilder("http", "localhost", 9999, "/hier/{daar}/maar")
	uri.PathParam("daar", "xyz")
	uri.QueryParam("unicode", "◎")
	uri.QueryParam("unicode", "X")
	uri.QueryParam("simple", "Y")
	println(uri.Build())
}