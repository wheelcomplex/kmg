package kmgHttp

import "net/url"

//sort url query to unique it
func SortUrlQuery(vurl *url.URL) {
	query := vurl.Query()
	vurl.RawQuery = query.Encode()
}
